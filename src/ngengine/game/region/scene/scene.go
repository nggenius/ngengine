package scene

import (
	"ngengine/game/gameobject"
	"ngengine/game/gameobject/entity"
	"ngengine/module/object"
	"ngengine/share"
)

type GameScene struct {
	*entity.Scene
	gameobject.SceneObject
	factory *object.Factory
	region  share.Region
}

func NewGameScene() *GameScene {
	s := new(GameScene)
	s.Scene = entity.NewScene()
	return s
}

func (s *GameScene) LoadRes(res string) bool {
	return true
}
