package service

import (
	"errors"
)

// 模块回调接口
type ModuleHandler interface {
	Name() string
	Init(core CoreAPI) bool
	Start()
	Shut()
	OnUpdate(t *Time)
	OnMessage(id int, args ...interface{})
}

type Module struct {
}

// Start 模块启动
func (m *Module) Start() {
}

// Shut 模块关闭
func (m *Module) Shut() {
}

// OnUpdate 模块Update
func (m *Module) OnUpdate(t *Time) {
}

// OnMessage 模块消息
func (m *Module) OnMessage(id int, args ...interface{}) {

}

// 模块集合
type modules struct {
	modules map[string]ModuleHandler
}

func NewModules() *modules {
	m := &modules{}
	m.modules = make(map[string]ModuleHandler)
	return m
}

// 增加一个模块
func (ms *modules) AddModule(m ModuleHandler) error {
	if _, dup := ms.modules[m.Name()]; dup {
		return errors.New("module is dup")
	}

	ms.modules[m.Name()] = m
	return nil
}

// 获取模块
func (ms *modules) Module(name string) ModuleHandler {
	if m, has := ms.modules[name]; has {
		return m
	}

	return nil
}
