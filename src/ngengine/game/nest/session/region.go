package session

import (
	"ngengine/core/rpc"
	"ngengine/protocol"
)

type RoleRegion struct {
	ctx *SessionModule
}

func NewRoleRegion(m *SessionModule) *RoleRegion {
	s := new(RoleRegion)
	s.ctx = m
	return s
}

func (s *RoleRegion) RegisterCallback(srv rpc.Servicer) {
	//srv.RegisterCallback("Method", s.Method)
}

func (s *RoleRegion) Prototype(src rpc.Mailbox, dest rpc.Mailbox, msg *protocol.Message) (int32, *protocol.Message) {
	return 0, nil
}
