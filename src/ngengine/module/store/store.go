package store

import (
	"ngengine/core/rpc"
	"ngengine/protocol"

	"github.com/go-xorm/xorm"
)

type Store struct {
	*rpc.Thread
	ctx *StoreModule
}

func NewStore(ctx *StoreModule) *Store {
	s := &Store{}
	s.ctx = ctx
	s.Thread = rpc.NewThread("store", 4, 10)
	return s
}

func (s *Store) RegisterCallback(svr rpc.Servicer) {
	svr.RegisterCallback("Get", s.Get)
	svr.RegisterCallback("Find", s.Find)
	svr.RegisterCallback("Insert", s.Insert)
	svr.RegisterCallback("Update", s.Update)
	svr.RegisterCallback("Delete", s.Delete)
	svr.RegisterCallback("Query", s.Query)
	svr.RegisterCallback("Execute", s.Execute)
}

func (s *Store) Get(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	tag, _ := m.ReadString()
	typ, _ := m.ReadString()
	var condition map[string]interface{}
	if err := m.Read(&condition); err != nil {
		return 1, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}
	obj := s.ctx.register.Create(typ)
	if obj == nil {
		return 1, protocol.ReplyMessage(protocol.TINY, tag, "object create failed")
	}

	var session *xorm.Session
	for k, v := range condition {
		if session == nil {
			session = s.ctx.sql.orm.Where(k, v)
			continue
		}
		session = session.And(k, v)
	}
	if session == nil {
		return 1, protocol.ReplyMessage(protocol.TINY, tag, "condition is nil")
	}
	_, err := session.Get(obj)
	if err != nil {
		return 1, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}

	return 0, protocol.ReplyMessage(protocol.DEF, tag, obj)
}

func (s *Store) Find(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	tag, _ := m.ReadString()
	typ, _ := m.ReadString()
	var condition map[string]interface{}
	if err := m.Read(&condition); err != nil {
		return 1, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}
	var limit, start int
	m.Read(&limit)
	m.Read(&start)
	obj := s.ctx.register.CreateSlice(typ)
	if obj == nil {
		return 1, protocol.ReplyMessage(protocol.TINY, tag, "object create failed")
	}

	var session *xorm.Session
	for k, v := range condition {
		if session == nil {
			session = s.ctx.sql.orm.Where(k, v)
			continue
		}
		session = session.And(k, v)
	}

	if limit != 0 || start != 0 {
		if session == nil {
			session = s.ctx.sql.orm.Limit(limit, start)
		} else {
			session = session.Limit(limit, start)
		}
	}
	var err error
	if session == nil {
		err = s.ctx.sql.orm.Find(obj)
	} else {
		err = session.Find(obj)
	}

	if err != nil {
		return 1, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}
	return 0, protocol.ReplyMessage(protocol.DEF, tag, obj)
}

func (s *Store) Insert(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	tag, _ := m.ReadString()
	typ, _ := m.ReadString()
	obj := s.ctx.register.Create(typ)
	if err := m.ReadObject(obj); err != nil {
		return 1, protocol.ReplyMessage(128, tag, err.Error())
	}

	id, err := s.ctx.sql.orm.Insert(obj)
	if err != nil {
		return 1, protocol.ReplyMessage(128, tag, err.Error())
	}

	return 0, protocol.ReplyMessage(128, tag, id)
}

func (s *Store) Update(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	return 0, nil
}

func (s *Store) Delete(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	return 0, nil
}

func (s *Store) Query(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	return 0, nil
}

func (s *Store) Execute(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	return 0, nil
}
