package object

import (
	"ngengine/core/rpc"
	"ngengine/protocol"
)

type SyncObject struct {
	owner *ObjectModule
}

func (s *SyncObject) RegisterCallback(svr rpc.Servicer) {
	svr.RegisterCallback("UpdateAttr", s.UpdateAttr)
	svr.RegisterCallback("UpdateTuple", s.UpdateTuple)
	svr.RegisterCallback("AddTableRow", s.AddTableRow)
	svr.RegisterCallback("AddTableRowValue", s.AddTableRowValue)
	svr.RegisterCallback("SetTableRowValue", s.SetTableRowValue)
	svr.RegisterCallback("DelTableRow", s.DelTableRow)
	svr.RegisterCallback("ClearTable", s.ClearTable)
	svr.RegisterCallback("ChangeTable", s.ChangeTable)
}

// 对象属性变动
func (s *SyncObject) UpdateAttr(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	return 0, nil
}

// 对象tupele属性变动
func (s *SyncObject) UpdateTuple(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	return 0, nil
}

// 对象表格增加一行
func (s *SyncObject) AddTableRow(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	return 0, nil
}

// 对象表格增加一行，并设置值
func (s *SyncObject) AddTableRowValue(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	return 0, nil
}

// 设置表格一行的值
func (s *SyncObject) SetTableRowValue(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	return 0, nil
}

// 对象表格删除一行
func (s *SyncObject) DelTableRow(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	return 0, nil
}

// 对象表格清除所有行
func (s *SyncObject) ClearTable(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	return 0, nil
}

// 对象表格单元格更新
func (s *SyncObject) ChangeTable(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	return 0, nil
}
