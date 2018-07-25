package session

import (
	"ngengine/common/fsm"
	"ngengine/share"
)

type createrole struct {
	fsm.Default
	owner *Session
	Idle  int32
}

func (c *createrole) Handle(event int, param interface{}) string {
	switch event {
	case ECREATED:
		errcode := param.(int32)
		if errcode != 0 {
			c.owner.Error(errcode)
			return SLOGGED
		}
		if !c.owner.QueryRoleInfo() {
			c.owner.Error(-1)
		}
		return SLOGGED
	case ETIMER:
		c.Idle++
		if c.Idle > 60 {
			c.owner.Error(share.ERR_CREATE_TIMEOUT)
			return SLOGGED
		}
	case EBREAK:
		c.owner.DestroySelf()
		return fsm.STOP
	default:
		c.owner.ctx.Core.LogWarnf("create role state receive error event(%d)", event)
	}
	return ""
}
