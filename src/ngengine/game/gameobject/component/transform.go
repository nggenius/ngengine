package component

import (
	"ngengine/game/gameobject"
)

type Transform struct {
	gameobject.Component
}

func (t *Transform) Create() {

}

func (t *Transform) LookAtTarget(target Transform) {

}

func (t *Transform) LookAtPoint(x, y, z float32) {

}

func (t *Transform) RotateEulerAngles(x, y, z float32) {

}

func (t *Transform) RotateDirAngle(x, y, z float32, angle float32) {

}

func (t *Transform) RotateAngle(xAngle, yAngle, zAngle float32) {

}

func (t *Transform) Translate(x, y, z float32) {

}
