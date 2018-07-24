package scene

import (
	"ngengine/core/rpc"
	"ngengine/game/gameobject"
	"ngengine/game/gameobject/entity"
	"ngengine/module/object"
	"ngengine/protocol"
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

func (s *GameScene) RegisterCallback(svr rpc.Servicer) {
	svr.RegisterCallback("Test", s.Test)
}

//srv.RegisterCallback("FunctionName", s.FunctionName)
func (s *GameScene) Test(src rpc.Mailbox, dest rpc.Mailbox, msg *protocol.Message) (int32, *protocol.Message) {
	var info string
	err := protocol.ParseArgs(msg, &info)
	if err != nil {
		
	}
	return 0, nil
}

type GameSceneCreater struct {
}

func (g *GameSceneCreater) Create() interface{} {
	s := NewGameScene()
	return s
}
