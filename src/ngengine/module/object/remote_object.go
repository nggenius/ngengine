package object

import (
	"ngengine/core/rpc"
	"ngengine/utils"
)

// 是否是复制对象
func (o *ObjectWitness) Dummy() bool {
	return o.dummy
}

// 设置为复制对象
func (o *ObjectWitness) SetDummy(c bool) {
	o.dummy = c
}

// 同步状态
func (o *ObjectWitness) Sync() bool {
	return o.sync
}

// 设置同步状态
func (o *ObjectWitness) SetSync(s bool) {
	o.sync = s
}

// 原始对象
func (o *ObjectWitness) Original() *rpc.Mailbox {
	return o.original
}

// 设置原始对象
func (o *ObjectWitness) SetOriginal(m *rpc.Mailbox) {
	o.original = m
}

func (o *ObjectWitness) RemoteUpdateAttr(name string, val interface{}) {
	if o.original == nil {
		o.factory.owner.Core.LogErr("original is nil")
		return
	}

	o.factory.owner.Core.Mailto(&o.objid, o.original, "object.UpdateAttr", name, val)
}

func (o *ObjectWitness) RemoteUpdateTuple(name string, val interface{}) {
	if o.original == nil {
		o.factory.owner.Core.LogErr("original is nil")
		return
	}
	o.factory.owner.Core.Mailto(&o.objid, o.original, "object.UpdateTuple", name, val)
}

func (o *ObjectWitness) RemoteAddTableRow(name string, row int) {
	if o.original == nil {
		o.factory.owner.Core.LogErr("original is nil")
		return
	}
	o.factory.owner.Core.Mailto(&o.objid, o.original, "object.AddTableRow", name, row)
}

func (o *ObjectWitness) RemoteAddTableRowValue(name string, row int, val ...interface{}) {
	if o.original == nil {
		o.factory.owner.Core.LogErr("original is nil")
		return
	}
	o.factory.owner.Core.Mailto(&o.objid, o.original, "object.AddTableRowValue", name, row, val)
}

func (o *ObjectWitness) RemoteSetTableRowValue(name string, row int, val ...interface{}) {
	if o.original == nil {
		o.factory.owner.Core.LogErr("original is nil")
		return
	}
	o.factory.owner.Core.Mailto(&o.objid, o.original, "object.SetTableRowValue", name, row, val)
}

func (o *ObjectWitness) RemoteDelTableRow(name string, row int) {
	if o.original == nil {
		o.factory.owner.Core.LogErr("original is nil")
		return
	}
	o.factory.owner.Core.Mailto(&o.objid, o.original, "object.DelTableRow", name, row)
}

func (o *ObjectWitness) RemoteClearTable(name string) {
	if o.original == nil {
		o.factory.owner.Core.LogErr("original is nil")
		return
	}
	o.factory.owner.Core.Mailto(&o.objid, o.original, "object.ClearTable", name)
}

func (o *ObjectWitness) RemoteChangeTable(name string, row, col int, val interface{}) {
	if o.original == nil {
		o.factory.owner.Core.LogErr("original is nil")
		return
	}

	o.factory.owner.Core.Mailto(&o.objid, o.original, "object.ChangeTable", name, row, col, val)
}

// RemoteLockObj 远程上锁
func (o *ObjectWitness) RemoteLockObj(lockID uint32) {
	if o.original == nil {
		o.factory.owner.Core.LogErr("original is nil")
		return
	}
	o.factory.owner.Core.Mailto(&o.objid, o.original, "object.LockObj", o.original, lockID)
}

// RemoteUnLockObj 远程解锁
func (o *ObjectWitness) RemoteUnLockObj(lockID uint32) error {
	if o.original == nil {
		o.factory.owner.Core.LogErr("original is nil")
		return nil
	}
	return o.factory.owner.Core.Mailto(&o.objid, o.original, "object.UnLockObj", o.original, lockID)
}

// RemoteLockObjSuccess 远程上锁成功通知
func (o *ObjectWitness) RemoteLockObjSuccess(lockID uint32) error {
	if o.locker == nil {
		o.factory.owner.Core.LogErr("locker is nil")
		return nil
	}
	return o.factory.owner.Core.MailtoAndCallback(&o.objid, &o.locker.Locker, "object.LockObjSuccess",
		func(p interface{}, e *rpc.Error, l *utils.LoadArchive) {
			if e != nil {
				// 如果远端已经没有这个对象了，解开锁
				o.UnLockObjSuccess(false)
			}
		},
		nil, o.locker.Locker, lockID)
}

// RemoteUnLockObjSuccess 远程解锁成功通知
func (o *ObjectWitness) RemoteUnLockObjSuccess() error {
	if o.locker == nil {
		o.factory.owner.Core.LogErr("locker is nil")
		return nil
	}
	return o.factory.owner.Core.Mailto(&o.objid, &o.locker.Locker, "object.UnLockObjSuccess", o.locker.Locker)
}
