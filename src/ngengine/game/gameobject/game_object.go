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
	// Init 初始化
	Init(object interface{})
	// Spirit 获取Object对象
	Spirit() object.Object
	// SetTransport 设置连接
	SetTransport(t *Transport)
	// Transport 获取连接
	Transport() *Transport
	// AddComponent 增加组件
	AddComponent(name string, com Component, update bool) error
	// RemoveComponent 移除组件
	RemoveComponent(name string)
	// 获取组件
	GetComponent(name string) Component
	// 获取gameobject
	GameObject() interface{}
	// Parent 父对象
	Parent() GameObject
	// SetParent 设置父对象
	SetParent(p GameObject)
	// ParentIndex 获取在父对象中的位置
	ParentIndex() int
	// SetParentIndex 设置在父对象中的位置
	SetParentIndex(pos int)
	// Cap 获取容量，返回-1表示不限容量
	Cap() int
}

// GameObjectEqual 判断两个对象是否相等
func GameObjectEqual(l GameObject, r GameObject) bool {
	if l.Spirit() == nil || r.Spirit() == nil {
		panic("object spirit is nil")
	}

	return l.Spirit().ObjId() == l.Spirit().ObjId()
}

type ComponentInfo struct {
	started   bool
	comp      Component
	useUpdate bool
}

type BaseObject struct {
	c *Container
	object.CacheData
	typ        int
	delete     bool
	index      int // 在factory中的索引
	client     rpc.Mailbox
	spirit     object.Object
	delegate   object.Delegate
	component  map[string]ComponentInfo
	transport  *Transport
	gameObject GameObject
	parent     GameObject
	pi         int  // 在父对象容器中的位置
	update     bool // 是否每一帧进行调用
}

// Init 初始化
func (b *BaseObject) Init(object interface{}) {
	if g, ok := object.(GameObject); ok {
		b.gameObject = g
		return
	}

	panic("object not implement GameObject")
}

// Parent 父对象
func (b *BaseObject) Parent() GameObject {
	return b.parent
}

// SetParent 设置父对象
func (b *BaseObject) SetParent(p GameObject) {
	b.parent = p
}

// ParentIndex 获取在父对象中的位置
func (b *BaseObject) ParentIndex() int {
	return b.pi
}

// SetParentIndex 设置在父对象中的位置
func (b *BaseObject) SetParentIndex(pi int) {
	b.pi = pi
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
