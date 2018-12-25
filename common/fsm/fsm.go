package fsm

import (
	"fmt"
)

const (
	STOP   = "-"
	ETIMER = -1
)

type StateRegister interface {
	AddHandle(event int, handle StateHandle)
}

type State interface {
	Init(StateRegister)
	Enter()
	OnHandle(event int, param interface{}) string
	OnTimer() string
	Exit()
}

type Default struct {
}

func (d *Default) Init(StateRegister)                           {}
func (d *Default) Enter()                                       {}
func (d *Default) OnHandle(event int, param interface{}) string { return STOP }
func (d *Default) OnTimer() string                              { return "" }
func (d *Default) Exit()                                        {}

type StateHandle func(event int, param interface{}) string
type StateWrapper struct {
	s  State
	sf map[int]StateHandle
}

func NewState(s State) *StateWrapper {
	d := new(StateWrapper)
	d.s = s
	d.sf = make(map[int]StateHandle)
	return d
}

func (d *StateWrapper) Init() {
	d.s.Init(d)
}

func (d *StateWrapper) AddHandle(event int, handle StateHandle) {
	d.sf[event] = handle
}

func (d *StateWrapper) Enter() {
	d.s.Enter()
}

func (d *StateWrapper) Handle(event int, param interface{}) string {
	if event == ETIMER {
		return d.s.OnTimer()
	}
	if h, ok := d.sf[event]; ok {
		return h(event, param)
	}
	return d.s.OnHandle(event, param)
}

func (d *StateWrapper) Exit() {
	d.s.Exit()
}

// 有限状态机
type FSM struct {
	state    map[string]*StateWrapper
	def      string
	current  *StateWrapper
	curstate string
}

func NewFSM() *FSM {
	f := &FSM{}
	f.state = make(map[string]*StateWrapper)
	return f
}

// 注册状态
func (f *FSM) Register(name string, state State) {
	if _, dup := f.state[name]; dup {
		panic("register state twice")
	}
	s := NewState(state)
	s.Init()
	f.state[name] = s
}

// 启动状态机, state是初始状态
func (f *FSM) Start(state string) error {
	if _, exist := f.state[state]; !exist {
		panic("state not found")
	}

	f.def = state

	s, exist := f.state[f.def]
	if !exist {
		return fmt.Errorf("state not found, %s", f.def)
	}

	f.curstate = f.def
	f.current = s
	f.current.Enter()
	return nil
}

// 超时
func (f *FSM) Timeout() (bool, error) {
	return f.Dispatch(ETIMER, nil)
}

// 派发事件
func (f *FSM) Dispatch(event int, param interface{}) (bool, error) {
	if f.current == nil {
		return true, fmt.Errorf("current state is nil")
	}

	ret := f.current.Handle(event, param)
	if ret == "" {
		return false, nil
	}

	if ret == STOP {
		f.current.Exit()
		f.current = nil
		f.curstate = ""
		return true, nil
	}

	next, exist := f.state[ret]
	if !exist {
		f.current = nil
		f.curstate = ""
		return true, fmt.Errorf("state is nil, %s", ret)
	}

	f.current.Exit()
	f.curstate = ret
	f.current = next
	f.current.Enter()
	return false, nil
}
