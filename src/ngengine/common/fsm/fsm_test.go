package fsm

import (
	"fmt"
	"testing"
)

type State1 struct {
}

func (s *State1) Enter() {
	fmt.Println("state1 enter")
}

func (s *State1) Handle(event int, param interface{}) string {
	fmt.Println("execute", event)
	if event == 1 {
		return "State2"
	}
	return ""

}

func (s *State1) Exit() {
	fmt.Println("state1 exit")
}

type State2 struct {
}

func (s *State2) Enter() {
	fmt.Println("state2 enter")
}

func (s *State2) Handle(event int, param interface{}) string {
	if event == 1 {
		return "State3"
	}
	return ""
}

func (s *State2) Exit() {
	fmt.Println("state2 exit")
}

type State3 struct {
}

func (s *State3) Enter() {
	fmt.Println("state3 enter")
}

func (s *State3) Handle(event int, param interface{}) string {
	if event == 1 {
		return STOP
	}
	return ""
}

func (s *State3) Exit() {
	fmt.Println("state3 exit")
}
func TestFSM(t *testing.T) {
	fsm := NewFSM()
	fsm.Register("State1", &State1{})
	fsm.Register("State2", &State2{})
	fsm.Register("State3", &State3{})
	fsm.Start("State1")
	fsm.Dispatch(2, nil)
	if fsm.curstate != "State1" {
		t.Fatalf("current state is %s, need State1", fsm.curstate)
	}
	fsm.Dispatch(1, nil)
	if fsm.curstate != "State2" {
		t.Fatalf("current state is %s, need State1", fsm.curstate)
	}
	fsm.Dispatch(2, nil)
	if fsm.curstate != "State2" {
		t.Fatalf("current state is %s, need State2", fsm.curstate)
	}
	fsm.Dispatch(1, nil)
	if fsm.curstate != "State3" {
		t.Fatalf("current state is %s, need State3", fsm.curstate)
	}
	fsm.Dispatch(2, nil)
	if fsm.curstate != "State3" {
		t.Fatalf("current state is %s, need State3", fsm.curstate)
	}
	ret, _ := fsm.Dispatch(1, nil)
	if !ret {
		t.Fatalf("current state is not stop, need true")
	}

}
