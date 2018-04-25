package entity

import (
	"ngengine/game/gameobject/entity/inner"
	"ngengine/module/store"
)

type Register interface {
	Register(name string, creater store.DataCreater) error
}

func RegisterToDB(r Register) {
	r.Register("inner.Account", &inner.AccountCreater{})
	r.Register("inner.Role", &inner.RoleCreater{})
	r.Register("entity.Player", &PlayerArchiveCreater{})
}
