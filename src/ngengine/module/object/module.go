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
	regs           map[string]ObjectCreate
	entitydelegate map[string]*EventDelegate
}

func New() *ObjectModule {
	o := &ObjectModule{}
	o.defaultFactory = newFactory(o, share.OBJECT_TYPE_OBJECT)
	o.factorys = make(map[string]*Factory)
	o.regs = make(map[string]ObjectCreate)
	o.entitydelegate = make(map[string]*EventDelegate)
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

// Start 模块启动
func (o *ObjectModule) Start() {

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

// 注册对象事件回调
func (o *ObjectModule) AddEventCallback(typ, event string, f callback, priority int) error {
	if delegate, has := o.entitydelegate[typ]; has {
		delegate.AddListener(event, f, priority)
		return nil
	}

	return fmt.Errorf("object %s type not register ", typ)
}

// 移除事件回调
func (o *ObjectModule) RemoveEventCallback(typ, event string, f callback) error {
	if delegate, has := o.entitydelegate[typ]; has {
		delegate.RemoveListener(event, f)
		return nil
	}

	return fmt.Errorf("object %s type not register ", typ)
}

// 注册对象
func (o *ObjectModule) Register(name string, oc ObjectCreate) {

	if oc == nil {
		panic("object: Register object is nil")
	}
	if _, dup := o.regs[name]; dup {
		panic("object: Register called twice for object " + name)
	}

	o.regs[name] = oc
	o.entitydelegate[name] = NewEventDelegate()
}
