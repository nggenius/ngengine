package store

import (
	"ngengine/core/service"
	"ngengine/game/gameobject/entity"
	"ngengine/game/store/extension"
	"ngengine/module/store"
)

type Store struct {
	service.BaseService
	store *store.StoreModule
	role  *extension.Role
}

func (d *Store) Prepare(core service.CoreAPI) error {
	d.CoreAPI = core
	d.store = store.New()
	d.role = extension.NewRole(d.CoreAPI, d.store)
	return nil
}

func (d *Store) Init(opt *service.CoreOption) error {
	d.CoreAPI.AddModule(d.store)
	d.store.SetMode(store.STORE_SERVER)
	d.store.Extend("role", d.role)
	entity.RegisterToDB(d.store)
	return nil
}

func (d *Store) Start() error {
	d.CoreAPI.Watch("all")
	return nil
}
