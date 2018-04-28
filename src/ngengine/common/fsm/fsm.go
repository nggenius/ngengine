package fsm

import (
	"fmt"
)

const (
	STOP = "-"
)

type State interface {
	Enter()
	Handle(event int, param interface{}) string
	Exit()
}

type Default struct{}

func (d *Default) Enter() {

}

func (d *Default) Handle(event int, param interface{}) string {
	return STOP
}

func (d *Default) Exit() {

}

// 有限状态机
type FSM struct {
	state    map[string]State
	def      string
	current  State
	curstate string
}

func NewFSM() *FSM {
	f := &FSM{}
	f.state = make(map[string]State)
	return f
}

// 注册状态
func (f *FSM) Register(name string, state State) {
	if _, dup := f.state[name]; dup {
		panic("register state twice")
	}

	f.state[name] = state
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
