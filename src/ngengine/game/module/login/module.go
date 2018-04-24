package login

import "ngengine/core/service"

type LoginModule struct {
	service.Module
	core    service.CoreAPI
	account *Account
}

func New() *LoginModule {
	l := &LoginModule{}
	l.account = &Account{ctx: l}
	return l
}

func (l *LoginModule) Init(core service.CoreAPI) bool {
	l.core = core
	l.core.RegisterHandler("Account", l.account)
	return true
}
