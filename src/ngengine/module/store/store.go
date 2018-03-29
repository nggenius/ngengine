package store

import (
	"ngengine/core/rpc"
	"ngengine/protocol"
)

type Store struct {
	*rpc.Thread
	ctx *StoreModule
}

func NewStore(ctx *StoreModule) *Store {
	s := &Store{}
	s.ctx = ctx
	s.Thread = rpc.NewThread("store", 16, 10)
	return s
}

func (s *Store) RegisterCallback(svr rpc.Servicer) {
	svr.RegisterCallback("Load", s.Load)
	svr.RegisterCallback("Insert", s.Insert)
	svr.RegisterCallback("Update", s.Update)
	svr.RegisterCallback("Delete", s.Delete)
}

func (s *Store) Load(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	return 0, nil
}

func (s *Store) Insert(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	return 0, nil
}

func (s *Store) Update(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	return 0, nil
}

func (s *Store) Delete(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	return 0, nil
}
