package models

import "ngengine/module/object"

type reg interface {
	Register(oc object.ObjectCreate)
}

func Register(r reg) {
	r.Register(new(GamePlayer))
}
