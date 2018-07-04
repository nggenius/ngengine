package session

import "ngengine/common/fsm"

type idlestate struct {
	fsm.Default
	owner *Session
	Idle  int32
}

func (s *idlestate) Handle(event int, param interface{}) string {
	switch event {
	case LOGIN:
		token := param.(string)
		if s.owner.ValidToken(token) {
			if !s.owner.QueryRoleInfo() {
				s.owner.Error(-1)
				return ""
			}
			return SLOGGED
		}
		// 验证失败直接踢下线
		s.owner.Break()
		return ""
	case TIMER:
		s.Idle++
		if s.Idle > 60 {
			s.owner.Break()
			return ""
		}
	case BREAK:
		s.owner.DestroySelf()
		return fsm.STOP
	default:
		s.owner.ctx.core.LogWarnf("idle state receive error event(%d)", event)
	}
	return ""
}
