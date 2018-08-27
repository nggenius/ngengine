package entity

import (
	"ngengine/game/gameobject/entity/inner"
	"ngengine/module/object"
)

const (
	ACCOUNT       = "inner.Account"
	ROLE          = "inner.Role"
	ROLE_SAVE     = "entyty.Player"
	ROLE_SAVE_BAK = "entity.PlayerBak"
)

var objreg = make(map[string]func() object.Object)

type Register interface {
	Register(name string, obj interface{}, objslice interface{}) error
}

func RegisterToDB(r Register) {
	r.Register(ACCOUNT, &inner.Account{}, []*inner.Account{})
	r.Register(ROLE, &inner.Role{}, []*inner.Role{})
	r.Register(ROLE_SAVE, &PlayerArchive{}, []*PlayerArchive{})
	r.Register(ROLE_SAVE_BAK, &PlayerArchiveBak{}, []*PlayerArchiveBak{})
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
