package main

import (
	"ngengine/common/event"
	"ngengine/core/service"
	"ngengine/game/gameobject/entity"
	"ngengine/game/module/login"
	"ngengine/module/store"
	"ngengine/share"
)

var startlogin = `{
	"ServId":3,
	"ServType": "login",
	"AdminAddr":"127.0.0.1",
	"AdminPort":12500,
	"ServName": "login_1",
	"ServAddr": "127.0.0.1",
	"ServPort": 0,
	"Expose": true,
	"HostAddr": "0.0.0.0",
	"HostPort": 2002,
	"LogFile":"test1.log",
	"Args": {}
}`

type LoginResult struct {
	Result string
}

// service
type Login struct {
	service.BaseService
	login *login.LoginModule
	store *store.StoreModule
}

func (l *Login) Prepare(core service.CoreAPI) error {
	l.CoreAPI = core
	l.login = login.New()
	l.store = store.New()
	return nil
}

func (l *Login) Init(opt *service.CoreOption) error {
	l.CoreAPI.AddModule(l.store)
	l.store.SetMode(store.STORE_CLIENT)
	entity.RegisterToDB(l.store)
	l.CoreAPI.AddModule(l.login)
	return nil
}

func (l *Login) Start() error {
	l.CoreAPI.Watch("all")
	return nil
}

func (l *Login) OnEvent(e string, args event.EventArgs) {
	switch e {
	case share.EVENT_USER_CONNECT:
		l.CoreAPI.LogDebug("new user")
	case share.EVENT_USER_LOST:
		l.CoreAPI.LogDebug("lost user")
	}
}
