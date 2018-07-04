package session

import (
	"ngengine/common/fsm"
	"ngengine/share"
)

type idlestate struct {
	fsm.Default
	owner *Session
	Idle  int32
}

func (s *idlestate) Handle(event int, param interface{}) string {
	switch event {
	case ELOGIN:
		token := param.(string)
		if s.owner.ValidToken(token) {
			// TODO: 这里要进行排队检查
			if !s.owner.QueryRoleInfo() {
				s.owner.Error(share.ERR_SYSTEM_ERROR)
				return ""
			}
			return SLOGGED
		}
		// 验证失败直接踢下线
		s.owner.Break()
		return ""
	case ETIMER:
		s.Idle++
		if s.Idle > 60 {
			s.owner.Break()
			return ""
		}
	case EBREAK:
		s.owner.DestroySelf()
		return fsm.STOP
	default:
		s.owner.ctx.core.LogWarnf("idle state receive error event(%d)", event)
	}
	return ""
}
