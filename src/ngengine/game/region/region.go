package region

import (
	"ngengine/core/service"
	"ngengine/game/region/scene"
	"ngengine/module/object"
	"ngengine/module/timer"
)

type Region struct {
	service.BaseService
	region *scene.SceneModule
	object *object.ObjectModule
	timer  *timer.TimerModule
}

func (s *Region) Prepare(core service.CoreAPI) error {
	s.CoreAPI = core
	s.region = scene.New()
	s.object = object.New()
	s.timer = timer.New()
	return nil
}

func (s *Region) Init(opt *service.CoreOption) error {
	s.AddModule(s.region)
	s.AddModule(s.object)
	s.AddModule(s.timer)
	return nil
}

func (s *Region) Start() error {
	s.BaseService.Start()
	return nil
}
