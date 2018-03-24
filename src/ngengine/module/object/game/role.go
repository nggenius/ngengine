package game

import (
	"ngengine/module/object"
	"ngengine/module/object/entity"
	"time"
)

type Role struct {
	TransformObject
	*entity.Player
}

func NewRole() *Role {
	r := &Role{}
	r.Player = entity.NewPlayer()
	r.SetSpirit(r.Player)
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

type RoleCreater struct {
}

// create player
func (o *RoleCreater) Create() interface{} {
	return NewRole()
}

func init() {
	object.Register("Role", &RoleCreater{})
}
