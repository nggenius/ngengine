package session

import (
	"ngengine/common/fsm"
	"ngengine/share"
)

type deleting struct {
	fsm.Default
	owner *Session
	Idle  int32
}

func (s *deleting) Handle(event int, param interface{}) string {
	switch event {
	case ETIMER:
		s.Idle++
		if s.Idle > 60 {
			s.owner.Error(share.ERR_TIME_OUT)
			return SLOGGED
		}
	case EDELETED:
		s.owner.QueryRoleInfo()
		return SLOGGED
	case EBREAK:
		s.owner.DestroySelf()
		return fsm.STOP
	default:
		s.owner.ctx.core.LogWarnf("deleting state receive error event(%d)", event)
	}
	return ""
}
