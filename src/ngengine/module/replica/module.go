// package replica
// 对象同步模块
package replica

import "ngengine/core/service"

type ReplicaModule struct {
	core service.CoreApi
}

func New() *ReplicaModule {
	o := &ReplicaModule{}
	return o
}

// Name 模块名
func (o *ReplicaModule) Name() string {
	return "Replica"
}

// Init 模块初始化
func (o *ReplicaModule) Init(core service.CoreApi) bool {
	o.core = core
	return true
}

// Shut 模块关闭
func (o *ReplicaModule) Shut() {

}

// OnUpdate 模块Update
func (o *ReplicaModule) OnUpdate(t *service.Time) {
}

// OnMessage 模块消息
func (o *ReplicaModule) OnMessage(id int, args ...interface{}) {
}

// RegisterObject 注册一个对象
func (o *ReplicaModule) RegisterObject(typ string, obj interface{}) {

}
