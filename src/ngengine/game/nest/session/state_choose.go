package session

import (
	"ngengine/common/fsm"
	"ngengine/game/gameobject"
	"ngengine/share"
)

type chooserole struct {
	fsm.Default
	owner *Session
	Idle  int32
}

func (c *chooserole) Handle(event int, param interface{}) string {
	switch event {
	case ECHOOSED:
		args := param.([2]interface{})
		errcode := args[0].(int32)
		if errcode != 0 {
			c.owner.Error(errcode)
			return SLOGGED
		}
		player := args[1].(gameobject.GameObject)
		if player == nil {
			c.owner.Error(share.ERR_CHOOSE_ROLE)
			return SLOGGED
		}

		ls, ok := player.(LandInfo)
		if !ok {
			c.owner.DestroySelf()
			c.owner.ctx.Core.LogErr("entity not define landpos ", c.owner.ctx.mainEntity)
			return ""
		}
		c.owner.ctx.Core.LogDebug("enter game")
		c.owner.SetGameObject(player)
		x, y, z, o := ls.LandPosXYZOrient()
		c.owner.SetLandInfo(ls.LandScene(), x, y, z, o)
		c.owner.FindRegion()
		return SONLINE
	case EBREAK:
		c.owner.DestroySelf()
		return fsm.STOP
	case ETIMER:
		c.Idle++
		if c.Idle > 60 {
			c.owner.Error(share.ERR_CHOOSE_TIMEOUT)
			return SLOGGED
		}
	default:
		c.owner.ctx.Core.LogWarnf("choose role state receive error event(%d)", event)
	}
	return ""
}
