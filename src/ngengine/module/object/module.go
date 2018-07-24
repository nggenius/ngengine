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
	service.Module
	core           service.CoreAPI
	defaultFactory *Factory // 默认对象工厂
	factorys       map[int]*Factory
	regs           map[string]ObjectCreate
	entitydelegate map[string]*EventDelegate
	sync           *SyncObject
	router         *ObjectRouter
}

func New() *ObjectModule {
	o := &ObjectModule{}
	o.defaultFactory = newFactory(o, share.OBJECT_TYPE_NONE)
	o.factorys = make(map[int]*Factory)
	o.regs = make(map[string]ObjectCreate)
	o.entitydelegate = make(map[string]*EventDelegate)
	o.sync = &SyncObject{o}
	o.router = NewObjectRouter(o)
	return o
}

// Name 模块名
func (o *ObjectModule) Name() string {
	return "Object"
}

// Init 模块初始化
func (o *ObjectModule) Init(core service.CoreAPI) bool {
	o.core = core
	o.core.RegisterRemote("object", o.sync)
	o.core.RegisterRemote("ObjectRouter", o.router)
	return true
}

// OnUpdate 模块Update
func (o *ObjectModule) OnUpdate(t *service.Time) {
	o.defaultFactory.ClearDelete()
	for _, f := range o.factorys {
		f.ClearDelete()
	}
}

// 获取日志接口
func (o *ObjectModule) Logger() logger.Logger {
	return o.core
}

// 增加一个对象工厂
func (o *ObjectModule) AddFactory(identity int) error {
	if _, has := o.factorys[identity]; has {
		return fmt.Errorf("factory already created")
	}

	o.factorys[identity] = newFactory(o, identity)

	return nil
}

// 获取工厂
func (o *ObjectModule) Factory(identity int) *Factory {
	if f, has := o.factorys[identity]; has {
		return f
	}
	return nil
}

func (o *ObjectModule) FactoryCreate(identity int, typ string) (interface{}, error) {
	if f, has := o.factorys[identity]; has {
		return f.Create(typ)
	}

	return nil, fmt.Errorf("factory %d not found", identity)
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
	if mb.Identity() == o.defaultFactory.identity {
		return o.defaultFactory.FindObject(mb)
	}

	if f, ok := o.factorys[mb.Identity()]; ok {
		return f.FindObject(mb)
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
	o.router.Register(name, oc.Create())
}
