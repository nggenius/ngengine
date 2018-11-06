package models

import (
	"ngengine/game/gameobject"
	"ngengine/game/gameobject/entity"
)

const (
	GAME_PLAYER = "GamePlayer"
)

type GamePlayer struct {
	gameobject.RoleObject
	*entity.Player
}

func (p *GamePlayer) Ctor() {
	p.Player = entity.NewPlayer()
	p.SetSpirit(p.Player)
}

func (p *GamePlayer) EntityType() string {
	return GAME_PLAYER
}
