package store

import (
	"fmt"
	"reflect"
)

type objectcreater struct {
	obj   interface{}
	value reflect.Value
	typ   reflect.Type
	stype reflect.Type
}

func (oc objectcreater) Create() interface{} {
	return reflect.New(oc.typ.Elem()).Interface()
}

func (oc objectcreater) CreateSlice() interface{} {
	return reflect.New(oc.stype).Interface()
}

type Register struct {
	types map[string]objectcreater
}

func newRegister() *Register {
	r := &Register{}
	r.types = make(map[string]objectcreater)
	return r
}

func (r *Register) Register(name string, obj interface{}, objslice interface{}) error {
	if _, dup := r.types[name]; dup {
		return fmt.Errorf("register data twice, %s", name)
	}

	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		panic("register object must be a pointer")
	}

	oc := objectcreater{
		obj:   obj,
		value: v,
		typ:   v.Type(),
		stype: reflect.TypeOf(objslice),
	}
	r.types[name] = oc
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
		ctx.Core.LogInfo("sync ", k, " ok")
	}
	return nil
}
