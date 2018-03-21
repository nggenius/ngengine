package game

import (
	"fmt"
	"ngengine/core/rpc"
	"ngengine/module/object"
)

type delegate interface {
	Invoke(event string, self, sender object.Object, args ...interface{}) int
}

type GameObject interface {
	Spirit() object.Object
	AddComponent(name string, com Component) error
	RemoveComponent(name string)
}

type ComponentInfo struct {
	comp      Component
	useUpdate bool
}

type BaseObject struct {
	delete    bool
	index     int
	mailbox   rpc.Mailbox
	spirit    object.Object
	delegate  delegate
	component map[string]ComponentInfo
}

func (b *BaseObject) Create() {
	b.component = make(map[string]ComponentInfo)
}

// 准备销毁
func (b *BaseObject) Destroy() {
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
func (b *BaseObject) SetDelegate(d delegate) {
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

// 邮箱地址
func (b *BaseObject) Mailbox() rpc.Mailbox {
	return b.mailbox
}

// 设置邮箱地址
func (b *BaseObject) SetMailbox(mb rpc.Mailbox) {
	b.mailbox = mb
}

// 增加一个组件
func (b *BaseObject) AddComponent(name string, com Component, update bool) error {
	if _, has := b.component[name]; has {
		return fmt.Errorf("component has register twice, %s ", name)
	}

	b.component[name] = ComponentInfo{com, update}
	return nil
}

// 移除一个组件
func (b *BaseObject) RemoveComponent(name string) {
	if _, has := b.component[name]; has {
		delete(b.component, name)
	}
}
