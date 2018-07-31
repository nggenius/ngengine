package store

import (
	"errors"
	"ngengine/core/rpc"
	"ngengine/protocol"
	"ngengine/share"
)

var (
	ErrNoRows      = errors.New("no row found")
	ErrObject      = errors.New("object type error")
	ErrNoCondition = errors.New("get condition is empty")
)

type Store struct {
	*rpc.Thread
	ctx       *StoreModule
	extension map[string]Extension
}

type getsetid interface {
	DBId() int64
	SetId(val int64)
}

type Extension interface {
	RegisterCallback(svr rpc.Servicer)
}

func NewStore(ctx *StoreModule) *Store {
	s := &Store{}
	s.ctx = ctx
	s.Thread = rpc.NewThread("store", 4, 10)
	s.extension = make(map[string]Extension)
	return s
}

func (s *Store) AddExtension(name string, ext Extension) {
	s.extension[name] = ext
}

func (s *Store) RegisterCallback(svr rpc.Servicer) {
	svr.RegisterCallback("Get", s.Get)
	svr.RegisterCallback("Find", s.Find)
	svr.RegisterCallback("Insert", s.Insert)
	svr.RegisterCallback("MultiInsert", s.MultiInsert)
	svr.RegisterCallback("Update", s.Update)
	svr.RegisterCallback("MultiUpdate", s.Update)
	svr.RegisterCallback("Delete", s.Delete)
	svr.RegisterCallback("Delete2", s.Delete2)
	svr.RegisterCallback("Delete3", s.Delete3)
	svr.RegisterCallback("Query", s.Query)
	svr.RegisterCallback("Exec", s.Exec)
	for k, v := range s.extension {
		v.RegisterCallback(svr)
		s.ctx.Core.LogInfo("register extension ", k)
	}
}

func (s *Store) Get(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	tag, _ := m.ReadString()
	typ, _ := m.ReadString()
	var condition map[string]interface{}
	if err := m.Read(&condition); err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, err.Error(), tag)
	}
	obj := s.ctx.register.Create(typ)
	if obj == nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, ErrObject.Error(), tag)
	}

	session := s.ctx.sql.Session()
	defer session.Close()
	first := true
	for k, v := range condition {
		if first {
			session.Where(k, v)
			first = false
			continue
		}
		session.And(k, v)
	}
	if first {
		return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, ErrNoCondition.Error(), tag)
	}

	has, err := session.Get(obj)
	if err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_STORE_SQL, err.Error(), tag)
	}

	if !has {
		return protocol.ReplyError(protocol.DEF, share.ERR_STORE_NOROW, ErrNoRows.Error(), tag)
	}

	return protocol.Reply(protocol.DEF, tag, obj)
}

func (s *Store) Find(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	tag, _ := m.ReadString()
	typ, _ := m.ReadString()
	var condition map[string]interface{}
	if err := m.Read(&condition); err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, err.Error(), tag)
	}
	var limit, start int
	m.Read(&limit)
	m.Read(&start)
	obj := s.ctx.register.CreateSlice(typ)
	if obj == nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, ErrObject.Error(), tag)
	}

	session := s.ctx.sql.Session()
	defer session.Close()
	first := true
	for k, v := range condition {
		if first {
			session.Where(k, v)
			continue
		}
		session.And(k, v)
	}

	if limit != 0 || start != 0 {
		session.Limit(limit, start)
	}

	if err := session.Find(obj); err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_STORE_SQL, err.Error(), tag)
	}

	return protocol.Reply(protocol.DEF, tag, obj)
}

func (s *Store) Insert(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	tag, _ := m.ReadString()
	typ, _ := m.ReadString()
	obj := s.ctx.register.Create(typ)
	if err := m.ReadObject(obj); err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, err.Error(), tag)
	}

	affected, err := s.ctx.sql.orm.Insert(obj)
	if err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_STORE_SQL, err.Error(), tag)
	}

	var id int64
	if get, ok := obj.(getsetid); ok {
		id = get.DBId()
	}
	return protocol.Reply(protocol.TINY, tag, affected, id)
}

func (s *Store) MultiInsert(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	tag, _ := m.ReadString()
	var typ []string
	if err := m.Read(&typ); err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, err.Error(), tag)
	}

	var object []interface{}
	for k := range typ {
		obj := s.ctx.register.Create(typ[k])
		if err := m.ReadObject(obj); err != nil {
			return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, err.Error(), tag)
		}
		object = append(object, obj)
	}

	session := s.ctx.sql.Session()
	defer session.Close()
	session.Begin()

	for k := range object {
		_, err := session.Insert(object[k])
		if err != nil {
			return protocol.ReplyError(protocol.TINY, share.ERR_STORE_SQL, err.Error(), tag)
		}
	}

	session.Commit()

	return protocol.Reply(protocol.TINY, tag)
}

