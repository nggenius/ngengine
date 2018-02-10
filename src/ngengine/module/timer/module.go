package timer

import (
	"ngengine/core/service"
)

type TimerModule struct {
	core service.CoreApi
	t    *timerManager
}

func New() *TimerModule {
	m := &TimerModule{}
	return m
}

func (m *TimerModule) Name() string {
	return "TimerModule"
}

func (m *TimerModule) Init(core service.CoreApi) bool {
	m.core = core
	m.t = newManager()
	return true
}

func (m *TimerModule) Shut() {

}

func (m *TimerModule) OnUpdate(t *service.Time) {
	m.t.run()
}

func (m *TimerModule) OnMessage(id int, args ...interface{}) {

}

func (m *TimerModule) AddTimer(delta int64, args interface{}, cb timerCallBack) (id int64) {
	return m.t.addTimer(delta, args, cb)
}

func (m *TimerModule) AddCountTimer(amount int, delta int64, args interface{}, cb timerCallBack) (id int64) {
	return m.t.addCountTimer(amount, delta, args, cb)
}

func (m *TimerModule) RemoveTimer(id int64) bool {
	return m.t.removeTimer(id)
}

func (m *TimerModule) FindTimer(id int64) (bool, int) {
	return m.t.findTimer(id)
}

func (m *TimerModule) GetTimerDelta(id int64) int64 {
	return m.t.getTimerDelta(id)
}
