package game

type TransformObject struct {
	VisibleObject
}

func (t *TransformObject) Create() {
	t.VisibleObject.Create()
}

func (t *TransformObject) Destroy() {
	t.VisibleObject.Destroy()
}
