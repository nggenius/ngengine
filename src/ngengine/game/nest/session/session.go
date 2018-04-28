package session

import (
	"ngengine/common/fsm"
	"ngengine/core/rpc"
	"ngengine/game/gameobject/entity/inner"
	"ngengine/protocol/proto/s2c"
)

type SessionDB map[uint64]*Session

const (
	NONE      = iota
	TIMER     // 1秒钟的定时器
	BREAK     // 客户端断开连接
	LOGIN     // 客户端登录
	ROLE_INFO // 角色列表
)

type idlestate struct {
	fsm.Default
	owner *Session
	Idle  int32
}

func (s *idlestate) Handle(event int, param interface{}) string {
	switch event {
	case LOGIN:
		token := param.(string)
		if s.owner.ValidToken(token) {
			return "logged"
		}
		// 验证失败直接踢下线
		s.owner.Break()
		return ""
	case TIMER:
		s.Idle++
		if s.Idle > 60 {
			s.owner.Break()
			return ""
		}
	case BREAK:
		s.owner.DestroySelf()
		return fsm.STOP
	}
	return ""
}

type logged struct {
	fsm.Default
	owner *Session
	Idle  int32
}

func (s *logged) Handle(event int, param interface{}) string {
	switch event {
	case TIMER:
		s.Idle++
		if s.Idle > 60 {
			s.owner.Break()
			return ""
		}
	case BREAK:
		s.owner.DestroySelf()
		return fsm.STOP
	case ROLE_INFO:
		args := param.([2]interface{})
		errcode := args[0].(int32)
		roles := args[1].([]*inner.Role)
		if errcode != 0 {
			s.owner.Error(errcode)
			return ""
		}

		s.owner.OnRoleInfo(roles)
		s.Idle = 0
	}
	return ""
}

type Session struct {
	*fsm.FSM
	ctx     *SessionModule
	id      uint64
	Account string
	Mailbox *rpc.Mailbox
	delete  bool
}

func NewSession(id uint64, ctx *SessionModule) *Session {
	s := &Session{}
	s.ctx = ctx
	s.id = id
	s.FSM = fsm.NewFSM()
	s.FSM.Register("idle", &idlestate{owner: s})
	s.FSM.Register("logged", &logged{owner: s})
	s.FSM.Start("idle")
	return s
}

func (s *Session) DestroySelf() {
	s.delete = true
	s.ctx.deleted.PushBack(s.id)
}

func (s *Session) Break() {
	s.ctx.core.Break(s.id)
}

func (s *Session) ValidToken(token string) bool {
	if s.ctx.cache.Valid(s.Account, token) {
		if err := s.ctx.account.requestRoleInfo(s); err == nil {
			return true
		}
	}
	return false
}

func (s *Session) OnRoleInfo(role []*inner.Role) {
	s.ctx.core.LogDebug("role info", role)
	roles := &s2c.RoleInfo{}
	roles.Roles = make([]s2c.Role, 0, len(role))
	for k := range role {
		r := s2c.Role{}
		r.Index = role[k].Index
		r.Name = role[k].Account
		roles.Roles = append(roles.Roles, r)
	}

	s.ctx.core.Mailto(nil, s.Mailbox, "Account.Roles", roles)
}

func (s *Session) Error(errcode int32) {

}
