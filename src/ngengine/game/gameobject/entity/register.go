package entity

import (
	"ngengine/game/gameobject/entity/inner"
	"ngengine/module/object"
	"ngengine/module/store"
)

var objreg = make(map[string]func() object.Object)

type Register interface {
	Register(name string, creater store.DataCreater) error
}

func RegisterToDB(r Register) {
	r.Register("inner.Account", &inner.AccountCreater{})
	r.Register("inner.Role", &inner.RoleCreater{})
	r.Register("entity.Player", &PlayerArchiveCreater{})
}

func registObject(typ string, f func() object.Object) {
	if _, has := objreg[typ]; has {
		panic("register object twice")
	}

	objreg[typ] = f
}

func Create(typ string) object.Object {
	if c, ok := objreg[typ]; ok {
		return c()
	}
	return nil
}
