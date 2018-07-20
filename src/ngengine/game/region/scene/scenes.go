package scene

import (
	"ngengine/core/rpc"
	"ngengine/share"
)

type Scenes struct {
	ctx    *SceneModule
	scenes map[int]*GameScene
}

func NewScenes(ctx *SceneModule) *Scenes {
	s := new(Scenes)
	s.ctx = ctx
	s.scenes = make(map[int]*GameScene)
	return s
}

func (s *Scenes) CreateScene(r share.Region) (rpc.Mailbox, error) {
	err := s.ctx.object.AddFactory(r.Id)
	if err != nil {
		return rpc.NullMailbox, err
	}

	f := s.ctx.object.Factory(r.Id)
	scene, err := f.Create("GameScene")
	if err != nil {
		return rpc.NullMailbox, err
	}

	gamescene := scene.(*GameScene)
	gamescene.factory = f
	gamescene.region = r

	s.scenes[r.Id] = gamescene

	return gamescene.Scene.ObjId(), nil
}
