package store

import (
	"errors"
	"fmt"
	"ngengine/core/rpc"
	"ngengine/protocol"
	"ngengine/share"
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
	srv := s.ctx.core.LookupOneServiceByType("store")
	if srv == nil {
		s.db = nil
		return
	}

	mb := rpc.GetServiceMailbox(srv.Id)
	s.db = &mb
}

// 解析查询回调的参数
func ParseGetReply(reply *protocol.Message, object interface{}) (errcode int32, err error, tag string) {
	errcode, ar := protocol.ParseReply(reply)
	tag, err = ar.ReadString()
	if err != nil {
		errcode = share.ERR_ARGS_ERROR
		return
	}
	if int(errcode) == share.ERR_TIME_OUT {
		err = rpc.ErrTimeout
		return
	}
	if errcode != 0 {
		errstr, _ := ar.ReadString()
		err = errors.New(errstr)
		return
	}

	if err = ar.Read(object); err != nil {
		errcode = share.ERR_ARGS_ERROR
		return
	}

	return
}

// 查询单条记录，tag为返回的标识符,typ查询的数据类型，condition为条件{column:value}
func (s *StoreClient) Get(tag string, typ string, condition map[string]interface{}, reply rpc.ReplyCB) error {
	if s.db == nil {
		return fmt.Errorf("store not connected")
	}
	if reply == nil {
		return s.ctx.core.Mailto(nil, s.db, "Store.Get", tag, typ, condition)
	}
	return s.ctx.core.MailtoAndCallback(nil, s.db, "Store.Get", reply, tag, typ, condition)
}

func ParseFindReply(reply *protocol.Message, object interface{}) (errcode int32, err error, tag string) {
	errcode, ar := protocol.ParseReply(reply)
	tag, err = ar.ReadString()
	if err != nil {
		errcode = share.ERR_ARGS_ERROR
		return
	}
	if int(errcode) == share.ERR_TIME_OUT {
		err = rpc.ErrTimeout
		return
	}
	if errcode != 0 {
		errstr, _ := ar.ReadString()
		err = errors.New(errstr)
		return
	}

	if err = ar.Read(object); err != nil {
		errcode = share.ERR_ARGS_ERROR
		return
	}

	return
}

// 查找多条记录，tag为返回的标识符,typ查询的数据类型，condition为条件{column:value}，
func (s *StoreClient) Find(tag string, typ string, condition map[string]interface{}, limit int, start int, reply rpc.ReplyCB) error {
	if s.db == nil {
		return fmt.Errorf("store not connected")
	}
	if reply == nil {
		return s.ctx.core.Mailto(nil, s.db, "Store.Find", tag, typ, condition, limit, start)
	}
	return s.ctx.core.MailtoAndCallback(nil, s.db, "Store.Find", reply, tag, typ, condition, limit, start)
}

func ParseInsertReply(reply *protocol.Message) (errcode int32, err error, tag string, affected, id int64) {
	errcode, ar := protocol.ParseReply(reply)
	tag, err = ar.ReadString()
	if err != nil {
		errcode = share.ERR_ARGS_ERROR
		return
	}
	if int(errcode) == share.ERR_TIME_OUT {
		err = rpc.ErrTimeout
		return
	}
	if errcode != 0 {
		errstr, _ := ar.ReadString()
		err = errors.New(errstr)
		return
	}

	affected, err = ar.ReadInt64()
	if err != nil {
		errcode = share.ERR_ARGS_ERROR
		return
	}
	id, err = ar.ReadInt64()
	if err != nil {
		errcode = share.ERR_ARGS_ERROR
	}
	return
}

// 查找一条记录，tag为返回的标识符,typ查询的数据类型，object待插入的数据
func (s *StoreClient) Insert(tag string, typ string, object interface{}, reply rpc.ReplyCB) error {
	if s.db == nil {
		return fmt.Errorf("store not connected")
	}
	if reply == nil {
		return s.ctx.core.Mailto(nil, s.db, "Store.Insert", tag, typ, object)
	}
	return s.ctx.core.MailtoAndCallback(nil, s.db, "Store.Insert", reply, tag, typ, object)
}

