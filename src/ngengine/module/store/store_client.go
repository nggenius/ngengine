package store

import (
	"fmt"
	"ngengine/core/rpc"
	"ngengine/protocol"
)

type StoreClient struct {
	ctx *StoreModule
	db  *rpc.Mailbox
}

func NewStoreClient(ctx *StoreModule) *StoreClient {
	s := &StoreClient{ctx: ctx}
	return s
}

func (s *StoreClient) OnDatabaseReady(evt string, args ...interface{}) {
	srv := s.ctx.core.LookupOneServiceByType("database")
	if srv == nil {
		s.db = nil
		return
	}

	mb := rpc.GetServiceMailbox(srv.Id)
	s.db = &mb
}

func (s *StoreClient) Load() error {
	if s.db == nil {
		return fmt.Errorf("database not connected")
	}

	err := s.ctx.core.MailtoAndCallback(nil, s.db, "store.Load", s.LoadBack)
	return err
}

func (s *StoreClient) LoadBack(reply *protocol.Message) {

}

func (s *StoreClient) Insert() error {
	return nil
}

func (s *StoreClient) Update() error {
	return nil
}

func (s *StoreClient) Delete() error {
	return nil
}
