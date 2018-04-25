package fsm

import (
	"fmt"
)

const (
	STOP = "-"
)

type State interface {
	Enter()
	Execute(event int, param []interface{}) string
	Exit()
}

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

func (f *FSM) AddState(name string, state State) {
	if _, dup := f.state[name]; dup {
		panic("register state twice")
	}

	f.state[name] = state
}

func (f *FSM) SetDefaultState(state string) {
	if _, exist := f.state[state]; !exist {
		panic("state not found")
	}

	f.def = state
}

func (f *FSM) Start() error {
	s, exist := f.state[f.def]
	if !exist {
		return fmt.Errorf("state not found, %s", f.def)
	}

	f.curstate = f.def
	f.current = s
	f.current.Enter()
	return nil
}

func (f *FSM) Exec(event int, param []interface{}) (bool, error) {
	if f.current == nil {
		return true, fmt.Errorf("current state is nil")
	}

	ret := f.current.Execute(event, param)
	if ret == "" {
		return false, nil
	}

	if ret == STOP {
		f.current.Exit()
		return true, nil
	}

	next, exist := f.state[ret]
	if !exist {
		return true, fmt.Errorf("state is nil, %s", ret)
	}

	f.current.Exit()
	f.curstate = ret
	f.current = next
	f.current.Enter()
	return false, nil
}
