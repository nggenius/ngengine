package session

import (
	"ngengine/common/fsm"
	"ngengine/core/rpc"
	"ngengine/share"
)

type online struct {
	fsm.Default
	owner  *Session
	Idle   int32
	online bool
}

func (o *online) Handle(event int, param interface{}) string {
	switch event {
	case EFREGION:
		args := param.([2]interface{})
		errcode := args[0].(int32)
		if errcode != 0 {
			o.owner.Error(errcode)
			return SLOGGED
		}

		r := args[1].(rpc.Mailbox)
		if r == rpc.NullMailbox {
			o.owner.Error(share.ERR_REGION_NOT_FOUND)
			return SLOGGED
		}

		if err := o.owner.EnterRegion(r); err != nil {
			o.owner.Error(share.ERR_ENTER_REGION_FAILED)
			return SLOGGED
		}
	case EONLINE:
		o.online = true
		o.Idle = 0
	case EBREAK:
		//o.owner.DestroySelf()
		o.owner.ctx.Core.LogInfo("client break")
		return SLEAVING
	case ETIMER:
		o.Idle++
		if !o.online {
			if o.Idle > 60 {
				o.owner.Error(share.ERR_ENTER_REGION_FAILED)
				return SLOGGED
			}
			break
		}
	default:
		o.owner.ctx.Core.LogWarnf("online state receive error event(%d)", event)
	}
	return ""
}

func (o *online) Exit() {
	if !o.online {
		if o.owner.gameobject != nil {
			o.owner.ctx.factory.Destroy(o.owner.gameobject)
			o.owner.gameobject = nil
		}
	}
}
