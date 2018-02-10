package crontab

import (
	"ngengine/core/service"
	"time"
)

// CrontabModule 时间事件模块
type CrontabModule struct {
	core  service.CoreApi
	crtab *crontab
}

// 构造一个CrontabModule
func New() *CrontabModule {
	m := &CrontabModule{}
	return m
}

// Name 模块名
func (m *CrontabModule) Name() string {
	return "CrontabModule"
}

// Init 模块初始化
func (m *CrontabModule) Init(core service.CoreApi) bool {
	m.core = core
	m.core.LogInfo("CrontabModule is init")
	m.crtab = newCrontab()
	return true
}

// Shut 模块关闭
func (m *CrontabModule) Shut() {

}

// OnUpdate 模块Update
func (m *CrontabModule) OnUpdate(t *service.Time) {
	m.check()
}

// OnMessage 模块消息
func (m *CrontabModule) OnMessage(id int, args ...interface{}) {
}

// Check crontab插件的主调用方法
func (m *CrontabModule) check() {
	if m.crtab == nil {
		return
	}
	now := time.Now().Unix()
	if now-m.crtab.lastTime >= duration {
		m.crtab.checkTriggerEvent(time.Now())
		m.crtab.lastTime = now - int64(time.Now().Second())
	}
}

// RegistEvent crontab插件事件注册接口
// 调用来注册时间事件
func (m *CrontabModule) RegistEvent(timeStr string, cb callback, args interface{}) error {
	if m.crtab == nil {
		return ErrCrontabNotInit
	}
	evt, err := parseEventTime(timeStr)
	if err != nil {
		return err
	}

	if cb == nil {
		return ErrCbNil
	}

	if args == nil {
		return ErrArgNil
	}

	evt.cb = cb
	evt.args = args
	m.crtab.evts = append(m.crtab.evts, evt)

	return nil
}
