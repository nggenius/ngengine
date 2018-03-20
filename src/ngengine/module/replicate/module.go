// package replica
// 对象同步模块
package replicate

import "ngengine/core/service"

type ReplicateModule struct {
	core service.CoreApi
}

func New() *ReplicateModule {
	o := &ReplicateModule{}
	return o
}

// Name 模块名
func (o *ReplicateModule) Name() string {
	return "Replicate"
}

// Init 模块初始化
func (o *ReplicateModule) Init(core service.CoreApi) bool {
	o.core = core
	return true
}

// Shut 模块关闭
func (o *ReplicateModule) Shut() {

}

// OnUpdate 模块Update
func (o *ReplicateModule) OnUpdate(t *service.Time) {
}

// OnMessage 模块消息
func (o *ReplicateModule) OnMessage(id int, args ...interface{}) {
}

// RegisterObject 注册一个对象
func (o *ReplicateModule) RegisterObject(typ string, obj interface{}) {

}
