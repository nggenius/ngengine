package service

import (
	"errors"
)

// 模块回调接口
type ModuleHandler interface {
	SetCore(c CoreAPI)
	Name() string
	Prepare()
	Init() bool
	Start()
	Shut()
	Update(t *Time)
	OnUpdate(t *Time)
	OnMessage(id int, args ...interface{})
}

type Module struct {
	*Period
	Core CoreAPI
}

func (m *Module) Prepare() {
	m.Period = newPeriod()
}

// 设置核心
func (m *Module) SetCore(c CoreAPI) {
	m.Core = c
}

// Start 模块启动
func (m *Module) Start() {
}

// Shut 模块关闭
func (m *Module) Shut() {
}

// Update 周期更新
func (m *Module) Update(t *Time) {
	m.Period.Update(t)
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
