package service

import (
	"errors"
)

// 模块回调接口
type ModuleHandler interface {
	Name() string
	Init(core CoreApi) bool
	Shut()
	OnUpdate(t *Time)
	OnMessage(id int, args ...interface{})
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
