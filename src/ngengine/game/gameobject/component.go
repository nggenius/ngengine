package gameobject

import (
	"time"
)

type Component interface {
	SetGameObject(interface{})
	GameObject() interface{}
	Create()
	Start()
	Update(delta time.Duration)
	Destroy()
	Enable() bool
	SetEnable(e bool)
}

type GameComponent struct {
	gameObject interface{}
	enable     bool
}

func (g *GameComponent) SetGameObject(obj interface{}) {
	g.gameObject = obj
}

func (g *GameComponent) GameObject() interface{} {
	return g.gameObject
}

func (g *GameComponent) Enable() bool {
	return g.enable
}

func (g *GameComponent) SetEnable(e bool) {
	g.enable = e
}

func (g *GameComponent) Create() {

}

func (g *GameComponent) Start() {

}

func (g *GameComponent) Update(delta time.Duration) {

}

func (g *GameComponent) Destroy() {

}
