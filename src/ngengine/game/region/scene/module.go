package scene

import (
	"ngengine/core/service"
	"ngengine/module/object"
	"ngengine/module/timer"
	"time"
)

type SceneModule struct {
	service.Module
	creater *RegionCreate
	object  *object.ObjectModule
	timer   *timer.TimerModule
	scenes  *Scenes
}

func New() *SceneModule {
	m := new(SceneModule)
	m.creater = NewRegionCreate(m)
	m.scenes = NewScenes(m)
	return m
}

func (m *SceneModule) Init() bool {
	m.object = m.Core.MustModule("Object").(*object.ObjectModule)
	m.timer = m.Core.MustModule("Timer").(*timer.TimerModule)
	m.object.Register(new(GameScene))
	m.Core.RegisterRemote("Region", m.creater)
	m.AddPeriod(time.Second)
	m.AddCallback(time.Second, m.PerSecondCheck)
	return true
}

func (m *SceneModule) Name() string {
	return "Scene"
}

func (s *SceneModule) PerSecondCheck(d time.Duration) {
	s.scenes.UpdateAllScene(d)
}

func (s *SceneModule) OnUpdate(t *service.Time) {
	s.Module.Update(t)
}