func ParseMultiInsertReply(reply *protocol.Message) (errcode int32, err error, tag string) {
	errcode, ar := protocol.ParseReply(reply)
	tag, err = ar.ReadString()
	if err != nil {
		errcode = share.ERR_ARGS_ERROR
		return
	}
	if int(errcode) == share.ERR_TIME_OUT {
		err = rpc.ErrTimeout
		return
	}
	if errcode != 0 {
		errstr, _ := ar.ReadString()
		err = errors.New(errstr)
		return
	}

	return
}

// 批量插入，tag为返回的标识符,typ查询的object数据类型集合,object待插入的数据集合
func (s *StoreClient) MultiInsert(tag string, reply rpc.ReplyCB, typ []string, object []interface{}) error {
	if s.db == nil {
		return fmt.Errorf("store not connected")
	}

	if len(typ) != len(object) {
		return fmt.Errorf("typ and object count not equal")
	}

	var params []interface{}
	params = append(params, tag)
	params = append(params, typ)
	for k := range object {
		params = append(params, object[k])
	}
	if reply == nil {
		return s.ctx.core.Mailto(nil, s.db, "Store.MultiInsert", params...)
	}

	return s.ctx.core.MailtoAndCallback(nil, s.db, "Store.MultiInsert", reply, params...)
}

func ParseUpdateReply(reply *protocol.Message) (errcode int32, err error, tag string, affected int64) {
	errcode, ar := protocol.ParseReply(reply)
	tag, err = ar.ReadString()
	if err != nil {
		errcode = share.ERR_ARGS_ERROR
		return
	}
	if int(errcode) == share.ERR_TIME_OUT {
		err = rpc.ErrTimeout
		return
	}
	if errcode != 0 {
		errstr, _ := ar.ReadString()
		err = errors.New(errstr)
		return
	}
	affected, err = ar.ReadInt64()
	if err != nil {
		errcode = share.ERR_ARGS_ERROR
	}
	return
}

// 更新一条记录，tag为返回的标识符,typ查询的数据类型，cols更新的列，condition为条件{column:value}，object待插入的数据
func (s *StoreClient) Update(tag string, typ string, cols []string, condition map[string]interface{}, object interface{}, reply rpc.ReplyCB) error {
	if s.db == nil {
		return fmt.Errorf("store not connected")
	}
	if reply == nil {
		return s.ctx.core.Mailto(nil, s.db, "Store.Update", tag, typ, cols, condition, object)
	}
	return s.ctx.core.MailtoAndCallback(nil, s.db, "Store.Update", reply, tag, typ, cols, condition, object)
}

// 批量更新，tag为返回的标识符,typ查询的object数据类型集合,object待插入的数据集合
func (s *StoreClient) MultiUpdate(tag string, typ []string, object []interface{}, reply rpc.ReplyCB) error {
	if s.db == nil {
		return fmt.Errorf("store not connected")
	}

	var params []interface{}
	params = append(params, tag)
	params = append(params, typ)
	for k := range object {
		params = append(params, object[k])
	}

	if reply == nil {
		return s.ctx.core.Mailto(nil, s.db, "Store.MultiUpdate", params...)
	}
	return s.ctx.core.MailtoAndCallback(nil, s.db, "Store.MultiUpdate", reply, params...)
}

func ParseDeleteReply(reply *protocol.Message) (errcode int32, err error, tag string, affected int64) {
	errcode, ar := protocol.ParseReply(reply)
	tag, err = ar.ReadString()
	if err != nil {
		errcode = share.ERR_ARGS_ERROR
		return
	}
	if int(errcode) == share.ERR_TIME_OUT {
		err = rpc.ErrTimeout
		return
	}
	if errcode != 0 {
		errstr, _ := ar.ReadString()
		err = errors.New(errstr)
		return
	}
	affected, err = ar.ReadInt64()
	if err != nil {
		errcode = share.ERR_ARGS_ERROR
	}
	return
}

// 删除一条记录，tag为返回的标识符,typ查询的数据类型，待删除对象的id
func (s *StoreClient) Delete(tag string, typ string, id int64, reply rpc.ReplyCB) error {
	if s.db == nil {
		return fmt.Errorf("store not connected")
	}
	if reply == nil {
		return s.ctx.core.Mailto(nil, s.db, "Store.Delete", tag, typ, id)
	}
	return s.ctx.core.MailtoAndCallback(nil, s.db, "Store.Delete", reply, tag, typ, id)
}

