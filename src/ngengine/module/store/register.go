package store

import (
	"fmt"
)

type Register struct {
	types map[string]DataCreater
}

type DataCreater interface {
	Create() interface{}
	CreateSlice() interface{}
}

func newRegister() *Register {
	r := &Register{}
	r.types = make(map[string]DataCreater)
	return r
}

func (r *Register) Register(name string, creater DataCreater) error {
	if _, dup := r.types[name]; dup {
		return fmt.Errorf("register data twice, %s", name)
	}

	r.types[name] = creater
	return nil
}

func (r *Register) Create(name string) interface{} {
	if c, has := r.types[name]; has {
		return c.Create()
	}
	return nil
}

func (r *Register) CreateSlice(name string) interface{} {
	if c, has := r.types[name]; has {
		return c.CreateSlice()
	}
	return nil
}

func (r *Register) Sync(ctx *StoreModule) error {
	for k, v := range r.types {
		err := ctx.sql.Sync(v.Create())
		if err != nil {
			return err
		}
		ctx.core.LogInfo("sync ", k, " ok")
	}
	return nil
}
