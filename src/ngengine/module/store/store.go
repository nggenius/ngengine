package store

import (
	"errors"
	"ngengine/core/rpc"
	"ngengine/protocol"
	"ngengine/share"

	"github.com/go-xorm/xorm"
)

var (
	ErrNoRows      = errors.New("no row found")
	ErrObject      = errors.New("object type error")
	ErrNoCondition = errors.New("get condition is empty")
)

type Store struct {
	*rpc.Thread
	ctx *StoreModule
}

type getid interface {
	DBId() int64
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
	svr.RegisterCallback("Exec", s.Exec)
}

func (s *Store) Get(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	tag, _ := m.ReadString()
	typ, _ := m.ReadString()
	var condition map[string]interface{}
	if err := m.Read(&condition); err != nil {
		return share.ERR_ARGS_ERROR, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}
	obj := s.ctx.register.Create(typ)
	if obj == nil {
		return share.ERR_ARGS_ERROR, protocol.ReplyMessage(protocol.TINY, tag, ErrObject.Error())
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
		return share.ERR_ARGS_ERROR, protocol.ReplyMessage(protocol.TINY, tag, ErrNoCondition.Error())
	}

	has, err := session.Get(obj)
	if err != nil {
		return share.ERR_STORE_SQL, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}

	if !has {
		return share.ERR_STORE_NOROW, protocol.ReplyMessage(protocol.DEF, tag, ErrNoRows.Error())
	}

	return share.ERR_REPLY_SUCCEED, protocol.ReplyMessage(protocol.DEF, tag, obj)
}

func (s *Store) Find(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	tag, _ := m.ReadString()
	typ, _ := m.ReadString()
	var condition map[string]interface{}
	if err := m.Read(&condition); err != nil {
		return share.ERR_ARGS_ERROR, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}
	var limit, start int
	m.Read(&limit)
	m.Read(&start)
	obj := s.ctx.register.CreateSlice(typ)
	if obj == nil {
		return share.ERR_ARGS_ERROR, protocol.ReplyMessage(protocol.TINY, tag, ErrObject.Error())
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
		session = s.ctx.sql.orm.NewSession()
	}

	if limit != 0 || start != 0 {
		session = session.Limit(limit, start)
	}

	if err := session.Find(obj); err != nil {
		return share.ERR_STORE_SQL, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}

	return share.ERR_REPLY_SUCCEED, protocol.ReplyMessage(protocol.DEF, tag, obj)
}

func (s *Store) Insert(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	tag, _ := m.ReadString()
	typ, _ := m.ReadString()
	obj := s.ctx.register.Create(typ)
	if err := m.ReadObject(obj); err != nil {
		return share.ERR_ARGS_ERROR, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}

	affected, err := s.ctx.sql.orm.Insert(obj)
	if err != nil {
		return share.ERR_STORE_SQL, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}

	var id int64
	if get, ok := obj.(getid); ok {
		id = get.DBId()
	}
	return share.ERR_REPLY_SUCCEED, protocol.ReplyMessage(protocol.TINY, tag, affected, id)
}

func (s *Store) Update(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	tag, _ := m.ReadString()
	typ, _ := m.ReadString()
	var cols []string
	if err := m.Read(&cols); err != nil {
		return share.ERR_ARGS_ERROR, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}
	var condition map[string]interface{}
	if err := m.Read(&condition); err != nil {
		return share.ERR_ARGS_ERROR, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}
	obj := s.ctx.register.Create(typ)
	if err := m.ReadObject(obj); err != nil {
		return share.ERR_ARGS_ERROR, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}

	session := s.ctx.sql.orm.NewSession()
	if len(cols) > 0 {
		session = session.Cols(cols...)
	}

	var affected int64
	var err error
	if condition == nil || len(condition) == 0 {
		affected, err = session.Update(obj)
	} else {
		affected, err = session.Update(obj, condition)
	}

	if err != nil {
		return share.ERR_STORE_SQL, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}
	return share.ERR_REPLY_SUCCEED, protocol.ReplyMessage(protocol.TINY, tag, affected)
}

func (s *Store) Delete(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	tag, _ := m.ReadString()
	typ, _ := m.ReadString()
	obj := s.ctx.register.Create(typ)
	if obj == nil {
		return share.ERR_ARGS_ERROR, protocol.ReplyMessage(protocol.TINY, tag, ErrObject.Error())
	}

	if err := m.ReadObject(obj); err != nil {
		return share.ERR_ARGS_ERROR, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}

	affected, err := s.ctx.sql.orm.Delete(obj)
	if err != nil {
		return share.ERR_STORE_SQL, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}
	return share.ERR_REPLY_SUCCEED, protocol.ReplyMessage(protocol.TINY, tag, affected)
}

func (s *Store) Query(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	tag, _ := m.ReadString()
	sql, _ := m.ReadString()
	var args []interface{}
	if err := m.Read(&args); err != nil {
		return share.ERR_ARGS_ERROR, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}

	sqlorargs := make([]interface{}, 0, 1+len(args))
	sqlorargs = append(sqlorargs, sql)
	if len(args) > 0 {
		sqlorargs = append(sqlorargs, args...)
	}
	result, err := s.ctx.sql.orm.Query(sqlorargs...)
	if err != nil {
		return share.ERR_STORE_SQL, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}
	return share.ERR_REPLY_SUCCEED, protocol.ReplyMessage(protocol.DEF, tag, result)
}

func (s *Store) Exec(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	tag, _ := m.ReadString()
	sql, _ := m.ReadString()
	var args []interface{}
	if err := m.Read(&args); err != nil {
		return 1, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}
	res, err := s.ctx.sql.orm.Exec(sql, args...)
	if err != nil {
		return share.ERR_STORE_SQL, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return share.ERR_STORE_ERROR, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}
	return share.ERR_REPLY_SUCCEED, protocol.ReplyMessage(protocol.TINY, tag, affected)
}
