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
	"reflect"
)

const (
	GLOBAL_EVENT = "g_event"

	GLOBAL_ADD_DUMY = "g_add_dumy"
)

type ReplicateCB func(params interface{}, err error)

type create struct {
	obj interface{}
	v   reflect.Value
	t   reflect.Type
}

func (c create) Create() interface{} {
	o := reflect.New(c.t).Interface()
	if oc, ok := o.(ObjectCreate); ok {
		oc.Ctor()
	}

	return o
}

type ObjectModule struct {
	service.Module
	defaultFactory *Factory // 默认对象工厂
	factorys       map[int]*Factory
	regs           map[string]create
	entitydelegate map[string]*EventDelegate
	sync           *SyncObject
	router         *ObjectRouter
}

func New() *ObjectModule {
	o := &ObjectModule{}
	o.defaultFactory = newFactory(o, share.OBJECT_TYPE_NONE)
	o.factorys = make(map[int]*Factory)
	o.regs = make(map[string]create)
	o.entitydelegate = make(map[string]*EventDelegate)
	o.sync = &SyncObject{o}
	o.router = NewObjectRouter(o)
	// 注册全局事件
	o.entitydelegate[GLOBAL_EVENT] = NewEventDelegate()
	o.AddFactory(share.OBJECT_TYPE_OBJECT)
	o.AddFactory(share.OBJECT_TYPE_GHOST)
	o.AddFactory(share.OBJECT_TYPE_SHARE)
	return o
}

// Name 模块名
func (o *ObjectModule) Name() string {
	return "Object"
}

// Init 模块初始化
func (o *ObjectModule) Init() bool {
	o.Core.RegisterRemote("object", o.sync)
	o.Core.RegisterRemote("ObjectRouter", o.router)

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
	return o.Core
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

// 指定在某个工厂创建
func (o *ObjectModule) FactoryCreate(identity int, typ string) (interface{}, error) {
	if f, has := o.factorys[identity]; has {
		return f.Create(typ)
	}

	return nil, fmt.Errorf("factory %d not found", identity)
}

// 指定在某个工厂创建容器
func (o *ObjectModule) FactoryCreateWithCap(identity int, typ string, cap int) (interface{}, error) {
	if f, has := o.factorys[identity]; has {
		return f.CreateWithCap(typ, cap)
	}

	return nil, fmt.Errorf("factory %d not found", identity)
}

// 创建
func (o *ObjectModule) Create(typ string) (interface{}, error) {
	return o.defaultFactory.Create(typ)
}

// 创建一个容器对象
func (o *ObjectModule) CreateWithCap(typ string, cap int) (interface{}, error) {
	return o.defaultFactory.CreateWithCap(typ, cap)
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

// FireObjectEvent 发送对象事件
func (o *ObjectModule) FireObjectEvent(event string, sender rpc.Mailbox, target rpc.Mailbox, args ...interface{}) (int, error) {
	tobj, _ := o.FindObject(target)
	sobj, _ := o.FindObject(sender)
	if sobj == nil || tobj == nil {
		return 0, fmt.Errorf("object is nil")
	}

	typ := tobj.(ObjectCreate).EntityType()
	if delegate, has := o.entitydelegate[typ]; has {
		ret := delegate.Invoke(event, sender, target, args...)
		return ret, nil
	}
	return 0, fmt.Errorf("object %s type not register ", typ)
}

// fireGlobalEvent 发送事件，但是不检查对象是否存在
func (o *ObjectModule) fireGlobalEvent(event string, sender rpc.Mailbox, target rpc.Mailbox, args ...interface{}) (int, error) {
	g := o.entitydelegate[GLOBAL_EVENT]
	ret := g.Invoke(event, sender, target, args...)
	return ret, nil
}

// 注册对象
func (o *ObjectModule) Register(oc ObjectCreate) {

	if oc == nil {
		panic("object: Register object is nil")
	}

	name := oc.EntityType()

	if name == "" {
		panic("entity name is empty")
	}

	if _, dup := o.regs[name]; dup {
		panic("object: Register called twice for object " + name)
	}

	v := reflect.ValueOf(oc)

	if v.Kind() != reflect.Ptr {
		panic("object: Register object must be pointer")
	}

	o.regs[name] = create{
		obj: oc,
		v:   v,
		t:   v.Type().Elem(),
	}
	o.entitydelegate[name] = NewEventDelegate()
	o.router.Register(name, oc)
}

// Replicate 将一个对象复制到另一个服务
func (o *ObjectModule) Replicate(objid rpc.Mailbox, dest rpc.Mailbox, tag int, cb ReplicateCB, cbparams interface{}) error {
	obj, err := o.FindObject(objid)
	if err != nil {
		return err
	}

	if oo, ok := obj.(Object); ok {
		f := oo.Factory()
		if f == nil {
			return fmt.Errorf("object factory is nil")
		}

		return f.Replicate(obj, dest, tag, cb, cbparams)
	}

	return fmt.Errorf("object is not implement Object")
}
