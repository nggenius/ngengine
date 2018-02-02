package timer

import (
	"ngengine/core/service"
)

type TimerModule struct {
	core service.CoreApi
	t    *TimerManager
}

func (m *TimerModule) Name() string {
	return "TimerModule"
}

func (m *TimerModule) Init(core service.CoreApi) bool {
	m.core = core
	m.t = NewManager()
	return true
}

func (m *TimerModule) Shut() {

}

func (m *TimerModule) OnUpdate(t *service.Time) {
	m.t.Run()
}

func (m *TimerModule) OnMessage(id int, args ...interface{}) {

}

func (m *TimerModule) AddTimer(delta int64, args interface{}, cb TimerCallBack) (id int64) {
	return m.t.AddTimer(delta, args, cb)
}

func (m *TimerModule) AddCountTimer(amount int, delta int64, args interface{}, cb TimerCallBack) (id int64) {
	return m.t.AddCountTimer(amount, delta, args, cb)
}

func (m *TimerModule) RemoveTimer(id int64) bool {
	return m.t.RemoveTimer(id)
}

func (m *TimerModule) FindTimer(id int64) (bool, int) {
	return m.t.FindTimer(id)
}

func (m *TimerModule) GetTimerDelta(id int64) int64 {
	return m.t.GetTimerDelta(id)
}