// 删除一条记录，tag为返回的标识符,typ查询的数据类型，object待删除的数据(非零值为条件)
func (s *StoreClient) DeleteByObject(tag string, typ string, object interface{}, reply rpc.ReplyCB) error {
	if s.db == nil {
		return fmt.Errorf("store not connected")
	}
	if reply == nil {
		return s.ctx.core.Mailto(nil, s.db, "Store.Delete2", tag, typ, object)
	}
	return s.ctx.core.MailtoAndCallback(nil, s.db, "Store.Delete2", reply, tag, typ, object)
}

// 删除一条记录，tag为返回的标识符,typ查询的数据类型，待删除对象的id
func (s *StoreClient) MultiDelete(tag string, typ []string, id []int64, reply rpc.ReplyCB) error {
	if s.db == nil {
		return fmt.Errorf("store not connected")
	}

	if reply == nil {
		return s.ctx.core.Mailto(nil, s.db, "Store.Delete3", tag, typ, id)
	}
	return s.ctx.core.MailtoAndCallback(nil, s.db, "Store.Delete3", reply, tag, typ, id)
}

func ParseQueryReply(reply *protocol.Message) (errcode int32, err error, tag string, result []map[string][]byte) {
	errcode, ar := protocol.ParseReply(reply)
	tag, err = ar.ReadString()
	if err != nil {
		errcode = share.ERR_ARGS_ERROR
		return
	}
	if int(errcode) == share.ERR_TIME_OUT {
		err = rpc.ErrTimeout
		return
	}
	if errcode != 0 {
		errstr, _ := ar.ReadString()
		err = errors.New(errstr)
		return
	}

	err = ar.Read(&result)
	return
}

// 原生sql查询，tag为返回的标识符，sql为查询语句，args是参数
func (s *StoreClient) Query(tag string, sql string, args []interface{}, reply rpc.ReplyCB) error {
	if s.db == nil {
		return fmt.Errorf("store not connected")
	}
	if reply == nil {
		return s.ctx.core.Mailto(nil, s.db, "Store.Query", tag, sql, args)
	}
	return s.ctx.core.MailtoAndCallback(nil, s.db, "Store.Query", reply, tag, sql, args)
}

func ParseExecReply(reply *protocol.Message) (errcode int32, err error, tag string, affected int64) {
	errcode, ar := protocol.ParseReply(reply)
	tag, err = ar.ReadString()
	if err != nil {
		errcode = share.ERR_ARGS_ERROR
		return
	}
	if int(errcode) == share.ERR_TIME_OUT {
		err = rpc.ErrTimeout
		return
	}
	if errcode != 0 {
		errstr, _ := ar.ReadString()
		err = errors.New(errstr)
		return
	}
	affected, err = ar.ReadInt64()
	if err != nil {
		errcode = share.ERR_ARGS_ERROR
	}
	return
}

// 执行原生sql语句，tag为返回的标识符，sql为执行语句，args是参数
func (s *StoreClient) Exec(tag string, sql string, args []interface{}, reply rpc.ReplyCB) error {
	if s.db == nil {
		return fmt.Errorf("store not connected")
	}
	if reply == nil {
		return s.ctx.core.Mailto(nil, s.db, "Store.Exec", tag, sql, args)
	}
	return s.ctx.core.MailtoAndCallback(nil, s.db, "Store.Exec", reply, tag, sql, args)
}

// 批量插入，tag为返回的标识符,typ查询的object数据类型集合,object待插入的数据集合
func (s *StoreClient) Custom(tag string, reply rpc.ReplyCB, method string, args ...interface{}) error {
	if s.db == nil {
		return fmt.Errorf("store not connected")
	}

	var params []interface{}
	params = append(params, tag)
	params = append(params, args...)
	if reply == nil {
		return s.ctx.core.Mailto(nil, s.db, method, params...)
	}

	return s.ctx.core.MailtoAndCallback(nil, s.db, method, reply, params...)
}
