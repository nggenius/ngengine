package object

import (
	"container/list"
	"errors"
	"fmt"
	"ngengine/core/rpc"
	"ngengine/core/service"
	"ngengine/share"
	"ngengine/utils"
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
	OnCreate()
	OnDestroy()
	OnDelete()
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

func (f *Factory) CreateWithCap(typ string, cap int) (interface{}, error) {
	o, err := f.Create(typ)
	if err != nil {
		return nil, err
	}

	if c, ok := o.(Container); ok {
		c.SetCap(cap)
		return o, nil
	}

	f.Destroy(o)
	return nil, fmt.Errorf("%s not have SetCap", typ)
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
			o.OnCreate()

			f.owner.Core.LogDebug("create object ", o.ObjId())
			return inst, nil
		}

		return nil, fmt.Errorf("new object type %s not implement FactoryObject", typ)
	}
	return nil, fmt.Errorf("object %s not found", typ)
}

// 销毁一个对象
func (f *Factory) Destroy(object interface{}) error {
	if fo, ok := object.(FactoryObject); ok {
		if fo.Alive() {
			fo.OnDestroy()
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
		fo.OnDelete()
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

// Replicate 将一个对象复制到另一个服务
func (f *Factory) Replicate(object interface{}, dest rpc.Mailbox, tag int, cb ReplicateCB, cbparams interface{}) error {
	if dest.IsNil() {
		return fmt.Errorf("dest is nil")
	}

	if object == nil {
		return fmt.Errorf("object is nil")
	}

	data, err := f.Encode(object)
	if err != nil {
		return err
	}

	return f.owner.Core.MailtoAndCallback(nil, &dest, "object.Replicate", f.onReplicate, []interface{}{cb, cbparams, object.(Object).ObjId()}, tag, data)
}

func (f *Factory) onReplicate(param interface{}, replyerr *rpc.Error, ar *utils.LoadArchive) {
	if param == nil {
		return
	}

	args := param.([]interface{})
	cb := args[0].(ReplicateCB)
	id := args[2].(rpc.Mailbox)
	o, _ := f.FindObject(id)
	if o == nil {
		f.owner.Core.LogErr("object not found")
		return
	}

	if replyerr != nil {
		if cb != nil {
			cb(args[1], replyerr)
		}
		return
	}
	var dummy rpc.Mailbox
	if err := ar.Read(&dummy); err != nil {
		if cb != nil {
			cb(args[1], err)
		}
		return
	}
	obj := o.(Object)
	obj.AddDummy(dummy, DUMMY_STATE_READY)

	if cb != nil {
		cb(args[1], nil)
	}
}
