package models

import "ngengine/module/object"

type reg interface {
	Register(name string, oc object.ObjectCreate)
}

func Register(r reg) {
	r.Register("entity.Player", new(GamePlayerCreater))
}
