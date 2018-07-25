package session

import "ngengine/common/fsm"

type offline struct {
	fsm.Default
	owner *Session
	Idle  int32
}

func (o *offline) Handle(event int, param interface{}) string {
	switch event {
	case EBREAK:
		o.owner.DestroySelf()
		return fsm.STOP
	case ETIMER:
	default:
		o.owner.ctx.Core.LogWarnf("offline state receive error event(%d)", event)
	}
	return ""
}
