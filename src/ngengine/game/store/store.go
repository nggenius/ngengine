package store

import (
	"ngengine/core/service"
	"ngengine/game/gameobject/entity"
	"ngengine/module/store"
)

type Store struct {
	service.BaseService
	store *store.StoreModule
}

func (d *Store) Prepare(core service.CoreAPI) error {
	d.CoreAPI = core
	d.store = store.New()
	return nil
}

func (d *Store) Init(opt *service.CoreOption) error {
	d.CoreAPI.AddModule(d.store)
	d.store.SetMode(store.STORE_SERVER)
	entity.RegisterToDB(d.store)
	return nil
}

func (d *Store) Start() error {
	d.CoreAPI.Watch("all")
	return nil
}
