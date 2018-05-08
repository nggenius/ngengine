package session

import (
	"ngengine/core/rpc"
	"ngengine/protocol"
	"ngengine/protocol/proto/c2s"
)

// 客户端连接点
type proxy struct {
	ctx *SessionModule
}

func NewProxy(ctx *SessionModule) *proxy {
	p := &proxy{}
	p.ctx = ctx
	return p
}

func (p *proxy) RegisterCallback(s rpc.Servicer) {
	s.RegisterCallback("Login", p.Login)
}

// token 登录
func (p *proxy) Login(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	login := &c2s.LoginNest{}
	if err := p.ctx.core.ParseProto(msg, login); err != nil {
		p.ctx.core.LogErr("login parse error,", err)
		return 0, nil
	}

	session := p.ctx.FindSession(mailbox.Id())
	if session == nil {
		p.ctx.core.LogErr("session not found")
		return 0, nil
	}

	session.Account = login.Account
	session.Mailbox = &mailbox
	session.Dispatch(LOGIN, login.Token)
	return 0, nil
}
