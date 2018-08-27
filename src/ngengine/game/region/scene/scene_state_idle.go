package scene

import "ngengine/common/fsm"

type Idle struct {
	fsm.Default
	owner *GameScene
	Idle  int32
}

func newIdle(o *GameScene) *Idle {
	s := new(Idle)
	s.owner = o
	return s
}

func (s *Idle) Handle(event int, param interface{}) string {
	switch event {
	case ETIMER:
		s.Idle++
	default:
	}
	return ""
}
