package object

import "ngengine/core/rpc"

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
		o.factory.owner.core.LogErr("original is nil")
		return
	}

	o.factory.owner.core.Mailto(&o.objid, o.original, "object.UpdateAttr", name, val)
}

func (o *ObjectWitness) RemoteUpdateTuple(name string, val interface{}) {
	if o.original == nil {
		o.factory.owner.core.LogErr("original is nil")
		return
	}
	o.factory.owner.core.Mailto(&o.objid, o.original, "object.UpdateTuple", name, val)
}

func (o *ObjectWitness) RemoteAddTableRow(name string, row int) {
	if o.original == nil {
		o.factory.owner.core.LogErr("original is nil")
		return
	}
	o.factory.owner.core.Mailto(&o.objid, o.original, "object.AddTableRow", name, row)
}

func (o *ObjectWitness) RemoteAddTableRowValue(name string, row int, val ...interface{}) {
	if o.original == nil {
		o.factory.owner.core.LogErr("original is nil")
		return
	}
	o.factory.owner.core.Mailto(&o.objid, o.original, "object.AddTableRowValue", name, row, val)
}

func (o *ObjectWitness) RemoteSetTableRowValue(name string, row int, val ...interface{}) {
	if o.original == nil {
		o.factory.owner.core.LogErr("original is nil")
		return
	}
	o.factory.owner.core.Mailto(&o.objid, o.original, "object.SetTableRowValue", name, row, val)
}

func (o *ObjectWitness) RemoteDelTableRow(name string, row int) {
	if o.original == nil {
		o.factory.owner.core.LogErr("original is nil")
		return
	}
	o.factory.owner.core.Mailto(&o.objid, o.original, "object.DelTableRow", name, row)
}

func (o *ObjectWitness) RemoteClearTable(name string) {
	if o.original == nil {
		o.factory.owner.core.LogErr("original is nil")
		return
	}
	o.factory.owner.core.Mailto(&o.objid, o.original, "object.ClearTable", name)
}

func (o *ObjectWitness) RemoteChangeTable(name string, row, col int, val interface{}) {
	if o.original == nil {
		o.factory.owner.core.LogErr("original is nil")
		return
	}

	o.factory.owner.core.Mailto(&o.objid, o.original, "object.ChangeTable", name, row, col, val)
}
