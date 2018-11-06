package object

import (
	"ngengine/core/rpc"
	"ngengine/protocol"
	"ngengine/share"
)

type SyncObject struct {
	owner *ObjectModule
}

func (s *SyncObject) RegisterCallback(svr rpc.Servicer) {
	svr.RegisterCallback("Rebase", s.Rebase)
	svr.RegisterCallback("Flip", s.Flip)
	svr.RegisterCallback("Replicate", s.Replicate)
	svr.RegisterCallback("UpdateAttr", s.UpdateAttr)
	svr.RegisterCallback("UpdateTuple", s.UpdateTuple)
	svr.RegisterCallback("AddTableRow", s.AddTableRow)
	svr.RegisterCallback("AddTableRowValue", s.AddTableRowValue)
	svr.RegisterCallback("SetTableRowValue", s.SetTableRowValue)
	svr.RegisterCallback("DelTableRow", s.DelTableRow)
	svr.RegisterCallback("ClearTable", s.ClearTable)
	svr.RegisterCallback("ChangeTable", s.ChangeTable)
	svr.RegisterCallback("LockObj", s.LockObj)
	svr.RegisterCallback("UnLockObj", s.UnLockObj)
	svr.RegisterCallback("LockObjSuccess", s.LockObjSuccess)
	svr.RegisterCallback("UnLockObjSuccess", s.UnLockObjSuccess)
}

// Rebase 变基
func (s *SyncObject) Rebase(src rpc.Mailbox, dest rpc.Mailbox, msg *protocol.Message) (int32, *protocol.Message) {
	return 0, nil
}

// Flip 交换控制权(由origin处理)
func (s *SyncObject) Flip(src rpc.Mailbox, dest rpc.Mailbox, msg *protocol.Message) (int32, *protocol.Message) {
	obj, err := s.owner.FindObject(dest)
	if obj == nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_ORIGIN_NOT_FOUND, err.Error())
	}

	return 0, nil
}

// 对象复制
func (s *SyncObject) Replicate(src rpc.Mailbox, dest rpc.Mailbox, msg *protocol.Message) (int32, *protocol.Message) {

	var tag int
	var data []byte
	err := protocol.ParseArgs(msg, &tag, &data)
	if err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, err.Error())
	}

	f := s.owner.Factory(share.OBJECT_TYPE_GHOST)
	obj, err := f.Decode(data)
	if err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_OBJECT_REPLICATE, err.Error())
	}

	s.owner.fireGlobalEvent(GLOBAL_ADD_DUMY, obj.(Object).ObjId(), rpc.NullMailbox, tag)

	return protocol.Reply(protocol.TINY, obj.(Object).ObjId())
}

// 对象属性变动
func (s *SyncObject) UpdateAttr(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	return 0, nil
}

// 对象tupele属性变动
func (s *SyncObject) UpdateTuple(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	return 0, nil
}

// 对象表格增加一行
func (s *SyncObject) AddTableRow(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	return 0, nil
}

// 对象表格增加一行，并设置值
func (s *SyncObject) AddTableRowValue(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	return 0, nil
}

// 设置表格一行的值
func (s *SyncObject) SetTableRowValue(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	return 0, nil
}

// 对象表格删除一行
func (s *SyncObject) DelTableRow(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	return 0, nil
}

// 对象表格清除所有行
func (s *SyncObject) ClearTable(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	return 0, nil
}

// 对象表格单元格更新
func (s *SyncObject) ChangeTable(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {

	return 0, nil
}

// LockObj 给对象上锁
func (s *SyncObject) LockObj(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	var objmailbox rpc.Mailbox
	m.Read(&objmailbox)

	r, _ := m.ReadUInt32()
	o, err := s.owner.FindObject(objmailbox)
	if err != nil {
		s.owner.Core.LogErr(err)
	}

	if obj, ok := o.(Object); ok {
		obj.AddLocker(sender, r, true)
	}
	return 0, nil
}

// UnLockObj 给对象解锁
func (s *SyncObject) UnLockObj(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	var objmailbox rpc.Mailbox
	m.Read(&objmailbox)

	r, _ := m.ReadUInt32()
	o, err := s.owner.FindObject(objmailbox)
	if err != nil {
		s.owner.Core.LogErr(err)
	}
	if obj, ok := o.(Object); ok {
		obj.UnLockObj(sender, r, true)
	}
	return 0, nil
}

// 远程通知对象上锁成功
func (s *SyncObject) LockObjSuccess(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	var objmailbox rpc.Mailbox
	m.Read(&objmailbox)

	r, _ := m.ReadUInt32()
	o, err := s.owner.FindObject(objmailbox)
	if err != nil {
		s.owner.Core.LogErr(err)
	}

	if obj, ok := o.(Object); ok {
		// 这里远程回复上锁成功所以上锁的就是本对象切这里算上本地的锁
		obj.LockObjSuccess(objmailbox, r, false)
	}
	return 0, nil
}

// 远程通知对象解锁成功
func (s *SyncObject) UnLockObjSuccess(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	var objmailbox rpc.Mailbox
	m.Read(&objmailbox)

	o, err := s.owner.FindObject(objmailbox)
	if err != nil {
		s.owner.Core.LogErr(err)
	}
	if obj, ok := o.(Object); ok {
		obj.UnLockObjSuccess(false)
	}
	return 0, nil
}
