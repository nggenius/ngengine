package scene

type GameSceneCreater struct {
}

func (g *GameSceneCreater) Create() interface{} {
	s := NewGameScene()
	return s
}
