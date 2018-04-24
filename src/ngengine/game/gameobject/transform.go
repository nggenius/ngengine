package gameobject

import "time"

type Transform interface {
	SetPosXYZ(x float32, y float32, z float32)
	GetPosXYZ() (x float32, y float32, z float32)
	SetOrient(orient float32)
	Orient() float32
}

type TransformObject struct {
	VisibleObject
}

func (t *TransformObject) Create() {
	t.VisibleObject.Create()
}

func (t *TransformObject) Destroy() {
	t.VisibleObject.Destroy()
}

func (t *TransformObject) Update(delta time.Duration) {
	t.VisibleObject.Update(delta)
}
