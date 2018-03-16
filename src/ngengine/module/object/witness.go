package object

import (
	"fmt"
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

type ObjectWitness struct {
	object        Object
	factory       *Factory
	attrobserver  map[string]attrObserver
	tableobserver map[string]tableObserver
}

// 获取工厂
func (o *ObjectWitness) Factory() *Factory {
	return o.factory
}

// 所属的工厂
func (o *ObjectWitness) SetFactory(f *Factory) {
	o.factory = f
}

// 设置对象
func (o *ObjectWitness) Witness(obj Object) {
	o.object = obj
	o.attrobserver = make(map[string]attrObserver)
	o.tableobserver = make(map[string]tableObserver)
}

// 增加属性观察者,这里的name是观察者的标识符，不是属性名称
func (o *ObjectWitness) AddAttrObserver(name string, observer attrObserver) error {
	if _, dup := o.attrobserver[name]; dup {
		return fmt.Errorf("add attr observer twice %s", name)
	}

	o.attrobserver[name] = observer
	observer.Init(o.object)
	return nil
}

// 增加表格观察者,这里的name是观察者的标识符，不是表格名称
func (o *ObjectWitness) AddTableObserver(name string, observer tableObserver) error {
	if _, dup := o.tableobserver[name]; dup {
		return fmt.Errorf("add table observer twice %s", name)
	}

	o.tableobserver[name] = observer
	observer.Init(o.object)
	return nil
}

// 对象属性变动(由object调用)
func (o *ObjectWitness) UpdateAttr(name string, val interface{}, old interface{}) {
	for _, observer := range o.attrobserver {
		observer.UpdateAttr(name, val, old)
	}
}

// 对象tupele属性变动(由object调用)
func (o *ObjectWitness) UpdateTuple(name string, val interface{}, old interface{}) {
	for _, observer := range o.attrobserver {
		observer.UpdateTuple(name, val, old)
	}
}

// 对象表格增加一行(由object调用)
func (o *ObjectWitness) AddTableRow(name string, row int) {
	for _, observer := range o.tableobserver {
		observer.UpdateTable(name, TABLE_ADD_ROW, row, 0)
	}
}

// 对象表格增加一行，并设置值(由object调用)
func (o *ObjectWitness) AddTableRowValue(name string, row int, val ...interface{}) {
	for _, observer := range o.tableobserver {
		observer.UpdateTable(name, TABLE_ADD_ROW, row, 0)
		observer.UpdateTable(name, TABLE_SET_ROW, row, 0)
	}
}

// 对象表格删除一行(由object调用)
func (o *ObjectWitness) DelTableRow(name string, row int) {
	for _, observer := range o.tableobserver {
		observer.UpdateTable(name, TABLE_REMOVE_ROW, row, 0)
	}
}

// 对象表格清除所有行(由object调用)
func (o *ObjectWitness) ClearTable(name string) {
	for _, observer := range o.tableobserver {
		observer.UpdateTable(name, TABLE_CLEAR_ROW, 0, 0)
	}
}

// 对象表格单元格更新(由object调用)
func (o *ObjectWitness) ChangeTable(name string, row, col int, val interface{}) {
	for _, observer := range o.tableobserver {
		observer.UpdateTable(name, TABLE_GRID_CHANGE, row, col)
	}
}
