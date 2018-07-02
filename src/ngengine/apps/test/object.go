package main

// 测试对象模块
import (
	"ngengine/core/service"
	"ngengine/game/gameobject"
	"ngengine/game/gameobject/entity"
	"ngengine/module/object"
	"ngengine/module/replicate"
	"ngengine/module/store"
	"ngengine/module/timer"
	"ngengine/protocol"

	"github.com/mysll/toolkit"
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

func (o *Object) Prepare(core service.CoreAPI) error {
	o.CoreAPI = core
	o.object = object.New()
	o.replicate = replicate.New()
	o.store = store.New()
	o.timer = timer.New()
	return nil
}

func (o *Object) Init(opt *service.CoreOption) error {
	o.CoreAPI.AddModule(o.object)
	o.CoreAPI.AddModule(o.replicate)
	o.CoreAPI.AddModule(o.store)
	o.CoreAPI.AddModule(o.timer)
	// 设置store
	o.store.SetMode(store.STORE_CLIENT)
	o.store.Register("Player", &entity.PlayerArchiveCreater{})
	o.object.Register("Player", &GamePlayerCreater{})
	o.replicate.RegisterReplicate("Player")
	return nil
}

func (o *Object) Start() error {
	o.CoreAPI.Watch("all")
	o.timer.AddCountTimer(1, 5000, nil, o.Timer)
	o.CoreAPI.LogInfo("add timer")
	return nil
}

func (o *Object) Timer(id int64, count int, args interface{}) {
	//o.store.Client().Exec("object", "DELETE from `player`", []interface{}{}, o.ExecBack)
	o.store.Client().Query("object", "select * from player", []interface{}{}, o.QueryBack)
	o.CoreAPI.LogInfo("timer")
	p, _ := o.object.Create("Player")
	gp := p.(*GamePlayer)
	gp.Cache("hello", "hello")
	gp.SetSilence(true)
	gp.SetName("sll")
	gp.SetSilence(false)
	pos := gp.Pos()
	pos.Set(1, 1, 1)
	gp.SetPos(pos)
	o.CoreAPI.LogInfo(gp.Value("hello"))
	o.store.Client().Insert("object", "Player", gp.Archive(), o.InsertBack)
	o.store.Client().Insert("object", "Player", gp.Archive(), o.InsertBack)
	o.store.Client().Insert("object", "Player", gp.Archive(), o.InsertBack)
	o.object.Destroy(p)
	o.store.Client().Get("object", "Player", map[string]interface{}{"Id=?": -1}, o.LoadBack)
	o.store.Client().Get("object", "Player", map[string]interface{}{"Id=?": 32}, o.LoadBack)
	o.store.Client().Find("object", "Player", map[string]interface{}{"Name=?": "sll"}, 4, 0, o.LoadAllBack)
}

func (o *Object) QueryBack(reply *protocol.Message) {
	errcode, err, tag, result := store.ParseQueryReply(reply)
	if err != nil {
		o.CoreAPI.LogErr(errcode, err, tag)
		return
	}
	o.CoreAPI.LogInfo("query result:", result)
}

func (o *Object) ExecBack(reply *protocol.Message) {
	errcode, err, tag, affected := store.ParseExecReply(reply)
	if err != nil {
		o.CoreAPI.LogErr(errcode, err, tag)
		return
	}
	o.CoreAPI.LogInfo("exec result:", affected)
}

func (o *Object) InsertBack(reply *protocol.Message) {
	errcode, err, tag, affected, id := store.ParseInsertReply(reply)
	if err != nil {
		o.CoreAPI.LogErr(errcode, err, tag)
		return
	}
	o.CoreAPI.LogInfo("insert ", tag, ", ", affected, ",id: ", id)
}

func (o *Object) LoadBack(reply *protocol.Message) {

	load := &entity.PlayerArchive{}
	errcode, err, tag := store.ParseGetReply(reply, load)
	if err != nil {
		o.CoreAPI.LogErr(errcode, err, tag)
		return
	}
	o.CoreAPI.LogInfo("load result: ", load)

	load.Orient = toolkit.RandRangef(0, 3.1415926*2)
	o.store.Client().Update("object", "Player", []string{"Orient"}, map[string]interface{}{"Id": load.Id}, load, o.UpdateBack)
}

func (o *Object) UpdateBack(reply *protocol.Message) {
	errcode, err, tag, affected := store.ParseUpdateReply(reply)
	if err != nil {
		o.CoreAPI.LogErr(errcode, err, tag)
		return
	}
	o.CoreAPI.LogInfo("update result:", affected)
}

func (o *Object) LoadAllBack(reply *protocol.Message) {
	var load []*entity.PlayerArchive
	errcode, err, tag := store.ParseFindReply(reply, &load)
	if err != nil {
		o.CoreAPI.LogErr(errcode, err, tag)
		return
	}

	for k, v := range load {
		o.CoreAPI.LogInfo("load result: ", k, v)
	}

	o.store.Client().DeleteByObject("object", "Player", load[len(load)-1], o.DeleteBack)

}

func (o *Object) DeleteBack(reply *protocol.Message) {
	errcode, err, tag, affected := store.ParseDeleteReply(reply)
	if err != nil {
		o.CoreAPI.LogErr(errcode, err, tag)
		return
	}
	o.CoreAPI.LogInfo("delete result:", affected)
}

type GamePlayer struct {
	*gameobject.Role
	*entity.Player
}

func NewGamePlayer() *GamePlayer {
	p := &GamePlayer{}
	p.Role = gameobject.NewRole()
	p.Player = entity.NewPlayer()
	p.Role.SetSpirit(p.Player)
	return p
}

type GamePlayerCreater struct {
}

func (g *GamePlayerCreater) Create() interface{} {
	return NewGamePlayer()
}
