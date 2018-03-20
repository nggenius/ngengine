package object

import (
	"reflect"
	"sort"
)

type callback func(self Object, sender Object, args ...interface{}) int

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
func (e *EventDelegate) Invoke(event string, self, sender Object, args ...interface{}) int {
	ret := 0
	if l, has := e.event[event]; has {
		for _, v := range l {
			if v == nil || v.cb == nil {
				continue
			}
			ret = v.cb(self, sender, args...)
			if ret != 0 { // 不等于0，则终止调用
				break
			}
		}
	}
	return ret
}

// 执行指定事件的回调
func (e *EventDelegate) InvokeNoReturn(event string, self, sender Object, args ...interface{}) {
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

// 移除事件监听
func (e *EventDelegate) RemoveListener(event string, f callback) {
	del := make([]int, 0, 2)
	if l, has := e.event[event]; has {
		for k, v := range l {
			if v == nil || v.cb == nil {
				continue
			}

			sf1 := reflect.ValueOf(v.cb)
			sf2 := reflect.ValueOf(f)
			if sf1.Pointer() == sf2.Pointer() {
				l[k] = nil
				del = append(del, k)
			}
		}
	}
	for _, d := range del {
		copy(e.event[event][d:], e.event[event][d+1:])
	}

	size := len(e.event[event])
	e.event[event] = e.event[event][:size-len(del)]
}
