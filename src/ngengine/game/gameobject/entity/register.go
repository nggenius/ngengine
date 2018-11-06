package entity

import (
	"ngengine/game/gameobject/entity/inner"
	"ngengine/module/object"
)

const (
	ACCOUNT       = "inner.Account"
	ROLE          = "inner.Role"
	DB_PLAYER     = "entity.Player"
	DB_PLAYER_BAK = "entity.PlayerBak"
)

var objreg = make(map[string]func() object.Object)

type Register interface {
	Register(name string, obj interface{}, objslice interface{}) error
}

func RegisterToDB(r Register) {
	r.Register(ACCOUNT, &inner.Account{}, []*inner.Account{})
	r.Register(ROLE, &inner.Role{}, []*inner.Role{})
	r.Register(DB_PLAYER, &PlayerArchive{}, []*PlayerArchive{})
	r.Register(DB_PLAYER_BAK, &PlayerArchiveBak{}, []*PlayerArchiveBak{})
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
