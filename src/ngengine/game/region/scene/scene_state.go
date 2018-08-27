package scene

import "ngengine/common/fsm"

const (
	NONE = iota
	ETIMER
)

const (
	SIDLE = "idle"
)

func initState(s *GameScene) *fsm.FSM {
	fsm := fsm.NewFSM()
	fsm.Register(SIDLE, newIdle(s))
	fsm.Start(SIDLE)
	return fsm
}
