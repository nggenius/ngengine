package object

import (
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
}

type Factory struct {
	objType int
	serial  uint32
	owner   *ObjectModule
	pool    *ObjectList
}

func newFactory(owner *ObjectModule, typ int) *Factory {
	f := &Factory{}
	f.objType = typ
	f.owner = owner
	f.pool = NewObjectList(128, 0x1000000)
	return f
}

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
			o.SetFactory(f)
			return inst, nil
		}

		return nil, fmt.Errorf("new obj is not Object")
	}
	return nil, fmt.Errorf("object not found")
}

func (f *Factory) Destroy(object interface{}) error {
	if fo, ok := object.(FactoryObject); ok {
		return f.pool.Remove(fo.Index(), object)
	}

	return errors.New("object is not FactoryObject")
}
