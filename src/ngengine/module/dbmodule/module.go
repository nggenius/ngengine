// 用来注册数据库操作相关函数
// for example:
//
package dbmodule

import (
	"ngengine/core/service"
)

type DbModule struct {
	Core service.CoreApi
}

// New 获取一个DbModule的指针
func New() *DbModule {
	o := &DbModule{}
	return o
}

// Name 模块的名字
func (m *DbModule) Name() string {
	return "DbModule"
}

// Init 模块初始化
func (m *DbModule) Init(core service.CoreApi) bool {
	m.Core = core
	m.Core.RegisterRemote("Database", &DbCallBack{
		DbModule: DbModule{core}})
	return true
}

// Shut 模块关闭
func (m *DbModule) Shut() {

}

// OnUpdate 模块Update
func (m *DbModule) OnUpdate(t *service.Time) {

}

// OnMessage 模块消息
func (m *DbModule) OnMessage(id int, args ...interface{}) {

}
