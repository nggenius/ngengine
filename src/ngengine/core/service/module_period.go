package service

import (
	"container/list"
	"time"
)

type periodcb func(time.Duration)

type periodInfo struct {
	id      int
	handler periodcb
}

type Period struct {
	p    map[time.Duration]*list.List
	tick map[time.Duration]time.Time
	pid  int
}

func newPeriod() *Period {
	s := new(Period)
	s.p = make(map[time.Duration]*list.List)
	s.tick = make(map[time.Duration]time.Time)
	return s
}

// AddPeriod 增加一个周期
func (p *Period) AddPeriod(t time.Duration) {
	if _, ok := p.p[t]; ok {
		return
	}

	p.p[t] = list.New()
	p.tick[t] = time.Now()
}

// RemovePeriod 移除周期
func (p *Period) RemovePeriod(t time.Duration) {
	delete(p.p, t)
	delete(p.tick, t)
}

// Add 注册周期回调
func (p *Period) AddCallback(t time.Duration, c periodcb) int {
	if l, ok := p.p[t]; ok {
		p.pid++
		i := &periodInfo{id: p.pid, handler: c}
		l.PushBack(i)
		return p.pid
	}

	return 0
}

// Remove 移除回调
func (p *Period) RemoveCallback(t time.Duration, pid int) {
	if l, ok := p.p[t]; ok {
		for e := l.Front(); e != nil; e = e.Next() {
			if e.Value.(*periodInfo).id == pid {
				l.Remove(e)
				return
			}
		}
	}
}

func (p *Period) Update(t *Time) {
	n := time.Now()
	for k, v := range p.tick {
		if d := n.Sub(v); d >= k {
			for e := p.p[k].Front(); e != nil; e = e.Next() {
				e.Value.(*periodInfo).handler(d)
			}
			p.tick[k] = n
		}
	}
}
