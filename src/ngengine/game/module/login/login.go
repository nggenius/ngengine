package login

import (
	"ngengine/core/rpc"
	"ngengine/protocol"
)

type Account struct {
	ctx *LoginModule
}

func (a *Account) RegisterCallback(s rpc.Servicer) {
	s.RegisterCallback("Login", a.Login)
}

func (a *Account) Login(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	return 0, nil
}
