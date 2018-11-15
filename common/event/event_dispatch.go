package event

type Dispatcher interface {
	ClearEvent()
	AddListener(event string, f callback) *EventListener
	RemoveListener(event string, l *EventListener)
	DispatchEvent(event string, args ...interface{})
}

type EventListener struct {
	handler callback
}

func newEventListener(c callback) *EventListener {
	l := new(EventListener)
	l.handler = c
	return l
}

type callback func(event string, args ...interface{})

type delegate []*EventListener

// 事件代理，按优先级排序
type EventDispatch struct {
	delegate map[string]delegate
}

func NewEventDispatch() *EventDispatch {
	e := &EventDispatch{}
	e.delegate = make(map[string]delegate)
	return e
}

// 初始化事件
func (e *EventDispatch) ClearEvent() {
	e.delegate = make(map[string]delegate)
}

// 执行指定事件的回调,并带返回值
func (e *EventDispatch) DispatchEvent(event string, args ...interface{}) {
	if l, has := e.delegate[event]; has {
		for _, f := range l {
			if f == nil {
				continue
			}
			f.handler(event, args...)
		}
	}
}

// 增加事件监听
func (e *EventDispatch) AddListener(event string, f callback) *EventListener {
	if _, has := e.delegate[event]; !has {
		e.delegate[event] = make(delegate, 0, 4)
	}
	l := newEventListener(f)
	e.delegate[event] = append(e.delegate[event], l)
	return l
}

// 移除事件监听
func (e *EventDispatch) RemoveListener(event string, l *EventListener) {
	del := make([]int, 0, 2)
	if d, has := e.delegate[event]; has {
		for k, f1 := range d {
			if f1 == nil {
				continue
			}
			if f1 == l {
				d[k] = nil
				del = append(del, k)
			}
		}
	}
	for _, d := range del {
		copy(e.delegate[event][d:], e.delegate[event][d+1:])
	}

	size := len(e.delegate[event])
	e.delegate[event] = e.delegate[event][:size-len(del)]
}
