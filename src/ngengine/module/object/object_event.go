package object

import (
	"ngengine/core/rpc"
	"sort"
)

const (
	PRIORITY_LOWEST  = -1024
	PRIORITY_LOWER   = -512
	PRIORITY_NORMAL  = 0
	PRIORITY_HIGH    = 512
	PRIORITY_HIGHEST = 1024
)

type Delegate interface {
	Invoke(event string, self, sender rpc.Mailbox, args ...interface{}) int
	InvokeNoReturn(event string, self, sender rpc.Mailbox, args ...interface{})
}

type callback func(self, sender rpc.Mailbox, args ...interface{}) int

type PriorityDelegate struct {
	sort int
	cb   callback
}

type DelegateList []*PriorityDelegate

func (l DelegateList) Len() int {
	return len(l)
}

func (l DelegateList) Less(i, j int) bool {
	return l[i].sort < l[j].sort
}

func (l DelegateList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// 事件代理，按优先级排序
type EventDelegate struct {
	event map[string]DelegateList
}

func NewEventDelegate() *EventDelegate {
	e := &EventDelegate{}
	e.event = make(map[string]DelegateList)
	return e
}

// 执行指定事件的回调,并带返回值
func (e *EventDelegate) Invoke(event string, sender, target rpc.Mailbox, args ...interface{}) int {
	ret := 0
	if l, has := e.event[event]; has {
		for _, v := range l {
			if v == nil || v.cb == nil {
				continue
			}
			ret = v.cb(target, sender, args...)
			if ret != 0 { // 不等于0，则终止调用
				break
			}
		}
	}
	return ret
}

// 执行指定事件的回调
func (e *EventDelegate) InvokeNoReturn(event string, self, sender rpc.Mailbox, args ...interface{}) {
	if l, has := e.event[event]; has {
		for _, v := range l {
			if v == nil || v.cb == nil {
				continue
			}
			v.cb(self, sender, args...)
		}
	}
}

// 增加事件监听
func (e *EventDelegate) AddListener(event string, f callback, priority int) {
	if _, has := e.event[event]; !has {
		e.event[event] = make(DelegateList, 0, 4)
	}

	e.event[event] = append(e.event[event], &PriorityDelegate{priority, f})
	sort.Sort(sort.Reverse(e.event[event]))
}
