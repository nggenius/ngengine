// package object
// 对象管理模块
// for example:
//
package object

import "ngengine/core/service"

type ObjectModule struct {
	core           service.CoreApi
	defaultFactory *factory // 默认对象工厂
	factorys       map[string]*factory
}

func New() *ObjectModule {
	o := &ObjectModule{}
	o.defaultFactory = newFactory()
	o.factorys = make(map[string]*factory)
	return o
}

// Name 模块名
func (o *ObjectModule) Name() string {
	return "Object"
}

// Init 模块初始化
func (o *ObjectModule) Init(core service.CoreApi) bool {
	o.core = core
	return true
}

// Shut 模块关闭
func (o *ObjectModule) Shut() {

}

// OnUpdate 模块Update
func (o *ObjectModule) OnUpdate(t *service.Time) {
}

// OnMessage 模块消息
func (o *ObjectModule) OnMessage(id int, args ...interface{}) {
}

// RegisterObject 注册一个对象
func (o *ObjectModule) RegisterObject(typ string, obj interface{}) {

}

// 创建
func (o *ObjectModule) Create(typ string) (Object, error) {
	return o.defaultFactory.Create(typ)
}
