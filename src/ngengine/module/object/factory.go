package object

import (
	"container/list"
	"errors"
	"fmt"
	"ngengine/core/rpc"
)

var regs = map[string]ObjectCreate{}

func Register(name string, o ObjectCreate) {

	if o == nil {
		panic("object: Register object is nil")
	}
	if _, dup := regs[name]; dup {
		panic("object: Register called twice for object " + name)
	}

	regs[name] = o
}

type FactoryObject interface {
	Index() int
	SetIndex(int)
	SetMailbox(mb rpc.Mailbox)
	SetFactory(f *Factory)
	Create()
	Destroy()
	Delete()
	Alive() bool
}

type Factory struct {
	objType int
	serial  int
	owner   *ObjectModule
	pool    *ObjectList
	delete  *list.List
}

func newFactory(owner *ObjectModule, typ int) *Factory {
	f := &Factory{}
	f.objType = typ
	f.owner = owner
	f.pool = NewObjectList(128, 0x1000000)
	f.delete = list.New()
	return f
}

// 通过类型创建一个对象
func (f *Factory) Create(typ string) (interface{}, error) {
	if c, ok := regs[typ]; ok {
		inst := c.Create()
		if inst == nil {
			return nil, fmt.Errorf("object create failed")
		}

		if o, ok := inst.(FactoryObject); ok {
			index, err := f.pool.Add(inst)
			if err != nil {
				return nil, err
			}
			o.SetIndex(index)
			f.serial = (f.serial + 1) % 0xFF
			o.SetMailbox(f.owner.core.Mailbox().ObjectId(f.objType, f.serial, index))
			o.SetFactory(f)
			o.Create()
			return inst, nil
		}

		return nil, fmt.Errorf("new obj is not object")
	}
	return nil, fmt.Errorf("object not found")
}

// 销毁一个对象
func (f *Factory) Destroy(object interface{}) error {
	if fo, ok := object.(FactoryObject); ok {
		if fo.Alive() {
			fo.Destroy()
			f.delete.PushBack(object)
			return f.pool.Remove(fo.Index(), object)
		}
		return nil
	}
	return errors.New("object is not FactoryObject")
}

// 清理需要删除的对象
func (f *Factory) ClearDelete() {
	for ele := f.delete.Front(); ele != nil; ele = ele.Next() {
		fo := ele.Value.(FactoryObject)
		fo.Delete()
		e := ele
		ele = ele.Next()
		f.delete.Remove(e)
	}
}

// 获取对象
func (f *Factory) GetObject(mb rpc.Mailbox) (interface{}, error) {
	if f.owner.core.Mailbox().Sid != mb.Sid ||
		mb.Flag != 0 ||
		(mb.Id>>40)&0x7F != uint64(f.objType) {
		return nil, errors.New("mailbox error")
	}

	obj, err := f.pool.Get(int(mb.Id & 0xFFFFFFFF))
	if err != nil {
		return nil, err
	}

	if obj != nil && obj.(FactoryObject).Alive() {
		return obj, nil
	}

	return nil, fmt.Errorf("object has destroyed, %s", mb)
}
