package space

import "ngengine/common/fsm"

type Idle struct {
	fsm.Default
	owner *SpaceManage
	Idle  int32
}

func newIdle(o *SpaceManage) *Idle {
	s := new(Idle)
	s.owner = o
	return s
}

func (s *Idle) Handle(event int, param interface{}) string {
	switch event {
	case ETIMER:
		s.Idle++
	case EREGION_RESP:
		if s.owner.checkAllRegion() {
			return SCREATE
		}
	default:
	}
	return ""
}
