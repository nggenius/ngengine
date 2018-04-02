package main

import (
	"ngengine/core/service"
	"ngengine/module/object/entity"
	"ngengine/module/store"
)

var dbargs = `{
	"ServId":1,
	"ServType": "database",
	"AdminAddr":"127.0.0.1",
	"AdminPort":12500,
	"ServName": "db_1",
	"ServAddr": "127.0.0.1",
	"ServPort": 0,
	"Expose": false,
	"HostAddr": "",
	"HostPort": 0,
	"LogFile":"db.log",
	"Args": {
		"db":"mysql",
		"datasource":"sa:abc@tcp(192.168.1.52:3306)/test?charset=utf8"
	}
}`

// service
type Database struct {
	service.BaseService
	store *store.StoreModule
}

func (d *Database) Prepare(core service.CoreApi) error {
	d.CoreApi = core
	d.store = store.New()
	return nil
}

func (d *Database) Init(opt *service.CoreOption) error {
	d.CoreApi.AddModule(d.store)
	d.store.SetMode(store.STORE_SERVER)
	d.store.Register("Player", &entity.PlayerArchiveCreater{})
	return nil
}

func (d *Database) Start() error {
	d.CoreApi.Watch("all")
	return nil
}
