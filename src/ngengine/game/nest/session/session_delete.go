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
		errcode := param.(int32)
		if errcode != 0 {
			s.owner.Error(errcode)
			return SLOGGED
		}
		s.owner.QueryRoleInfo()
		return SLOGGED
	case EBREAK:
		s.owner.DestroySelf()
		return fsm.STOP
	default:
		s.owner.ctx.Core.LogWarnf("deleting state receive error event(%d)", event)
	}
	return ""
}
