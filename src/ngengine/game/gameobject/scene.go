package gameobject

import "time"

type SceneObject struct {
	BaseObject
}

func (s *SceneObject) Create() {
	s.typ = OBJECT_SCENE
	s.BaseObject.Create()
}

func (s *SceneObject) Destroy() {
	s.BaseObject.Destroy()
}

func (s *SceneObject) Update(delta time.Duration) {
	s.BaseObject.Update(delta)
}
