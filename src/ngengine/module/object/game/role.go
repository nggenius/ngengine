package game

import (
	"time"
)

type Role struct {
	TransformObject
}

func NewRole() *Role {
	r := &Role{}
	r.typ = OBJECT_PLAYER
	return r
}

func (r *Role) Create() {
	r.TransformObject.Create()
}

func (r *Role) Destroy() {
	r.TransformObject.Destroy()
}

func (r *Role) Update(delta time.Duration) {
	r.TransformObject.Update(delta)
}
