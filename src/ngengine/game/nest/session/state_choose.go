package session

import (
	"ngengine/common/fsm"
	"ngengine/game/gameobject"
	"ngengine/share"
	"os"

	"github.com/davecgh/go-spew/spew"
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

		f, err := os.OpenFile("dump.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err == nil {
			spew.Fdump(f, player)
			f.Close()
		}

		c.owner.ctx.Core.LogDebug("enter game")
		c.owner.SetGameObject(player)
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
