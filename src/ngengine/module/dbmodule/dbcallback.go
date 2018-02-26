package dbmodule

import (
	"fmt"
	"ngengine/core/rpc"
	"ngengine/module/dbmodule/dbdriver"
	"ngengine/protocol"
	"ngengine/share"
)

// DbCallBack rpc调用的函数
type DbCallBack struct {
	DbModule
}

// RegisterCallback 注册rpc回调
func (a *DbCallBack) RegisterCallback(s rpc.Servicer) {
	s.RegisterCallback("DBGet", a.Get)
	s.RegisterCallback("DBSqlQuery", a.SQLQuery)
	s.RegisterCallback("DBSqlExec", a.SQLExec)
	s.RegisterCallback("DBInsert", a.Insert)
	s.RegisterCallback("DBUpdate", a.Update)
	s.RegisterCallback("DBDelete", a.Delete)
	s.RegisterCallback("DBFind", a.Find)
}

// Get a.owner.CoreApi.MailtoAndCallback(nil, &dest, "DBModule.DBGet", a.DbCallBack, dbtable.NX_BASE, share.MessageBox{Message: &nx_base.NxChangename{Uid: "333"}})
func (a *DbCallBack) Get(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	DbName, _ := m.ReadString()

	mes := &share.MessageBox{}
	er := m.ReadObject(mes)
	if er != nil {
		fmt.Println(er)
	}

	g, err := dbdriver.Get(DbName, mes.Message)

	if err != nil {
		return 0, protocol.ReplyMessage(protocol.TINY, 1)
	}

	a.Core.LogDebug("查询出来的数据:", g)
	// 回复查询的结构
	return 0, protocol.ReplyMessage(protocol.DEF, share.MessageBox{Message: g})
}

// SQLQuery 直接传入查询类sql语
func (a *DbCallBack) SQLQuery(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	DbName, _ := m.ReadString()
	sql, _ := m.ReadString()

	g, err := dbdriver.SQLQuery(DbName, sql)

	if err != nil {
		return 0, protocol.ReplyMessage(protocol.TINY, 1)
	}

	a.Core.LogDebug("查询出来的数据:", g)
	// 回复查询的结构
	return 0, protocol.ReplyMessage(protocol.DEF, share.MessageBox{Message: g})
}

// SQLExec 直接传入执行类sql语言
func (a *DbCallBack) SQLExec(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	DbName, _ := m.ReadString()
	sql, _ := m.ReadString()

	g, err := dbdriver.SQLExec(DbName, sql)

	if err != nil {
		return 0, protocol.ReplyMessage(protocol.TINY, 1)
	}

	a.Core.LogDebug("执行结果:", g)
	// 回复查询的结构
	return 0, protocol.ReplyMessage(protocol.DEF, g)
}

//Insert  a.owner.CoreApi.Mailto(nil, &dest, "DBModule.DBInsert", dbtable.NX_BASE, share.MessageBox{Message: &nx_base.NxChangename{
// 	Name:    "东山日照",
// 	Uid:     "1234567",
// 	NewName: "飞舞乱倒",
// }})
func (a *DbCallBack) Insert(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	DbName, _ := m.ReadString()

	mes := &share.MessageBox{}
	er := m.ReadObject(mes)
	if er != nil {
		fmt.Println(er)
	}

	err := dbdriver.Insert(DbName, mes.Message)

	if err != nil {
		a.Core.LogDebug("错误:", err)
		return 0, protocol.ReplyMessage(protocol.TINY, 1)
	}

	// 回复查询的结构
	return 0, nil
}

// Update a.owner.CoreApi.Mailto(nil, &dest, "DBModule.DBUpdate", dbtable.NX_BASE, share.MessageBox{Message: &nx_base.NxChangename{
// 	Name:    "尚东龙兴",
// 	Uid:     "1234567",
// 	NewName: "茶田非君",
// }}, share.MessageBox{Message: &nx_base.NxChangename{
// 	Name: "东山日照",
// }})
func (a *DbCallBack) Update(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	DbName, _ := m.ReadString()

	mes := &share.MessageBox{}
	er := m.ReadObject(mes)
	if er != nil {
		fmt.Println(er)
	}

	old := &share.MessageBox{}
	err1 := m.ReadObject(old)
	if err1 != nil {
		fmt.Println(err1)
	}
	err := dbdriver.Update(DbName, mes.Message, old.Message)

	if err != nil {
		a.Core.LogDebug("错误:", err)
		return 0, protocol.ReplyMessage(protocol.TINY, 1)
	}
	// 回复查询的结构
	return 0, nil
}

// Delete a.owner.CoreApi.Mailto(nil, &dest, "DBModule.DBDelete", dbtable.NX_BASE, share.MessageBox{Message: &nx_base.NxChangename{
// 	Name: "尚东龙兴",
// }})
func (a *DbCallBack) Delete(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	DbName, _ := m.ReadString()

	mes := &share.MessageBox{}
	er := m.ReadObject(mes)
	if er != nil {
		fmt.Println(er)
	}

	err := dbdriver.Delete(DbName, mes.Message)
	if err != nil {
		a.Core.LogDebug("错误:", err)
		return 0, protocol.ReplyMessage(protocol.TINY, 1)
	}
	return 0, protocol.ReplyMessage(protocol.TINY, 1)
}

// Find a.owner.CoreApi.MailtoAndCallback(nil, &dest, "DBModule.DBFind", a.DbCallBack, dbtable.NX_BASE, share.MessageBox{Message: &nx_base.NxChangename{
// 	Uid: "333",
// }})
func (a *DbCallBack) Find(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	DbName, _ := m.ReadString()

	mes := &share.MessageBox{}
	er := m.ReadObject(mes)
	if er != nil {
		fmt.Println(er)
	}
	g, err := dbdriver.Find(DbName, mes.Message)
	if err != nil {
		fmt.Println(err)
		return 0, protocol.ReplyMessage(protocol.TINY, 1)
	}
	a.Core.LogDebug("执行结果:", g)
	// 这里是数组没有注册gob然后返回值失败
	return 0, protocol.ReplyMessage(protocol.DEF, share.MessageBox{Message: g})
}
