package main

import "ngengine/core/service"

type ModuleTest struct {
	core service.CoreApi
}

func (m *ModuleTest) Name() string {
	return "ModuleTest"
}

func (m *ModuleTest) Init(core service.CoreApi) bool {
	m.core = core
	return true
}

func (m *ModuleTest) Shut() {

}

func (m *ModuleTest) OnUpdate(t *service.Time) {
	if t.FrameCount()%5000 == 0 {
		m.core.CallModule("ModuleTest2", 1, "hello world")
	}
}

func (m *ModuleTest) OnMessage(id int, args ...interface{}) {

}

type ModuleTest2 struct {
	core service.CoreApi
}

func (m *ModuleTest2) Name() string {
	return "ModuleTest2"
}

func (m *ModuleTest2) Init(core service.CoreApi) bool {
	m.core = core
	return true
}

func (m *ModuleTest2) Shut() {

}

func (m *ModuleTest2) OnUpdate(t *service.Time) {

}

func (m *ModuleTest2) OnMessage(id int, args ...interface{}) {
	m.core.LogDebug("recv message, id:", id, ", msg:", args[0].(string))
}
