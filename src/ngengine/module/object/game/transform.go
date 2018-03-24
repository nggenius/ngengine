package game

import "time"

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
