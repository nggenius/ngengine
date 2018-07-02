package session

import (
	"ngengine/common/fsm"
	"ngengine/core/rpc"
	"ngengine/game/gameobject/entity/inner"
	"ngengine/protocol/proto/c2s"
	"ngengine/protocol/proto/s2c"
)

type SessionDB map[uint64]*Session

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
	s.FSM.Register("createrole", &createrole{owner: s})
	s.FSM.Register("chooserole", &chooserole{owner: s})
	s.FSM.Register("online", &online{owner: s})
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
		return true
	}
	return false
}

func (s *Session) QueryRoleInfo() bool {
	if err := s.ctx.account.requestRoleInfo(s); err == nil {
		return true
	}
	return false
}

func (s *Session) OnRoleInfo(role []*inner.Role) {
	s.ctx.core.LogDebug("role info", role)
	roles := &s2c.RoleInfo{}
	roles.Roles = make([]s2c.Role, 0, len(role))
	for k := range role {
		r := s2c.Role{}
		r.RoleId = role[k].RoleId
		r.Index = role[k].Index
		r.Name = role[k].Account
		roles.Roles = append(roles.Roles, r)
	}

	s.ctx.core.Mailto(nil, s.Mailbox, "Account.Roles", roles)
}

func (s *Session) CreateRole(info c2s.CreateRole) error {
	return s.ctx.account.CreateRole(s, info)
}

func (s *Session) ChooseRole(info c2s.ChooseRole) error {
	return s.ctx.account.ChooseRole(s, info)
}

func (s *Session) Error(errcode int32) {
	err := s2c.Error{}
	err.ErrCode = errcode
	s.ctx.core.Mailto(nil, s.Mailbox, "system.Error", &err)
}
