package nest

// Nest 模块

import (
	"ngengine/core/service"
	"ngengine/game/gameobject/entity"
	"ngengine/game/gameobject/models"
	"ngengine/game/nest/session"
	"ngengine/module/object"
	"ngengine/module/store"
	"ngengine/module/timer"
)

type Nest struct {
	service.BaseService
	store   *store.StoreModule
	session *session.SessionModule
	timer   *timer.TimerModule
	factory *object.ObjectModule
}

func (n *Nest) Prepare(core service.CoreAPI) error {
	n.CoreAPI = core
	n.store = store.New()
	n.session = session.New()
	n.timer = timer.New()
	n.factory = object.New()
	return nil
}

func (n *Nest) Init(opt *service.CoreOption) error {
	n.CoreAPI.AddModule(n.store)
	n.CoreAPI.AddModule(n.session)
	n.CoreAPI.AddModule(n.timer)
	n.CoreAPI.AddModule(n.factory)
	n.store.SetMode(store.STORE_CLIENT)
	entity.RegisterToDB(n.store)

	return nil
}

func (n *Nest) Start() error {
	n.BaseService.Start()
	models.Register(n.factory) // 注册gameobjet
	return nil
}
