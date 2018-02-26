package main

import (
	"ngengine/core/rpc"
	"ngengine/core/service"
	"ngengine/protocol"
)

var startargs = `{
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
	"LogFile":"test.log",
	"Args": {}
}`

// service
type Database struct {
	service.BaseService
	acc *Account
}

func (d *Database) Prepare(core service.CoreApi) error {
	d.CoreApi = core
	d.acc = &Account{owner: d}
	d.acc.Thread = rpc.NewThread("account", 5, 10)
	core.RegisterRemote("Account", d.acc)
	return nil
}

func (d *Database) Start() error {
	d.CoreApi.Watch("all")
	return nil
}

// rpc
type Account struct {
	*rpc.Thread
	owner *Database
}

func (a *Account) RegisterCallback(s rpc.Servicer) {
	s.RegisterCallback("Login", a.Login)
}

func (a *Account) Login(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	name, _ := m.ReadString()
	pass, _ := m.ReadString()
	mb := rpc.Mailbox{}
	m.ReadObject(&mb)
	a.owner.CoreApi.LogDebug("login:", name, ",pass:", pass)

	a.owner.CoreApi.Mailto(nil, &mb, "Client.Login", LoginResult{"ok"})

	if pass == "123" {
		return 0, protocol.ReplyMessage(protocol.TINY, "ok", mb)
	} else {
		return 0, protocol.ReplyMessage(protocol.TINY, "failed", mb)
	}

}
