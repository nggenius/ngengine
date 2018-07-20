package session

import "ngengine/common/fsm"

type online struct {
	fsm.Default
	owner *Session
	Idle  int32
}

func (o *online) Handle(event int, param interface{}) string {
	switch event {
	case EBREAK:
		//o.owner.DestroySelf()
		return SLEAVING
	case ETIMER:
	default:
		o.owner.ctx.core.LogWarnf("online state receive error event(%d)", event)
	}
	return ""
}
