package gameobject

import "time"

type SceneObject struct {
	BaseObject
}

func (s *SceneObject) OnCreate() {
	s.typ = OBJECT_SCENE
	s.BaseObject.OnCreate()
}

func (s *SceneObject) OnDestroy() {
	s.BaseObject.OnDestroy()
}

func (s *SceneObject) Update(delta time.Duration) {
	s.BaseObject.Update(delta)
}
