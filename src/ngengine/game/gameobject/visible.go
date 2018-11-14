package gameobject

import "time"

type VisibleObject struct {
	BaseObject
}

func (v *VisibleObject) OnCreate() {
	v.BaseObject.OnCreate()
}

func (v *VisibleObject) OnDestroy() {
	v.BaseObject.OnDestroy()
}

func (v *VisibleObject) Update(delta time.Duration) {
	v.BaseObject.Update(delta)
}
