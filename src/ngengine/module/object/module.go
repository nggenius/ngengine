// package object
// 对象管理模块
// for example:
//
package object

import (
	"fmt"
	"ngengine/core/rpc"
	"ngengine/core/service"
	"ngengine/logger"
	"ngengine/share"
)

type ObjectModule struct {
	core           service.CoreApi
	defaultFactory *Factory // 默认对象工厂
	factorys       map[string]*Factory
}

func New() *ObjectModule {
	o := &ObjectModule{}
	o.defaultFactory = newFactory(o, share.OBJECT_TYPE_OBJECT)
	o.factorys = make(map[string]*Factory)
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
	o.defaultFactory.ClearDelete()
	for _, f := range o.factorys {
		f.ClearDelete()
	}
}

// OnMessage 模块消息
func (o *ObjectModule) OnMessage(id int, args ...interface{}) {
}

// 获取日志接口
func (o *ObjectModule) Logger() logger.Logger {
	return o.core
}

func (o *ObjectModule) FactoryCreate(factory, typ string) (interface{}, error) {
	if f, has := o.factorys[factory]; has {
		return f.Create(typ)
	}

	return nil, fmt.Errorf("factory %s not found", factory)
}

// 创建
func (o *ObjectModule) Create(typ string) (interface{}, error) {
	return o.defaultFactory.Create(typ)
}

// 销毁一个对象
func (o *ObjectModule) Destroy(object interface{}) error {
	if fo, ok := object.(FactoryObject); ok {
		f := fo.Factory()
		if f != nil {
			return f.Destroy(object)
		}
	}
	return fmt.Errorf("destroy object failed")
}

// 查找对象
func (o *ObjectModule) FindObject(mb rpc.Mailbox) (interface{}, error) {
	if mb.ObjectType() == o.defaultFactory.objType {
		return o.defaultFactory.FindObject(mb)
	}

	for _, f := range o.factorys {
		if f.objType == mb.ObjectType() {
			return f.FindObject(mb)
		}
	}

	return nil, fmt.Errorf("object %s not found", mb)

}
