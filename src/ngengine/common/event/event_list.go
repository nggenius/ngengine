package event

var (
	eventCache = make(chan *Event, 512) //缓存
)

type EventArgs map[string]interface{}
type Event struct {
	Typ  string
	Args EventArgs
}

func (e *Event) Free() {
	for k, _ := range e.Args {
		delete(e.Args, k)
	}
	select {
	case eventCache <- e:
	default:
	}
}

// 异步事件队列
type AsyncEvent struct {
	priorQueue chan *Event //高优先级队列
	eventQueue chan *Event //普通队列
}

// 生成一个异步事件
func (e *AsyncEvent) Fire(t string, args EventArgs, priority bool) {
	var event *Event
	select {
	case event = <-eventCache:
	default:
		event = new(Event)
	}

	event.Typ = t
	event.Args = args
	if priority {
		e.priorQueue <- event
		return
	}
	e.eventQueue <- event
}

// 捕获一个异步事件
func (e *AsyncEvent) Capture() *Event {
	var event *Event
	select {
	case event = <-e.priorQueue:
	default:
		select {
		case event = <-e.priorQueue:
		case event = <-e.eventQueue:
		default:
		}
	}
	return event
}

func (e *AsyncEvent) getEvent() *Event {
	var event *Event
	select {
	case event = <-eventCache:
	default:
		event = &Event{}
	}
	return event
}

func NewAsyncEvent() *AsyncEvent {
	el := &AsyncEvent{}
	el.priorQueue = make(chan *Event, 32)
	el.eventQueue = make(chan *Event, 512)
	return el
}
