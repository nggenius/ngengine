package models

import (
	"ngengine/game/gameobject"
	"ngengine/game/gameobject/entity"
)

type GamePlayer struct {
	gameobject.RoleObject
	*entity.Player
}

func NewGamePlayer() *GamePlayer {
	var p GamePlayer
	p.Player = entity.NewPlayer()
	p.SetSpirit(p.Player)
	return &p
}

type GamePlayerCreater struct{}

func (c *GamePlayerCreater) Create() interface{} {
	p := NewGamePlayer()
	return p
}
