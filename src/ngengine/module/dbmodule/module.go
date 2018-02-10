package dbmodule

import (
	"ngengine/core/service"
)

type DbModule struct {
	Core service.CoreApi
}

func (m *DbModule) Name() string {
	return "DbModule"
}

func (m *DbModule) Init(core service.CoreApi) bool {
	m.Core = core
	m.Core.RegisterRemote("DBModule", &DbCallBack{
		DbModule: DbModule{core}})
	return true
}

func (m *DbModule) Shut() {

}

func (m *DbModule) OnUpdate(t *service.Time) {

}

func (m *DbModule) OnMessage(id int, args ...interface{}) {

}
