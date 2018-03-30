package store

import (
	"fmt"
	"ngengine/core/rpc"
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

// 从数据库中加载一个
func (s *StoreClient) Get(tag string, typ string, condition map[string]interface{}, reply rpc.ReplyCB) error {
	if s.db == nil {
		return fmt.Errorf("database not connected")
	}
	if reply == nil {
		return s.ctx.core.Mailto(nil, s.db, "store.Get", tag, typ, condition)
	}
	return s.ctx.core.MailtoAndCallback(nil, s.db, "store.Get", reply, tag, typ, condition)
}

// 从数据库中加载多个
func (s *StoreClient) Find(tag string, typ string, condition map[string]interface{}, limit int, start int, reply rpc.ReplyCB) error {
	if s.db == nil {
		return fmt.Errorf("database not connected")
	}
	if reply == nil {
		return s.ctx.core.Mailto(nil, s.db, "store.Find", tag, typ, condition, limit, start)
	}
	return s.ctx.core.MailtoAndCallback(nil, s.db, "store.Find", reply, tag, typ, condition, limit, start)
}

// 插入数据
func (s *StoreClient) Insert(tag string, typ string, object interface{}, reply rpc.ReplyCB) error {
	if s.db == nil {
		return fmt.Errorf("database not connected")
	}
	if reply == nil {
		return s.ctx.core.Mailto(nil, s.db, "store.Insert", tag, typ, object)
	}
	return s.ctx.core.MailtoAndCallback(nil, s.db, "store.Insert", reply, tag, typ, object)
}

// 更新数据
func (s *StoreClient) Update(tag string, reply rpc.ReplyCB) error {
	return nil
}

// 删除数据
func (s *StoreClient) Delete(tag string, reply rpc.ReplyCB) error {
	return nil
}

// 原生sql查询
func (s *StoreClient) Query(tag string, sql string, reply rpc.ReplyCB) error {
	return nil
}

// 原生sql执行
func (s *StoreClient) Execute(tag string, sql string, reply rpc.ReplyCB) error {
	return nil
}
