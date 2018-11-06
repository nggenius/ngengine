package scene

import (
	"ngengine/common/fsm"
	"ngengine/game/gameobject"
	"ngengine/game/gameobject/entity"
	"ngengine/module/object"
	"ngengine/share"
)

const GAME_SCENE = "GameScene"

type GameScene struct {
	*entity.Scene
	gameobject.SceneObject
	factory *object.Factory
	region  share.Region
	fsm     *fsm.FSM
}

func (s *GameScene) Ctor() {
	s.Scene = entity.NewScene()
	s.fsm = initState(s)
}

func (s *GameScene) EntityType() string {
	return GAME_SCENE
}

func (s *GameScene) LoadRes(res string) bool {
	return true
}
