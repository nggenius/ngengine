package gameobject

import (
	"time"
)

type RoleObject struct {
	TransformObject
}

func (r *RoleObject) OnCreate() {
	r.typ = OBJECT_PLAYER
	r.TransformObject.OnCreate()
}

func (r *RoleObject) OnDestroy() {
	r.TransformObject.OnDestroy()
}

func (r *RoleObject) Update(delta time.Duration) {
	r.TransformObject.Update(delta)
}
