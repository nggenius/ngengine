package main

import (
	"fmt"
	"ngengine/common/event"
	"ngengine/core/rpc"
	"ngengine/core/service"
	"ngengine/protocol"
	"ngengine/share"
)

var startargs2 = `{
	"ServId":2,
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
}

func (l *Login) Prepare(core service.CoreApi) error {
	l.CoreApi = core
	core.RegisterHandler("User", &User{l})
	core.AddModule(&ModuleTest{})
	core.AddModule(&ModuleTest2{})
	return nil
}

func (l *Login) Start() error {
	l.CoreApi.Watch("all")
	return nil
}

func (l *Login) OnEvent(e string, args event.EventArgs) {
	switch e {
	case share.EVENT_USER_CONNECT:
		l.CoreApi.LogDebug("new user")
	case share.EVENT_USER_LOST:
		l.CoreApi.LogDebug("lost user")
	}
}

func (l *Login) OnReply(reply *protocol.Message) {
	m := protocol.NewMessageReader(reply)
	r, _ := m.ReadString()
	l.CoreApi.LogDebug("login result:", r)
}

type User struct {
	owner *Login
}

func (u *User) RegisterCallback(s rpc.Servicer) {
	s.RegisterCallback("Login", u.Login)
}

func (u *User) Login(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {

	srv := u.owner.CoreApi.LookupOneServiceByType("database")
	if srv == nil {
		u.owner.CoreApi.LogErr("database not found")
		return 0, nil
	}

	dest := rpc.GetServiceMailbox(srv.Id)
	u.owner.CoreApi.LogDebug(mailbox, "request login")
	//err := u.owner.CoreApi.MailtoAndCallback(nil, &dest, "Account.Login", u.OnReply, "sll", "123", mailbox)
	err := u.owner.CoreApi.Mailto(nil, &dest, "Account.Login", "sll", "123", mailbox)
	if err != nil {
		fmt.Println(err)
	}

	return 0, nil
}

func (u *User) OnReply(reply *protocol.Message) {
	m := protocol.NewMessageReader(reply)
	r, _ := m.ReadString()
	mb := &rpc.Mailbox{}
	m.ReadObject(mb)
	if err := u.owner.CoreApi.Mailto(nil, mb, "Client.Login", LoginResult{r}); err != nil {
		u.owner.CoreApi.LogErr("send to client failed", mb)
		return
	}
	u.owner.CoreApi.LogDebug("login result:", r)
}
