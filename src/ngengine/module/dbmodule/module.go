//Package dbmodule 用来注册数据库操作相关函数
// for example:
//
package dbmodule

import (
	"ngengine/core/rpc"
	"ngengine/core/service"
)

// DbModule 模块结构体
type DbModule struct {
	service.Module
	*rpc.Thread
	Core service.CoreApi
}

// New 获取一个DbModule的指针
func New() *DbModule {
	o := &DbModule{}
	o.Thread = rpc.NewThread("Database", 10, 10)
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
		ctx: m,
	})
	return true
}
