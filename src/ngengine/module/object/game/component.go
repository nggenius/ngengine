package game

import (
	"time"
)

type Component interface {
	Create()
	Start()
	Stop()
	Update(delta time.Duration)
	Destroy()
}

type GameComponent struct {
	gameObject GameObject
}

func (g *GameComponent) Create() {

}

func (g *GameComponent) Start() {

}

func (g *GameComponent) Stop() {

}

func (g *GameComponent) Update(delta time.Duration) {

}

func (g *GameComponent) Destroy() {

}
