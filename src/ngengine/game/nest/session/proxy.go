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
	s.RegisterCallback("CreateRole", p.CreateRole)
	s.RegisterCallback("ChooseRole", p.ChooseRole)
	s.RegisterCallback("DeleteRole", p.DeleteRole)
}

// token 登录
func (p *proxy) Login(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	login := &c2s.LoginNest{}
	if err := p.ctx.Core.ParseProto(msg, login); err != nil {
		p.ctx.Core.LogErr("login parse error,", err)
		return 0, nil
	}

	session := p.ctx.FindSession(sender.Id())
	if session == nil {
		p.ctx.Core.LogErr("session not found")
		return 0, nil
	}

	session.Account = login.Account
	session.Mailbox = &sender
	session.Dispatch(ELOGIN, login.Token)
	return 0, nil
}

// 创建角色
func (p *proxy) CreateRole(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	var create c2s.CreateRole
	if err := p.ctx.Core.ParseProto(msg, &create); err != nil {
		p.ctx.Core.LogErr("create parse error,", err)
		return 0, nil
	}

	session := p.ctx.FindSession(sender.Id())
	if session == nil {
		p.ctx.Core.LogErr("session not found")
		return 0, nil
	}

	session.Dispatch(ECREATE, create)

	return 0, nil
}

// 选择角色
func (p *proxy) ChooseRole(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	var choose c2s.ChooseRole
	if err := p.ctx.Core.ParseProto(msg, &choose); err != nil {
		p.ctx.Core.LogErr("choose parse error,", err)
		return 0, nil
	}

	session := p.ctx.FindSession(sender.Id())
	if session == nil {
		p.ctx.Core.LogErr("session not found")
		return 0, nil
	}

	session.Dispatch(ECHOOSE, choose)

	return 0, nil
}

// 删除角色
func (p *proxy) DeleteRole(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	var d c2s.DeleteRole
	if err := p.ctx.Core.ParseProto(msg, &d); err != nil {
		p.ctx.Core.LogErr("delete parse error, ", err)
		return 0, nil
	}

	session := p.ctx.FindSession(sender.Id())
	if session == nil {
		p.ctx.Core.LogErr("session not found")
		return 0, nil
	}

	session.Dispatch(EDELETE, d)
	return 0, nil
}
