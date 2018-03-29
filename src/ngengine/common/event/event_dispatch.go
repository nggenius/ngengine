package event

import (
	"reflect"
)

type Dispatcher interface {
	ClearEvent()
	AddListener(event string, f callback)
	RemoveListener(event string, f callback)
	DispatchEvent(event string, args ...interface{})
}

type callback func(event string, args ...interface{})

type DelegateList []callback

// 事件代理，按优先级排序
type EventDispatch struct {
	event map[string]DelegateList
}

func NewEventDispatch() *EventDispatch {
	e := &EventDispatch{}
	e.event = make(map[string]DelegateList)
	return e
}

// 初始化事件
func (e *EventDispatch) ClearEvent() {
	e.event = make(map[string]DelegateList)
}

// 执行指定事件的回调,并带返回值
func (e *EventDispatch) DispatchEvent(event string, args ...interface{}) {
	if l, has := e.event[event]; has {
		for _, f := range l {
			if f == nil {
				continue
			}
			f(event, args...)
		}
	}
}

// 增加事件监听
func (e *EventDispatch) AddListener(event string, f callback) {
	if _, has := e.event[event]; !has {
		e.event[event] = make(DelegateList, 0, 4)
	}

	e.event[event] = append(e.event[event], f)

}

// 移除事件监听
func (e *EventDispatch) RemoveListener(event string, f callback) {
	del := make([]int, 0, 2)
	if l, has := e.event[event]; has {
		for k, f1 := range l {
			if f1 == nil {
				continue
			}

			sf1 := reflect.ValueOf(f1)
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
