package scene

import (
	"ngengine/core/service"
	"ngengine/module/object"
)

type SceneModule struct {
	service.Module
	creater *RegionCreate
	object  *object.ObjectModule
	scenes  *Scenes
}

func New() *SceneModule {
	m := new(SceneModule)
	m.creater = NewRegionCreate(m)
	m.scenes = NewScenes(m)
	return m
}

func (m *SceneModule) Init() bool {
	f := m.Core.Module("Object")
	if f == nil {
		panic("need object module")
	}
	m.object = f.(*object.ObjectModule)

	m.object.Register("GameScene", new(GameSceneCreater))
	m.Core.RegisterRemote("Region", m.creater)
	return true
}

func (m *SceneModule) Name() string {
	return "Scene"
}
