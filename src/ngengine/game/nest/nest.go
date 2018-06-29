package nest

import (
	"ngengine/core/service"
	"ngengine/game/gameobject/entity"
	"ngengine/game/nest/session"
	"ngengine/module/store"
	"ngengine/module/timer"
)

type Nest struct {
	service.BaseService
	store   *store.StoreModule
	session *session.SessionModule
	timer   *timer.TimerModule
}

func (n *Nest) Prepare(core service.CoreAPI) error {
	n.CoreAPI = core
	n.store = store.New()
	n.session = session.New()
	n.timer = timer.New()
	return nil
}

func (n *Nest) Init(opt *service.CoreOption) error {
	n.CoreAPI.AddModule(n.store)
	n.CoreAPI.AddModule(n.session)
	n.CoreAPI.AddModule(n.timer)
	n.store.SetMode(store.STORE_CLIENT)
	entity.RegisterToDB(n.store)
	return nil
}

func (n *Nest) Start() error {
	n.BaseService.Start()
	return nil
}
