package space

import "ngengine/common/fsm"

type CreateRegion struct {
	fsm.Default
	owner *SpaceManage
	Idle  int32
}

func newCreateRegion(o *SpaceManage) *CreateRegion {
	s := new(CreateRegion)
	s.owner = o
	return s
}

func (s *CreateRegion) Enter() {
	if err := s.owner.createAllRegions(); err != nil {
		s.owner.ctx.Core.LogErr(err)
	}
}

func (s *CreateRegion) Handle(event int, param interface{}) string {
	switch event {
	case ETIMER:
		s.Idle++
	default:
	}
	return ""
}
