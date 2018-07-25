package object

import (
	"container/list"
	"errors"
	"fmt"
	"ngengine/core/rpc"
	"ngengine/core/service"
	"ngengine/share"
)

type FactoryObject interface {
	Init(interface{})
	Index() int
	SetIndex(int)
	SetObjId(mb rpc.Mailbox)
	ObjId() rpc.Mailbox
	Factory() *Factory
	SetFactory(f *Factory)
	SetCore(c service.CoreAPI)
	Prepare()
	Create()
	Destroy()
	Delete()
	Alive() bool
	SetDelegate(d Delegate)
}

type Factory struct {
	identity int
	serial   int
	owner    *ObjectModule
	pool     *ObjectList
	delete   *list.List
}

func newFactory(owner *ObjectModule, identity int) *Factory {
	f := &Factory{}
	f.identity = identity
	f.owner = owner
	f.pool = NewObjectList(128, share.OBJECT_MAX)
	f.delete = list.New()
	return f
}

// 通过类型创建一个对象
func (f *Factory) Create(typ string) (interface{}, error) {
	if c, ok := f.owner.regs[typ]; ok {
		inst := c.Create()
		if inst == nil {
			return nil, fmt.Errorf("object %s create failed", typ)
		}

		if o, ok := inst.(FactoryObject); ok {
			index, err := f.pool.Add(inst)
			if err != nil {
				return nil, err
			}
			o.Init(inst)
			o.SetIndex(index)
			f.serial = (f.serial + 1) % 0xFF
			o.SetObjId(f.owner.Core.Mailbox().NewObjectId(f.identity, f.serial, index))
			o.SetFactory(f)
			o.SetCore(f.owner.Core)
			o.SetDelegate(f.owner.entitydelegate[typ])
			o.Prepare()
			o.Create()

			f.owner.Core.LogDebug("create object ", o.ObjId())
			return inst, nil
		}

		return nil, fmt.Errorf("new object type %s not implement FactoryObject", typ)
	}
	return nil, fmt.Errorf("object %s not found", typ)
}

// 销毁一个对象``
func (f *Factory) Destroy(object interface{}) error {
	if fo, ok := object.(FactoryObject); ok {
		if fo.Alive() {
			fo.Destroy()
			f.delete.PushBack(object)
			return f.pool.Remove(fo.Index(), object)
		}
		return nil
	}
	return errors.New("object is not implement FactoryObject")
}

// 清理需要删除的对象
func (f *Factory) ClearDelete() {
	for ele := f.delete.Front(); ele != nil; {
		fo := ele.Value.(FactoryObject)
		fo.Delete()
		e := ele
		ele = ele.Next()
		f.delete.Remove(e)
		f.owner.Core.LogDebug("delete object ", fo.ObjId())
	}
}

// 查找对象
func (f *Factory) FindObject(mb rpc.Mailbox) (interface{}, error) {
	if f.owner.Core.Mailbox().ServiceId() != mb.ServiceId() ||
		mb.Flag() != 0 ||
		mb.Identity() != f.identity {
		return nil, fmt.Errorf("mailbox %s error", mb)
	}

	obj, err := f.pool.Get(mb.ObjectIndex())
	if err != nil {
		return nil, err
	}

	if obj != nil && obj.(FactoryObject).Alive() {
		return obj, nil
	}

	return nil, fmt.Errorf("object has destroyed, %s", mb)
}
