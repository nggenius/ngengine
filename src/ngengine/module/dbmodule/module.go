//Package dbmodule 用来注册数据库操作相关函数
// for example:
//
package dbmodule

import (
	"ngengine/core/rpc"
	"ngengine/core/service"
	"ngengine/logger"
)

// DbModule 模块结构体
type DbModule struct {
	*rpc.Thread
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
		DbModule: DbModule{rpc.NewThread("Database", 10, 10, logger.New("dbLog", 1)), core}})
	return true
}

// NewJob 重写Thread里面的NewJob(为了保证同一个对象多次获取数据的时候数据获取的先后顺序)
func (m *DbModule) NewJob(r *rpc.RpcCall) bool {
	m.Queue[int(r.GetSrc().Id)%m.Pools] <- r
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
