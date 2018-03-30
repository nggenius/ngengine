package main

import (
	"ngengine/core/service"
	"ngengine/module/object"
	"ngengine/module/object/entity"
	"ngengine/module/object/game"
	"ngengine/module/replicate"
	"ngengine/module/store"
	"ngengine/module/timer"
	"ngengine/protocol"
)

var objectargs = `{
	"ServId":2,
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
	store     *store.StoreModule
	timer     *timer.TimerModule
}

func (o *Object) Prepare(core service.CoreApi) error {
	o.CoreApi = core
	o.object = object.New()
	o.replicate = replicate.New()
	o.store = store.New()
	o.timer = timer.New()
	return nil
}

func (o *Object) Init(opt *service.CoreOption) error {
	o.CoreApi.AddModule(o.object)
	o.CoreApi.AddModule(o.replicate)
	o.CoreApi.AddModule(o.store)
	o.CoreApi.AddModule(o.timer)
	// 设置store
	o.store.SetMode(store.STORE_CLIENT)
	o.store.Register().Register("Player", &GamePlayerData{})
	o.object.Register("Player", &GamePlayerCreater{})
	o.replicate.RegisterReplicate("Player")
	return nil
}

func (o *Object) Start() error {
	o.CoreApi.Watch("all")
	o.timer.AddCountTimer(1, 5000, nil, o.Timer)
	o.CoreApi.LogInfo("add timer")
	return nil
}

func (o *Object) Timer(id int64, count int, args interface{}) {
	o.CoreApi.LogInfo("timer")
	p, _ := o.object.Create("Player")
	gp := p.(*GamePlayer)
	gp.Cache("hello", "hello")
	gp.SetSilence(true)
	gp.SetName("sll")
	gp.SetSilence(false)
	pos := gp.Pos()
	pos.Set(1, 1, 1)
	gp.SetPos(pos)
	o.CoreApi.LogInfo(gp.Value("hello"))
	o.store.Client().Insert("object", "Player", gp.Archive(), o.InsertBack)
	o.object.Destroy(p)
	o.store.Client().Get("object", "Player", map[string]interface{}{"Id=?": 1}, o.LoadBack)
	o.store.Client().Find("object", "Player", map[string]interface{}{"Name=?": "sll"}, 4, 0, o.LoadAllBack)
}

func (o *Object) InsertBack(reply *protocol.Message) {
	err, ar := protocol.ParseReply(reply)
	if err != 0 {
		o.CoreApi.LogErr(err)
	}

	tag, _ := ar.ReadString()
	id, _ := ar.ReadInt64()
	o.CoreApi.LogInfo("insert ok", tag, " ", id)
}

func (o *Object) LoadBack(reply *protocol.Message) {
	err, ar := protocol.ParseReply(reply)
	tag, _ := ar.ReadString()
	if err != 0 {
		errstr, _ := ar.ReadString()
		o.CoreApi.LogErr(tag, errstr)
		return
	}

	load := &entity.PlayerArchive{}
	if err := ar.Read(load); err != nil {
		o.CoreApi.LogErr(err)
		return
	}

	o.CoreApi.LogInfo("load result: ", load)
}

func (o *Object) LoadAllBack(reply *protocol.Message) {
	err, ar := protocol.ParseReply(reply)
	tag, _ := ar.ReadString()
	if err != 0 {
		errstr, _ := ar.ReadString()
		o.CoreApi.LogErr(tag, errstr)
		return
	}

	var load []*entity.PlayerArchive
	if err := ar.Read(&load); err != nil {
		o.CoreApi.LogErr(err)
		return
	}

	for k, v := range load {
		o.CoreApi.LogInfo("load result: ", k, v)
	}

}

type GamePlayerData struct {
}

func (g *GamePlayerData) Create() interface{} {
	return &entity.PlayerArchive{}
}

func (g *GamePlayerData) CreateSlice() interface{} {
	return &[]*entity.PlayerArchive{}
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
