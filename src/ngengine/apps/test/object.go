package main

import (
	"ngengine/core/service"
	"ngengine/module/object"
	"ngengine/module/object/entity"
	"ngengine/module/object/game"
	"ngengine/module/replicate"
)

var objectargs = `{
	"ServId":1,
	"ServType": "object",
	"AdminAddr":"127.0.0.1",
	"AdminPort":12500,
	"ServName": "object",
	"ServAddr": "127.0.0.1",
	"ServPort": 0,
	"Expose": false,
	"HostAddr": "",
	"HostPort": 0,
	"LogFile":"object.log",
	"Args": {}
}`

type Object struct {
	service.BaseService
	object    *object.ObjectModule
	replicate *replicate.ReplicateModule
}

func (o *Object) Prepare(core service.CoreApi) error {
	o.CoreApi = core
	o.object = object.New()
	o.replicate = replicate.New()
	return nil
}

func (o *Object) Init(opt *service.CoreOption) error {
	o.CoreApi.AddModule(o.object)
	o.CoreApi.AddModule(o.replicate)

	o.object.Register("Player", &GamePlayerCreater{})
	o.replicate.RegisterReplicate("Player")
	return nil
}

func (o *Object) Start() error {
	p, _ := o.object.Create("Player")
	gp := p.(*GamePlayer)
	gp.Cache("hello", "hello")
	gp.SetSilence(true)
	gp.SetName("sll")
	gp.SetSilence(false)
	pos := gp.Pos()
	pos.Set(1, 1, 1)
	gp.SetPos(pos)
	o.object.Destroy(p)
	return nil
}

type GamePlayer struct {
	*game.Role
	*entity.Player
}

func NewGamePlayer() *GamePlayer {
	p := &GamePlayer{}
	p.Role = game.NewRole()
	p.Player = entity.NewPlayer()
	p.Role.SetSpirit(p.Player)
	return p
}

type GamePlayerCreater struct {
}

func (g *GamePlayerCreater) Create() interface{} {
	return NewGamePlayer()
}
