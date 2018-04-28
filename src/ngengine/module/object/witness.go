package object

// Witness的作用是监听数据对象的变动，作为一个目击者。由数据对象直接持有。
// 目击者作为一个事件集散地，对所有的数据变动的事件进行调度。由第三方注册。
// 目击者不关注变动的细节，只进行转发，由第三方进行细节的处理。
import (
	"container/list"
	"fmt"
	"ngengine/core/rpc"
)

const (
	TABLE_INIT = iota + 1
	TABLE_ADD_ROW
	TABLE_REMOVE_ROW
	TABLE_CLEAR_ROW
	TABLE_GRID_CHANGE
	TABLE_SET_ROW
)

type attrObserver interface {
	Init(object Object)
	UpdateAttr(name string, val interface{}, old interface{})
	UpdateTuple(name string, val interface{}, old interface{})
}

type tableObserver interface {
	Init(object Object)
	UpdateTable(name string, op_type, row, col int)
}

// LockCallBack 回调
type LockCallBack func()

type ObjectWitness struct {
	object        Object
	objid         rpc.Mailbox
	factory       *Factory
	original      *rpc.Mailbox
	dummy         bool // 是否是副本
	sync          bool // 同步状态
	silence       bool // 沉默状态
	attrobserver  map[string]attrObserver
	tableobserver map[string]tableObserver

	Islock      bool                    // 是否已加锁
	LockCount   uint32                  // 加锁计数
	LockCb      map[uint32]LockCallBack // 回调函数
	LockerQueue *list.List              // 加锁的队列
	locker      *Locker                 // 当前上锁的人以及信息
}

// Factory 获取工厂
func (o *ObjectWitness) Factory() *Factory {
	return o.factory
}

// SetFactory 所属的工厂
func (o *ObjectWitness) SetFactory(f *Factory) {
	o.factory = f
}

// ObjId 唯一ID
func (o *ObjectWitness) ObjId() rpc.Mailbox {
	return o.objid
}

// SetObjId 设置唯一ID
func (o *ObjectWitness) SetObjId(id rpc.Mailbox) {
	o.objid = id
}

// Silence 沉默状态
func (o *ObjectWitness) Silence() bool {
	return o.silence
}

// SetSilence 设置沉默状态
func (o *ObjectWitness) SetSilence(s bool) {
	o.silence = s
}

// Witness 设置对象
func (o *ObjectWitness) Witness(obj Object) {
	o.object = obj
	o.attrobserver = make(map[string]attrObserver)
	o.tableobserver = make(map[string]tableObserver)
	o.LockerQueue = list.New()
}

// AddAttrObserver 增加属性观察者,这里的name是观察者的标识符，不是属性名称
func (o *ObjectWitness) AddAttrObserver(name string, observer attrObserver) error {
	if _, dup := o.attrobserver[name]; dup {
		return fmt.Errorf("add attr observer twice %s", name)
	}

	o.attrobserver[name] = observer
	observer.Init(o.object)
	return nil
}

// RemoveAttrObserver 删除属性观察者
func (o *ObjectWitness) RemoveAttrObserver(name string) {
	delete(o.attrobserver, name)
}

// AddTableObserver 增加表格观察者,这里的name是观察者的标识符，不是表格名称
func (o *ObjectWitness) AddTableObserver(name string, observer tableObserver) error {
	if _, dup := o.tableobserver[name]; dup {
		return fmt.Errorf("add table observer twice %s", name)
	}

	o.tableobserver[name] = observer
	observer.Init(o.object)
	return nil
}

// RemoveTableObserver 删除表格观察者
func (o *ObjectWitness) RemoveTableObserver(name string) {
	delete(o.tableobserver, name)
}

// UpdateAttr 对象属性变动(由object调用)
func (o *ObjectWitness) UpdateAttr(name string, val interface{}, old interface{}) {
	if o.dummy && !o.sync { // 需要操作远程对象
		o.RemoteUpdateAttr(name, val)
		return
	}
	if o.silence {
		return
	}
	for _, observer := range o.attrobserver {
		observer.UpdateAttr(name, val, old)
	}
}

// UpdateTuple 对象tupele属性变动(由object调用)
func (o *ObjectWitness) UpdateTuple(name string, val interface{}, old interface{}) {
	if o.dummy && !o.sync { // 需要操作远程对象
		o.RemoteUpdateTuple(name, val)
		return
	}
	if o.silence {
		return
	}
	for _, observer := range o.attrobserver {
		observer.UpdateTuple(name, val, old)
	}
}

// AddTableRow 对象表格增加一行(由object调用)
func (o *ObjectWitness) AddTableRow(name string, row int) {
	if o.dummy && !o.sync { // 需要操作远程对象
		o.RemoteAddTableRow(name, row)
		return
	}
	if o.silence {
		return
	}
	for _, observer := range o.tableobserver {
		observer.UpdateTable(name, TABLE_ADD_ROW, row, 0)
	}
}

// AddTableRowValue 对象表格增加一行，并设置值(由object调用)
func (o *ObjectWitness) AddTableRowValue(name string, row int, val ...interface{}) {
	if o.dummy && !o.sync { // 需要操作远程对象
		o.RemoteAddTableRowValue(name, row, val...)
		return
	}
	if o.silence {
		return
	}
	for _, observer := range o.tableobserver {
		observer.UpdateTable(name, TABLE_ADD_ROW, row, 0)
	}
}

// SetTableRowValue 设置表格一行的值(由object调用)
func (o *ObjectWitness) SetTableRowValue(name string, row int, val ...interface{}) {
	if o.dummy && !o.sync { // 需要操作远程对象
		o.RemoteSetTableRowValue(name, row, val...)
		return
	}
	if o.silence {
		return
	}
	for _, observer := range o.tableobserver {
		observer.UpdateTable(name, TABLE_SET_ROW, row, 0)
	}
}

// DelTableRow 对象表格删除一行(由object调用)
func (o *ObjectWitness) DelTableRow(name string, row int) {
	if o.dummy && !o.sync { // 需要操作远程对象
		o.RemoteDelTableRow(name, row)
		return
	}
	if o.silence {
		return
	}
	for _, observer := range o.tableobserver {
		observer.UpdateTable(name, TABLE_REMOVE_ROW, row, 0)
	}
}

// ClearTable 对象表格清除所有行(由object调用)
func (o *ObjectWitness) ClearTable(name string) {
	if o.dummy && !o.sync { // 需要操作远程对象
		o.RemoteClearTable(name)
		return
	}
	if o.silence {
		return
	}
	for _, observer := range o.tableobserver {
		observer.UpdateTable(name, TABLE_CLEAR_ROW, 0, 0)
	}
}

// ChangeTable 对象表格单元格更新(由object调用)
func (o *ObjectWitness) ChangeTable(name string, row, col int, val interface{}) {
	if o.dummy && !o.sync { // 需要操作远程对象
		o.RemoteChangeTable(name, row, col, val)
		return
	}
	if o.silence {
		return
	}
	for _, observer := range o.tableobserver {
		observer.UpdateTable(name, TABLE_GRID_CHANGE, row, col)
	}
}
