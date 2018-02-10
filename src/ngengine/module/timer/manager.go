package timer

type timerCallBack func(id int64, count int, args interface{})

// 定时器类型
const (
	NONE = iota
	REPEAT
	COUNT
)

type timerTask struct {
	id      int64 // 任务id
	kind    int   // 类型
	delta   int64 // 心跳时间
	count   int   // 次数
	amount  int   // 总数
	args    interface{}
	cb      timerCallBack
	timer   *timerQueue
	manager *timerManager
}

func (t *timerTask) taskCallBack(id int64) {
	t.count++
	if t.kind == REPEAT {
		t.timer.schedule(t.delta, t)
		t.cb(t.id, t.count, t.args)
	} else if t.kind == COUNT {
		if t.count < t.amount {
			t.timer.schedule(t.delta, t)
			t.cb(t.id, t.count, t.args)
		} else if t.count == t.amount {
			t.cb(t.id, t.count, t.args)
		}
	}
}

type timerManager struct {
	genID   int64
	taskMap map[int64]*timerTask
	timer   *timerQueue
}

func newManager() *timerManager {
	manager := &timerManager{
		genID:   0,
		taskMap: make(map[int64]*timerTask),
		timer:   newTimerQueue(),
	}
	return manager
}

func (t *timerManager) run() {
	t.timer.run()
}

func (t *timerManager) generateID() int64 {
	t.genID++
	return t.genID
}

func (t *timerManager) addTimer(delta int64, args interface{}, cb timerCallBack) (id int64) {
	task := &timerTask{
		id:      t.generateID(),
		kind:    REPEAT,
		delta:   delta,
		count:   0,
		amount:  0,
		cb:      cb,
		args:    args,
		timer:   t.timer,
		manager: t,
	}
	t.taskMap[task.id] = task
	t.timer.schedule(delta, task)
	return task.id
}

func (t *timerManager) addCountTimer(amount int, delta int64, args interface{}, cb timerCallBack) (id int64) {
	task := &timerTask{
		id:      t.generateID(),
		kind:    COUNT,
		delta:   delta,
		count:   0,
		amount:  amount,
		cb:      cb,
		args:    args,
		timer:   t.timer,
		manager: t,
	}
	t.taskMap[task.id] = task
	t.timer.schedule(delta, task)
	return task.id
}

func (t *timerManager) removeTimer(id int64) bool {
	_, ok := t.taskMap[id]
	if ok {
		t.taskMap[id] = nil
		delete(t.taskMap, id)
		return true
	}
	return false
}

func (t *timerManager) findTimer(id int64) (bool, int) {
	task, ok := t.taskMap[id]
	if ok {
		return true, task.kind
	}
	return false, NONE
}

func (t *timerManager) getTimerDelta(id int64) int64 {
	task, ok := t.taskMap[id]
	if ok {
		return task.delta
	}
	return 0
}
