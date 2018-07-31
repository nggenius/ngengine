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
	// 初始化
	Init(object interface{})
	// 获取Object对象
	Spirit() object.Object
	// 设置连接
	SetTransport(t *Transport)
	// 获取连接
	Transport() *Transport
	// 增加组件
	AddComponent(name string, com Component, update bool) error
	// 移除组件
	RemoveComponent(name string)
	// 获取组件
	GetComponent(name string) Component
	// 获取gameobject
	GameObject() interface{}
	// Parent 父对象
	Parent() GameObject
	// SetParent 设置父对象
	SetParent(p GameObject)
	// Pos 获取在父对象中的位置
	ContainerPos() int
	// SetPos 设置在父对象中的位置
	SetContainerPos(pos int)
}

type ComponentInfo struct {
	started   bool
	comp      Component
	useUpdate bool
}

type BaseObject struct {
	Container
	object.CacheData
	typ        int
	delete     bool
	index      int
	client     rpc.Mailbox
	spirit     object.Object
	delegate   object.Delegate
	component  map[string]ComponentInfo
	transport  *Transport
	gameObject interface{}
	parent     GameObject
	pos        int //在父对象中的位置
}

// Init 初始化
func (b *BaseObject) Init(object interface{}) {
	b.gameObject = object
}

// Parent 父对象
func (b *BaseObject) Parent() GameObject {
	return b.parent
}

// SetParent 设置父对象
func (b *BaseObject) SetParent(p GameObject) {
	b.parent = p
}

// Pos 获取在父对象中的位置
func (b *BaseObject) ContainerPos() int {
	return b.pos
}

// SetPos 设置在父对象中的位置
func (b *BaseObject) SetContainerPos(pos int) {
	b.pos = pos
}

// GameObject 获取gameobject
func (b *BaseObject) GameObject() interface{} {
	return b.gameObject
}

// SetTransport 设置连接
func (b *BaseObject) SetTransport(t *Transport) {
	b.transport = t
}

// Transport 获取连接
func (b *BaseObject) Transport() *Transport {
	return b.transport
}

// Prepare 预处理
func (b *BaseObject) Prepare() {
	b.component = make(map[string]ComponentInfo)
	b.CacheData.Init()
}

// Create 构造函数
func (b *BaseObject) Create() {
	if b.delegate != nil && b.spirit != nil {
		b.delegate.Invoke(E_ON_CREATE, b.spirit.ObjId(), rpc.NullMailbox)
	}
}

// ObjectType 获取对象类型
func (b *BaseObject) ObjectType() int {
	return b.typ
}

// Destroy 准备销毁
func (b *BaseObject) Destroy() {
	if b.delegate != nil && b.spirit != nil {
		b.delegate.Invoke(E_ON_DESTROY, b.spirit.ObjId(), rpc.NullMailbox)
	}
	b.delete = true
}

// Alive 是否还活着
func (b *BaseObject) Alive() bool {
	return !b.delete
}

// Delete 正式开始删除
func (b *BaseObject) Delete() {

}

// SetIndex 设置索引，由factory调用，不要手工调用
func (b *BaseObject) SetIndex(index int) {
	b.index = index
}

// Index 获取索引
func (b *BaseObject) Index() int {
	return b.index
}

// SetDelegate 设置事件代理
func (b *BaseObject) SetDelegate(d object.Delegate) {
	b.delegate = d
}

// Spirit 精神实体，数据部分
func (b *BaseObject) Spirit() object.Object {
	return b.spirit
}

// SetSpirit 设置精神实体
func (b *BaseObject) SetSpirit(s object.Object) {
	b.spirit = s
}

// Client 客户端地址
func (b *BaseObject) Client() rpc.Mailbox {
	return b.client
}

// SetClient 设置客户端地址
func (b *BaseObject) SetClient(mb rpc.Mailbox) {
	b.client = mb
}

// Update
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

// GetComponent 获取组件
func (b *BaseObject) GetComponent(name string) Component {
	if comp, has := b.component[name]; has {
		return comp.comp
	}
	return nil
}

// AddComponent 增加组件
func (b *BaseObject) AddComponent(name string, com Component, update bool) error {
	if _, has := b.component[name]; has {
		return fmt.Errorf("component has register twice, %s ", name)
	}

	b.component[name] = ComponentInfo{
		started:   false,
		comp:      com,
		useUpdate: update,
	}

	com.SetGameObject(b.gameObject)
	com.SetEnable(true)
	// 调用初始化函数
	com.Create()
	return nil
}

// RemoveComponent 移除组件
func (b *BaseObject) RemoveComponent(name string) {
	if comp, has := b.component[name]; has {
		comp.comp.Destroy() // 销毁组件
		delete(b.component, name)
	}
}
