package login

import (
	"ngengine/common/fsm"
	"ngengine/core/rpc"
	"ngengine/game/gameobject/entity/inner"
	"ngengine/protocol/proto/c2s"
	"ngengine/protocol/proto/s2c"
	"ngengine/share"
)

const (
	NONE         = iota
	TIMER        // 1秒钟的定时器
	BREAK        // 客户端断开连接
	LOGIN        // 客户端登录
	LOGIN_RESULT // 登录结果
	NEST_RESULT  // nest 登录结果
)

type Idle struct {
	fsm.Default
	owner *Session
	idle  int32
}

func (s *Idle) Handle(event int, param interface{}) string {
	switch event {
	case LOGIN:
		s.owner.Login(param.(*c2s.Login))
		return "Logging"
	case TIMER:
		s.idle++
		if s.idle > 60 {
			s.owner.Break()
			return ""
		}
	case BREAK:
		s.owner.DestroySelf()
		return fsm.STOP
	}
	return ""
}

type Logging struct {
	fsm.Default
	owner *Session
	idle  int32
}

func (l *Logging) Handle(event int, param interface{}) string {
	switch event {
	case LOGIN_RESULT:
		args := param.([2]interface{})
		if l.owner.LoginResult(args[0].(int32), args[1].(*inner.Account)) {
			return "Logged"
		}
	case TIMER:
		l.idle++
		if l.idle > 60 {
			l.owner.Error(share.S2C_ERR_SERVICE_INVALID)
			l.owner.Break()
			return ""
		}
	case BREAK:
		l.owner.DestroySelf()
		return fsm.STOP
	}
	return ""
}

type Logged struct {
	fsm.Default
	owner *Session
	idle  int32
}

func (l *Logged) Handle(event int, param interface{}) string {
	switch event {
	case NEST_RESULT:
		args := param.([2]interface{})
		if l.owner.NestResult(args[0].(int32), args[1].(string)) {
			l.idle = 0 // 1 分钟后退出
			return ""
		}
		return "Start" //重新登录
	case TIMER:
		l.idle++
		if l.idle > 60 {
			l.owner.Error(share.S2C_ERR_SERVICE_INVALID)
			l.owner.Break()
			return ""
		}
	case BREAK:
		l.owner.DestroySelf()
		return fsm.STOP
	}
	return ""
}

type Session struct {
	*fsm.FSM
	ctx     *LoginModule
	id      uint64
	Account string
	Mailbox *rpc.Mailbox
	nest    share.ServiceId
	delete  bool
}

func NewSession(id uint64, ctx *LoginModule) *Session {
	c := &Session{}
	c.ctx = ctx
	c.id = id
	c.FSM = fsm.NewFSM()
	c.FSM.Register("Idle", &Idle{owner: c})
	c.FSM.Register("Logging", &Logging{owner: c})
	c.FSM.Register("Logged", &Logged{owner: c})
	c.FSM.Start("Idle")
	return c
}

func (c *Session) DestroySelf() {
	c.delete = true
	c.ctx.deleted.PushBack(c.id)
}

func (c *Session) SetAccount(acc string) {
	c.Account = acc
}

func (c *Session) SetMailbox(mb *rpc.Mailbox) {
	c.Mailbox = mb
}

func (c *Session) Login(login *c2s.Login) {
	c.ctx.account.sendLogin(c, login)
}

func (c *Session) LoginResult(errcode int32, accinfo *inner.Account) bool {
	if errcode != 0 {
		c.Error(errcode)
		return false
	}

	if accinfo.Id != 0 {
		srv := c.ctx.account.findNest(c)
		if srv != nil {
			c.nest = srv.Id
			return true
		}
		return false
	}
	c.Error(share.S2C_ERR_NAME_PASS)
	return false
}

func (c *Session) NestResult(errcode int32, token string) bool {
	if errcode != 0 {
		c.Error(errcode)
		return false
	}

	srv := c.ctx.core.LookupService(c.nest)
	if srv == nil {
		c.Error(share.S2C_ERR_SERVICE_INVALID)
		return false
	}

	nest := &s2c.NestInfo{}
	nest.Addr = srv.OuterAddr
	nest.Port = int32(srv.OuterPort)
	nest.Token = token

	if err := c.ctx.core.Mailto(nil, c.Mailbox, "Login.Nest", nest); err != nil {
		return false
	}

	return true
}

func (c *Session) Break() {
	c.ctx.core.Break(c.id)
}

func (c *Session) Error(err int32) {
	result := s2c.Error{}
	result.ErrCode = err
	c.ctx.core.Mailto(nil, c.Mailbox, "system.Error", result)
}
