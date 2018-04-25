package login

import (
	"ngengine/core/service"
	"ngengine/module/store"
	"ngengine/share"
)

type LoginModule struct {
	service.Module
	core        service.CoreAPI
	account     *Account
	storeClient *store.StoreClient
}

func New() *LoginModule {
	l := &LoginModule{}
	l.account = &Account{ctx: l}
	return l
}

func (l *LoginModule) Name() string {
	return "Login"
}

func (l *LoginModule) Init(core service.CoreAPI) bool {
	store := core.Module("Store").(*store.StoreModule)
	if store == nil {
		core.LogFatal("need Store module")
		return false
	}
	l.core = core
	l.storeClient = store.Client()
	l.core.Service().AddListener(share.EVENT_READY, l.account.OnDatabaseReady)
	l.core.RegisterHandler("Account", l.account)
	return true
}

// Shut 模块关闭
func (l *LoginModule) Shut() {
	l.core.Service().RemoveListener(share.EVENT_READY, l.account.OnDatabaseReady)
}
