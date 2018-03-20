package game

type VisibleObject struct {
	BaseObject
}

func (v *VisibleObject) Create() {
	v.BaseObject.Create()
}

func (v *VisibleObject) Destroy() {
	v.BaseObject.Destroy()
}
