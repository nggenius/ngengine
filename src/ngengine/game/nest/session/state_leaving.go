package session

import (
	"ngengine/common/fsm"
)

type leaving struct {
	fsm.Default
	owner *Session
	Idle  int32
}

func (l *leaving) Handle(event int, param interface{}) string {
	switch event {
	case ESTORED:
		l.owner.Break()
		l.owner.DestroySelf()
	case ETIMER:
		l.Idle++
		if l.Idle > 60 {
			l.owner.Break()
			l.owner.DestroySelf()
		}
	default:
		l.owner.ctx.core.LogWarnf("leaving state receive error event(%d)", event)
	}
	return ""
}