func (s *Store) Update(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	tag, _ := m.ReadString()
	typ, _ := m.ReadString()
	var cols []string
	if err := m.Read(&cols); err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, err.Error(), tag)
	}
	var condition map[string]interface{}
	if err := m.Read(&condition); err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, err.Error(), tag)
	}
	obj := s.ctx.register.Create(typ)
	if err := m.ReadObject(obj); err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, err.Error(), tag)
	}

	session := s.ctx.sql.Session()
	defer session.Close()
	if len(cols) > 0 {
		session.Cols(cols...)
	}

	var affected int64
	var err error
	if condition == nil || len(condition) == 0 {
		affected, err = session.Update(obj)
	} else {
		affected, err = session.Update(obj, condition)
	}

	if err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_STORE_SQL, err.Error(), tag)
	}
	return protocol.Reply(protocol.TINY, tag, affected)
}

func (s *Store) MultiUpdate(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	tag, _ := m.ReadString()

	var typ []string
	if err := m.Read(&typ); err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, err.Error(), tag)
	}

	var object []interface{}
	for k := range typ {
		obj := s.ctx.register.Create(typ[k])
		if err := m.ReadObject(obj); err != nil {
			return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, err.Error(), tag)
		}
		object = append(object, obj)
	}

	session := s.ctx.sql.Session()
	defer session.Close()
	session.Begin()

	var affected int64
	for k := range object {
		aff, err := session.Update(object[k])
		if err != nil {
			return protocol.ReplyError(protocol.TINY, share.ERR_STORE_SQL, err.Error(), tag)
		}
		affected += aff
	}

	session.Commit()

	return protocol.Reply(protocol.TINY, tag, affected)
}

func (s *Store) Delete(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	tag, _ := m.ReadString()
	typ, _ := m.ReadString()
	obj := s.ctx.register.Create(typ)
	if obj == nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, ErrObject.Error(), tag)
	}

	id, err := m.ReadInt64()
	if err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, err.Error(), tag)
	}

	if set, ok := obj.(getsetid); ok {
		set.SetId(id)
	}
	affected, err := s.ctx.sql.orm.Delete(obj)
	if err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_STORE_SQL, err.Error(), tag)
	}
	return protocol.Reply(protocol.TINY, tag, affected)
}

func (s *Store) Delete2(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	tag, _ := m.ReadString()
	typ, _ := m.ReadString()
	obj := s.ctx.register.Create(typ)
	if obj == nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, ErrObject.Error(), tag)
	}

	if err := m.ReadObject(obj); err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, err.Error(), tag)
	}

	affected, err := s.ctx.sql.orm.Delete(obj)
	if err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_STORE_SQL, err.Error(), tag)
	}
	return protocol.Reply(protocol.TINY, tag, affected)
}

func (s *Store) Delete3(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	tag, _ := m.ReadString()

	var typ []string
	if err := m.Read(&typ); err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, err.Error(), tag)
	}

	var ids []int64
	if err := m.Read(&ids); err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, err.Error(), tag)
	}

	session := s.ctx.sql.Session()
	defer session.Close()
	session.Begin()

	affected := int64(0)
	for k := range typ {
		obj := s.ctx.register.Create(typ[k])
		if obj == nil {
			return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, ErrObject.Error(), tag)
		}
		if set, ok := obj.(getsetid); ok {
			set.SetId(ids[k])
		}
		aff, err := session.Delete(obj)
		if err != nil {
			return protocol.ReplyError(protocol.TINY, share.ERR_STORE_SQL, err.Error(), tag)
		}
		affected += aff
	}

	session.Commit()

	return protocol.Reply(protocol.TINY, tag, affected)
}

func (s *Store) Query(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	tag, _ := m.ReadString()
	sql, _ := m.ReadString()
	var args []interface{}
	if err := m.Read(&args); err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, err.Error(), tag)
	}

	sqlorargs := make([]interface{}, 0, 1+len(args))
	sqlorargs = append(sqlorargs, sql)
	if len(args) > 0 {
		sqlorargs = append(sqlorargs, args...)
	}
	result, err := s.ctx.sql.orm.Query(sqlorargs...)
	if err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_STORE_SQL, err.Error(), tag)
	}
	return protocol.Reply(protocol.DEF, tag, result)
}

func (s *Store) Exec(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	tag, _ := m.ReadString()
	sql, _ := m.ReadString()
	var args []interface{}
	if err := m.Read(&args); err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, err.Error(), tag)
	}
	res, err := s.ctx.sql.orm.Exec(sql, args...)
	if err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_STORE_SQL, err.Error(), tag)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_STORE_ERROR, err.Error(), tag)
	}
	return protocol.Reply(protocol.TINY, tag, affected)
}
