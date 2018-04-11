package main

import (
	"fmt"
	"ngengine/common/event"
	"ngengine/core/rpc"
	"ngengine/core/service"
	"ngengine/module/timer"
	"ngengine/protocol"
	"ngengine/share"
)

var startargs2 = `{
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
}

func (l *Login) Prepare(core service.CoreAPI) error {
	l.CoreAPI = core
	core.RegisterHandler("User", &User{l})
	core.AddModule(&ModuleTest{})
	core.AddModule(&ModuleTest2{})
	core.AddModule(&timer.TimerModule{})
	return nil
}

func (l *Login) Start() error {
	l.CoreAPI.Watch("all")

	timer := l.Module("TimerModule").(*timer.TimerModule)
	id := timer.AddTimer(1000, 123, l.RepeatCallBack)
	id2 := timer.AddCountTimer(10, 3000, 999, l.CountCallBack)
	l.CoreAPI.LogDebugf("id:%d   id2:%d", id, id2)
	return nil
}

func (l *Login) RepeatCallBack(id int64, count int, args interface{}) {
	fmt.Print(id, count, args)
}

func (l *Login) CountCallBack(id int64, count int, args interface{}) {
	fmt.Print(id, count, args)

}

func (l *Login) OnEvent(e string, args event.EventArgs) {
	switch e {
	case share.EVENT_USER_CONNECT:
		l.CoreAPI.LogDebug("new user")
	case share.EVENT_USER_LOST:
		l.CoreAPI.LogDebug("lost user")
	}
}

func (l *Login) OnReply(reply *protocol.Message) {
	m := protocol.NewMessageReader(reply)
	r, _ := m.ReadString()
	l.CoreAPI.LogDebug("login result:", r)
}

type User struct {
	owner *Login
}

func (u *User) RegisterCallback(s rpc.Servicer) {
	s.RegisterCallback("Login", u.Login)
}

func (u *User) Login(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {

	srv := u.owner.CoreAPI.LookupOneServiceByType("database")
	if srv == nil {
		u.owner.CoreAPI.LogErr("database not found")
		return 0, nil
	}

	dest := rpc.GetServiceMailbox(srv.Id)
	u.owner.CoreAPI.LogDebug(mailbox, "request login")
	//err := u.owner.CoreApi.MailtoAndCallback(nil, &dest, "Account.Login", u.OnReply, "sll", "123", mailbox)
	err := u.owner.CoreAPI.Mailto(nil, &dest, "Account.Login", "sll", "123", mailbox)
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
	if err := u.owner.CoreAPI.Mailto(nil, mb, "Client.Login", LoginResult{r}); err != nil {
		u.owner.CoreAPI.LogErr("send to client failed", mb)
		return
	}
	u.owner.CoreAPI.LogDebug("login result:", r)
}
