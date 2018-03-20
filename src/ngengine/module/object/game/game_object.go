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
	index     int
	mailbox   rpc.Mailbox
	spirit    object.Object
	delegate  delegate
	component map[string]ComponentInfo
}

func (b *BaseObject) Create() {
	b.component = make(map[string]ComponentInfo)
}

func (b *BaseObject) Destroy() {

}

func (b *BaseObject) SetIndex(index int) {
	b.index = index
}

func (b *BaseObject) Index() int {
	return b.index
}

func (b *BaseObject) SetDelegate(d delegate) {
	b.delegate = d
}

func (b *BaseObject) Spirit() object.Object {
	return b.spirit
}

func (b *BaseObject) SetSpirit(s object.Object) {
	b.spirit = s
}

func (b *BaseObject) Mailbox() rpc.Mailbox {
	return b.mailbox
}

func (b *BaseObject) SetMailbox(mb rpc.Mailbox) {
	b.mailbox = mb
}

func (b *BaseObject) AddComponent(name string, com Component, update bool) error {
	if _, has := b.component[name]; has {
		return fmt.Errorf("component has register twice, %s ", name)
	}

	b.component[name] = ComponentInfo{com, update}
	return nil
}

func (b *BaseObject) RemoveComponent(name string) {
	if _, has := b.component[name]; has {
		delete(b.component, name)
	}
}
