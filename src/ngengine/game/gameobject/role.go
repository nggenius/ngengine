package gameobject

import (
	"time"
)

type RoleObject struct {
	TransformObject
}

func (r *RoleObject) Create() {
	r.typ = OBJECT_PLAYER
	r.TransformObject.Create()
}

func (r *RoleObject) Destroy() {
	r.TransformObject.Destroy()
}

func (r *RoleObject) Update(delta time.Duration) {
	r.TransformObject.Update(delta)
}
