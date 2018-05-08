package gameobject

import (
	"fmt"
	"ngengine/core/rpc"
	"ngengine/module/object"
	"time"
)

const (
	OBJECT_NONE = iota
	OBJECT_SCENE
	OBJECT_PLAYER
	OBJECT_ITEM
	OBJECT_NPC
	OBJECT_MAX
)

type GameObject interface {
	Spirit() object.Object
	AddComponent(name string, com Component) error
	RemoveComponent(name string)
}

type ComponentInfo struct {
	started   bool
	comp      Component
	useUpdate bool
}

type BaseObject struct {
	object.CacheData
	typ       int
	delete    bool
	index     int
	objid     rpc.Mailbox
	client    rpc.Mailbox
	spirit    object.Object
	delegate  object.Delegate
	component map[string]ComponentInfo
}

// 预处理
func (b *BaseObject) Prepare() {
	b.component = make(map[string]ComponentInfo)
	b.CacheData.Init()
}

// 构造函数
func (b *BaseObject) Create() {
	if b.delegate != nil {
		b.delegate.Invoke(E_ON_CREATE, b.objid, rpc.NullMailbox)
	}
}

func (b *BaseObject) ObjectType() int {
	return b.typ
}

// 准备销毁
func (b *BaseObject) Destroy() {
	if b.delegate != nil {
		b.delegate.Invoke(E_ON_DESTROY, b.objid, rpc.NullMailbox)
	}
	b.delete = true
}

// 是否还活着
func (b *BaseObject) Alive() bool {
	return !b.delete
}

// 正式开始删除
func (b *BaseObject) Delete() {

}

// 设置索引，由factory调用，不要手工调用
func (b *BaseObject) SetIndex(index int) {
	b.index = index
}

// 获取索引
func (b *BaseObject) Index() int {
	return b.index
}

// 设置事件代理
func (b *BaseObject) SetDelegate(d object.Delegate) {
	b.delegate = d
}

// 精神实体，数据部分
func (b *BaseObject) Spirit() object.Object {
	return b.spirit
}

// 设置精神实体
func (b *BaseObject) SetSpirit(s object.Object) {
	b.spirit = s
}

// 唯一ID
func (b *BaseObject) ObjId() rpc.Mailbox {
	return b.objid
}

// 设置唯一ID
func (b *BaseObject) SetObjId(id rpc.Mailbox) {
	b.objid = id
}

// 客户端地址
func (b *BaseObject) Client() rpc.Mailbox {
	return b.client
}

// 设置客户端地址
func (b *BaseObject) SetClient(mb rpc.Mailbox) {
	b.client = mb
}

// update
func (b *BaseObject) Update(delta time.Duration) {
	for _, comp := range b.component {
		if !comp.comp.Enable() {
			continue
		}
		if !comp.started {
			comp.comp.Start()
			comp.started = true
		}
		if comp.useUpdate {
			comp.comp.Update(delta)
		}
	}
}

// 获取组件
func (b *BaseObject) GetComponent(name string) Component {
	if comp, has := b.component[name]; has {
		return comp.comp
	}
	return nil
}

// 增加组件
func (b *BaseObject) AddComponent(name string, com Component, update bool) error {
	if _, has := b.component[name]; has {
		return fmt.Errorf("component has register twice, %s ", name)
	}

	b.component[name] = ComponentInfo{
		started:   false,
		comp:      com,
		useUpdate: update,
	}

	com.SetEnable(true)
	// 调用初始化函数
	com.Create()
	return nil
}

// 移除组件
func (b *BaseObject) RemoveComponent(name string) {
	if comp, has := b.component[name]; has {
		comp.comp.Destroy() // 销毁组件
		delete(b.component, name)
	}
}
