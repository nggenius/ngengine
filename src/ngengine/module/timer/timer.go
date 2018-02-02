package timer

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"
)

const (
	TVR_BITS          = 8
	TVN_BITS          = 6
	TVR_SIZE          = 1 << TVR_BITS // 256
	TVN_SIZE          = 1 << TVN_BITS // 64
	TVR_MASK          = TVR_SIZE - 1  // 255		1111 1111
	TVN_MASK          = TVN_SIZE - 1  // 63		    0011 1111
	MIN_TICK_INTERVAL = 1e6           // nanoseconds, 1ms
	MAXN_LEVEL        = 5
	FPS               = 50
)

type timer struct {
	id     int64
	expire int64
	node   *list.Element
	root   *list.List
	task   *TimerTask
}

type TimerQueue struct {
	tickTime      int64
	ticks         int64
	nextTimerId   int64
	tvec          [MAXN_LEVEL][]*list.List
	pendingTimers *list.List
	mutex         sync.Mutex
	last          int64
	curr          int64
	ti            int64
}

func New() *TimerQueue {
	tq := &TimerQueue{
		tickTime:      now(),
		ticks:         0,
		nextTimerId:   0,
		pendingTimers: list.New(),
		last:          0,
		curr:          0,
		ti:            int64(1e9 / FPS),
	}
	for i := 0; i < MAXN_LEVEL; i++ {
		if i == 0 {
			tq.tvec[i] = make([]*list.List, TVR_SIZE)
		} else {
			tq.tvec[i] = make([]*list.List, TVN_SIZE)
		}
		for j := 0; j < len(tq.tvec[i]); j++ {
			tq.tvec[i][j] = list.New()
		}
	}
	return tq
}

func (tq *TimerQueue) Schedule(delay int64, task *TimerTask) int64 {
	delay = delay * 1e6
	if delay < MIN_TICK_INTERVAL {
		delay = MIN_TICK_INTERVAL
	}
	ev := &timer{
		id:     tq.genID(),
		expire: atomic.LoadInt64(&(tq.tickTime)) + delay,
		node:   nil,
		root:   nil,
		task:   task,
	}
	tq.mutex.Lock()
	tq.pendingTimers.PushBack(ev)
	tq.mutex.Unlock()
	return ev.id
}

func (tq *TimerQueue) Run() {
	if tq.last == 0 {
		tq.last = now()
	}

	tq.curr = now()
	diff := tq.curr - tq.last
	if diff < tq.ti {
		return
	}
	tq.tick(diff)
	tq.last = tq.curr
}

func (tq *TimerQueue) genID() int64 {
	tq.nextTimerId++
	return tq.nextTimerId
}

func now() int64 {
	return time.Now().UnixNano()
}

func (tq *TimerQueue) addTimer(t *timer) int64 {
	ticks := (t.expire - tq.tickTime) / MIN_TICK_INTERVAL
	if ticks < 0 {
		ticks = 0
	}
	idx := tq.ticks + ticks
	level := 0

	var threshold1 int64 = TVR_SIZE
	var threshold2 int64 = 1 << (TVR_BITS + TVN_BITS)
	var threshold3 int64 = 1 << (TVR_BITS + 2*TVN_BITS)
	var threshold4 int64 = 1 << (TVR_BITS + 3*TVN_BITS)

	if ticks < threshold1 {
		idx = idx & TVR_MASK
		level = 0
	} else if ticks < threshold2 {
		idx = (idx >> (TVR_BITS)) & TVN_MASK
		level = 1
	} else if ticks < threshold3 {
		idx = (idx >> (TVR_BITS + TVN_BITS)) & TVN_MASK
		level = 2
	} else if ticks < threshold4 {
		idx = (idx >> (TVR_BITS + 2*TVN_BITS)) & TVN_MASK
		level = 3
	} else {
		idx = (idx >> (TVR_BITS + 3*TVN_BITS)) & TVN_MASK
		level = 4
	}

	vec := tq.tvec[level][idx]
	t.node = vec.PushBack(t)
	t.root = vec
	return t.id
}

func (tq *TimerQueue) cascade(n uint32) uint32 {
	idx := uint32(tq.ticks>>(TVR_BITS+(n-1)*TVN_BITS)) & TVN_MASK
	vec := tq.tvec[n][idx]
	tq.tvec[n][idx] = list.New()

	for e := vec.Front(); e != nil; e = e.Next() {
		t := e.Value.(*timer)
		tq.addTimer(t)
	}
	return idx
}

func (tq *TimerQueue) tick(dt int64) {
	// schedule pending timers
	tq.mutex.Lock()
	pendingTimers := tq.pendingTimers
	tq.pendingTimers = list.New()
	tq.mutex.Unlock()
	for e := pendingTimers.Front(); e != nil; e = e.Next() {
		t := e.Value.(*timer)
		tq.addTimer(t)
	}

	// tick
	for ticks := dt / MIN_TICK_INTERVAL; ticks > 0; ticks-- {
		idx := tq.ticks & TVR_MASK
		if idx == 0 &&
			tq.cascade(1) == 0 &&
			tq.cascade(2) == 0 {
			tq.cascade(3)
		}

		root := tq.tvec[0][idx]
		tq.tvec[0][idx] = list.New()
		for e := root.Front(); e != nil; e = e.Next() {
			t := e.Value.(*timer)
			t.node = nil
			t.root = nil

			if t.task != nil {
				t.task.TaskCallBack(t.id)
				t.task = nil
			}
		}
		tq.ticks++
		atomic.AddInt64(&(tq.tickTime), MIN_TICK_INTERVAL)
	}
}
