package store

import (
	"fmt"
	"ngengine/core/rpc"
	"ngengine/protocol"
	"ngengine/share"
	"ngengine/utils"
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
	srv := s.ctx.Core.LookupOneServiceByType("store")
	if srv == nil {
		s.db = nil
		return
	}

	mb := rpc.GetServiceMailbox(srv.Id)
	s.db = &mb
}

// ParseGetReply 解析查询回调的参数
func ParseGetReply(err *rpc.Error, ar *utils.LoadArchive, object interface{}) *rpc.Error {
	if err != nil && protocol.CheckRpcError(err) {
		return err
	}
	if e := ar.Read(object); e != nil {
		return rpc.NewError(share.ERR_ARGS_ERROR, e.Error())
	}

	return nil
}

// Get 查询单条记录，identity为返回的标识符,typ查询的数据类型，condition为条件{column:value}
func (s *StoreClient) Get(identity interface{}, typ string, condition map[string]interface{}, reply rpc.ReplyCB) error {
	if s.db == nil {
		return fmt.Errorf("store not connected")
	}
	if reply == nil {
		return s.ctx.Core.Mailto(nil, s.db, "Store.Get", typ, condition)
	}
	return s.ctx.Core.MailtoAndCallback(nil, s.db, "Store.Get", reply, identity, typ, condition)
}

func ParseFindReply(err *rpc.Error, ar *utils.LoadArchive, object interface{}) *rpc.Error {
	if err != nil && protocol.CheckRpcError(err) {
		return err
	}

	if e := ar.Read(object); e != nil {
		return rpc.NewError(share.ERR_ARGS_ERROR, e.Error())
	}

	return nil
}

// 查找多条记录，identity为返回的标识符,typ查询的数据类型，condition为条件{column:value}，
func (s *StoreClient) Find(identity interface{}, typ string, condition map[string]interface{}, limit int, start int, reply rpc.ReplyCB) error {
	if s.db == nil {
		return fmt.Errorf("store not connected")
	}
	if reply == nil {
		return s.ctx.Core.Mailto(nil, s.db, "Store.Find", typ, condition, limit, start)
	}
	return s.ctx.Core.MailtoAndCallback(nil, s.db, "Store.Find", reply, identity, typ, condition, limit, start)
}

// 返回值：err, affected, id
func ParseInsertReply(err *rpc.Error, ar *utils.LoadArchive) (*rpc.Error, int64, int64) {
	if err != nil && protocol.CheckRpcError(err) {
		return err, 0, 0
	}

	affected, e := ar.ReadInt64()
	if e != nil {
		return rpc.NewError(share.ERR_ARGS_ERROR, e.Error()), 0, 0
	}
	id, e := ar.ReadInt64()
	if e != nil {
		return rpc.NewError(share.ERR_ARGS_ERROR, e.Error()), 0, 0
	}
	return nil, affected, id
}

// 查找一条记录，tag为返回的标识符,typ查询的数据类型，object待插入的数据
func (s *StoreClient) Insert(identity interface{}, typ string, object interface{}, reply rpc.ReplyCB) error {
	if s.db == nil {
		return fmt.Errorf("store not connected")
	}
	if reply == nil {
		return s.ctx.Core.Mailto(nil, s.db, "Store.Insert", typ, object)
	}
	return s.ctx.Core.MailtoAndCallback(nil, s.db, "Store.Insert", reply, identity, typ, object)
}

// 返回值：err
func ParseMultiInsertReply(err *rpc.Error, ar *utils.LoadArchive) *rpc.Error {
	if err != nil && protocol.CheckRpcError(err) {
		return err
	}

	return nil
}

// 批量插入，identity为返回的标识符,typ查询的object数据类型集合,object待插入的数据集合
func (s *StoreClient) MultiInsert(identity interface{}, reply rpc.ReplyCB, typ []string, object []interface{}) error {
	if s.db == nil {
		return fmt.Errorf("store not connected")
	}

	if len(typ) != len(object) {
		return fmt.Errorf("typ and object count not equal")
	}

	var params []interface{}
	params = append(params, typ)
	for k := range object {
		params = append(params, object[k])
	}
	if reply == nil {
		return s.ctx.Core.Mailto(nil, s.db, "Store.MultiInsert", params...)
	}

	return s.ctx.Core.MailtoAndCallback(nil, s.db, "Store.MultiInsert", reply, identity, params...)
}

// 返回值:err *rpc.Error, affected int6
func ParseUpdateReply(err *rpc.Error, ar *utils.LoadArchive) (*rpc.Error, int64) {
	if err != nil && protocol.CheckRpcError(err) {
		return err, 0
	}
	affected, e := ar.ReadInt64()
	if e != nil {
		return rpc.NewError(share.ERR_ARGS_ERROR, e.Error()), 0
	}
	return nil, affected
}

// 更新一条记录，identity为返回的标识符,typ查询的数据类型，cols更新的列，condition为条件{column:value}，object待插入的数据
func (s *StoreClient) Update(identity interface{}, typ string, cols []string, condition map[string]interface{}, object interface{}, reply rpc.ReplyCB) error {
	if s.db == nil {
		return fmt.Errorf("store not connected")
	}
	if reply == nil {
		return s.ctx.Core.Mailto(nil, s.db, "Store.Update", typ, cols, condition, object)
	}
	return s.ctx.Core.MailtoAndCallback(nil, s.db, "Store.Update", reply, identity, typ, cols, condition, object)
}

