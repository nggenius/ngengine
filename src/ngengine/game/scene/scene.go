package scene

import "ngengine/core/service"

type Scene struct {
	service.BaseService
}

func (s *Scene) Prepare(core service.CoreAPI) error {
	return nil
}

func (s *Scene) Init(opt *service.CoreOption) error {
	return nil
}

func (s *Scene) Start() error {
	s.BaseService.Start()
	return nil
}
