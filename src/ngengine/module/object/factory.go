package object

import (
	"fmt"
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

type Factory struct {
	serial uint32
	owner  *ObjectModule
}

func newFactory(owner *ObjectModule) *Factory {
	f := &Factory{}
	f.owner = owner
	return f
}

func (f *Factory) Create(typ string) (Object, error) {
	return f.createObj(typ)
}

func (f *Factory) createObj(typ string) (Object, error) {
	if c, ok := regs[typ]; ok {
		inst := c.Create()
		if inst == nil {
			return nil, fmt.Errorf("object create failed")
		}

		if o, ok := inst.(Object); ok {
			o.SetFactory(f)
			return o, nil
		}

		return nil, fmt.Errorf("new obj is not Object")
	}
	return nil, fmt.Errorf("object not found")
}