// 批量更新，identity为返回的标识符,typ查询的object数据类型集合,object待插入的数据集合
func (s *StoreClient) MultiUpdate(identity interface{}, typ []string, object []interface{}, reply rpc.ReplyCB) error {
	if s.db == nil {
		return fmt.Errorf("store not connected")
	}

	var params []interface{}
	params = append(params, typ)
	for k := range object {
		params = append(params, object[k])
	}

	if reply == nil {
		return s.ctx.Core.Mailto(nil, s.db, "Store.MultiUpdate", params...)
	}
	return s.ctx.Core.MailtoAndCallback(nil, s.db, "Store.MultiUpdate", reply, identity, params...)
}

//返回值:err *rpc.Error,affected int64
func ParseDeleteReply(err *rpc.Error, ar *utils.LoadArchive) (*rpc.Error, int64) {
	if err != nil && protocol.CheckRpcError(err) {
		return err, 0
	}

	affected, e := ar.ReadInt64()
	if e != nil {
		return rpc.NewError(share.ERR_ARGS_ERROR, e.Error()), 0
	}
	return nil, affected
}

// 删除一条记录，identity为返回的标识符,typ查询的数据类型，待删除对象的id
func (s *StoreClient) Delete(identity interface{}, typ string, id int64, reply rpc.ReplyCB) error {
	if s.db == nil {
		return fmt.Errorf("store not connected")
	}
	if reply == nil {
		return s.ctx.Core.Mailto(nil, s.db, "Store.Delete", typ, id)
	}
	return s.ctx.Core.MailtoAndCallback(nil, s.db, "Store.Delete", reply, identity, typ, id)
}

// 删除一条记录，identity为返回的标识符,typ查询的数据类型，object待删除的数据(非零值为条件)
func (s *StoreClient) DeleteByObject(identity interface{}, typ string, object interface{}, reply rpc.ReplyCB) error {
	if s.db == nil {
		return fmt.Errorf("store not connected")
	}
	if reply == nil {
		return s.ctx.Core.Mailto(nil, s.db, "Store.Delete2", typ, object)
	}
	return s.ctx.Core.MailtoAndCallback(nil, s.db, "Store.Delete2", reply, identity, typ, object)
}

// 删除一条记录，identity为返回的标识符,typ查询的数据类型，待删除对象的id
func (s *StoreClient) MultiDelete(identity interface{}, typ []string, id []int64, reply rpc.ReplyCB) error {
	if s.db == nil {
		return fmt.Errorf("store not connected")
	}

	if reply == nil {
		return s.ctx.Core.Mailto(nil, s.db, "Store.Delete3", typ, id)
	}
	return s.ctx.Core.MailtoAndCallback(nil, s.db, "Store.Delete3", reply, identity, typ, id)
}

// 返回值：err *rpc.Error, result []map[string][]byte
func ParseQueryReply(err *rpc.Error, ar *utils.LoadArchive) (*rpc.Error, []map[string][]byte) {
	if err != nil && protocol.CheckRpcError(err) {
		return err, nil
	}

	var result []map[string][]byte
	if e := ar.Read(&result); e != nil {
		return rpc.NewError(share.ERR_ARGS_ERROR, e.Error()), nil
	}
	return nil, result
}

// 原生sql查询，identity为返回的标识符，sql为查询语句，args是参数
func (s *StoreClient) Query(identity interface{}, sql string, args []interface{}, reply rpc.ReplyCB) error {
	if s.db == nil {
		return fmt.Errorf("store not connected")
	}
	if reply == nil {
		return s.ctx.Core.Mailto(nil, s.db, "Store.Query", sql, args)
	}
	return s.ctx.Core.MailtoAndCallback(nil, s.db, "Store.Query", reply, identity, sql, args)
}

//返回值：err *rpc.Error, affected int64
func ParseExecReply(err *rpc.Error, ar *utils.LoadArchive) (*rpc.Error, int64) {
	if err != nil && protocol.CheckRpcError(err) {
		return err, 0
	}
	affected, e := ar.ReadInt64()
	if e != nil {
		return rpc.NewError(share.ERR_ARGS_ERROR, e.Error()), 0
	}
	return nil, affected
}

// 执行原生sql语句，identity为返回的标识符，sql为执行语句，args是参数
func (s *StoreClient) Exec(identity interface{}, sql string, args []interface{}, reply rpc.ReplyCB) error {
	if s.db == nil {
		return fmt.Errorf("store not connected")
	}
	if reply == nil {
		return s.ctx.Core.Mailto(nil, s.db, "Store.Exec", sql, args)
	}
	return s.ctx.Core.MailtoAndCallback(nil, s.db, "Store.Exec", reply, identity, sql, args)
}

// 批量插入，identity为返回的标识符,typ查询的object数据类型集合,object待插入的数据集合
func (s *StoreClient) Custom(identity interface{}, reply rpc.ReplyCB, method string, args ...interface{}) error {
	if s.db == nil {
		return fmt.Errorf("store not connected")
	}

	if reply == nil {
		return s.ctx.Core.Mailto(nil, s.db, method, args...)
	}

	return s.ctx.Core.MailtoAndCallback(nil, s.db, method, reply, identity, args...)
}
