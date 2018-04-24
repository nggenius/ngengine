package gameobject

import "time"

type VisibleObject struct {
	BaseObject
}

func (v *VisibleObject) Create() {
	v.BaseObject.Create()
}

func (v *VisibleObject) Destroy() {
	v.BaseObject.Destroy()
}

func (v *VisibleObject) Update(delta time.Duration) {
	v.BaseObject.Update(delta)
}
