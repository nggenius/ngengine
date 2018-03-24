package object

// trigger提供了一个通用的属性触发器，在数据产生变动的时候，进行回调。
// 由使用方对具体的属性进行挂钩注册。需要注册进witness进行使用。
import (
	"fmt"
)

const (
	FLAG_ALTER         = 0x1 // 变动通知
	FLAG_REALTIME      = 0x2 // 实时的
	FLAG_LOD           = 0x4 // LOD控制
	FLAG_ALTER_RUNNING = 0x8 // 变动回调中
)

// 属性变动触发器
type AttrTrigger struct {
	object       Object
	flag         []byte
	attrTrigger  map[string]*AttrNotifier
	tableTrigger map[string]*TableNotifier
}

// 构造函数
func NewAttrTrigger() *AttrTrigger {
	o := &AttrTrigger{}
	return o
}

// 初始化，被witness回调
func (a *AttrTrigger) Init(object Object) {
	if a.object != object { // 只初始化一次
		attrs := object.AllAttr()
		a.object = object
		a.flag = make([]byte, len(attrs))
		a.attrTrigger = make(map[string]*AttrNotifier)
		a.tableTrigger = make(map[string]*TableNotifier)
	}
}

// 增加某个属性的回调
func (a *AttrTrigger) AddCallback(attr string, cbname string, cb AttrAlter) error {
	index := a.object.AttrIndex(attr)
	if index == -1 {
		return fmt.Errorf("attr not found %s", attr)
	}

	if a.object.GetAttrType(attr) == "table" {
		return fmt.Errorf("attr is table %s", attr)
	}

	a.flag[index] |= FLAG_ALTER

	if _, has := a.attrTrigger[attr]; !has {
		a.attrTrigger[attr] = NewAttrNotifier()
	}

	return a.attrTrigger[attr].Add(cbname, cb)
}

// 移除某个属性的回调
func (a *AttrTrigger) RemoveCallback(attr string, cbname string) error {
	index := a.object.AttrIndex(attr)
	if index == -1 {
		return fmt.Errorf("attr not found %s", attr)
	}

	if _, has := a.attrTrigger[attr]; !has {
		return fmt.Errorf("attr callback not found %s", attr)
	}

	return a.attrTrigger[attr].Remove(cbname)
}

// 增加某个表格的回调
func (a *AttrTrigger) AddTableCallback(table string, cbname string, cb TableAlter) error {
	index := a.object.AttrIndex(table)
	if index == -1 {
		return fmt.Errorf("attr not found %s", table)
	}

	if a.object.GetAttrType(table) != "table" {
		return fmt.Errorf("attr is not table %s", table)
	}

	a.flag[index] |= FLAG_ALTER

	if _, has := a.tableTrigger[table]; !has {
		a.tableTrigger[table] = NewTableNotifier()
	}

	return a.tableTrigger[table].Add(cbname, cb)
}

// 移除某个表格的回调
func (a *AttrTrigger) RemoveTableCallback(attr string, cbname string) error {
	index := a.object.AttrIndex(attr)
	if index == -1 {
		return fmt.Errorf("attr not found %s", attr)
	}

	if _, has := a.tableTrigger[attr]; !has {
		return fmt.Errorf("attr callback not found %s", attr)
	}

	return a.tableTrigger[attr].Remove(cbname)
}

// 属性变动时的回调函数，由witness回调
func (a *AttrTrigger) UpdateAttr(attr string, val interface{}, old interface{}) {
	index := a.object.AttrIndex(attr)
	if index == -1 {
		panic("attr not found " + attr)
	}

	if a.flag[index]&FLAG_ALTER == 0 { //没有回调
		return
	}

	if trigger, has := a.attrTrigger[attr]; has {
		if a.flag[index]&FLAG_ALTER_RUNNING == 0 {
			a.flag[index] |= FLAG_ALTER_RUNNING
			trigger.Invoke(a.object, attr, val, old)
			a.flag[index] &= ^byte(FLAG_ALTER_RUNNING)
		}
	}
}

// tuple属性变动时的回调函数，由witness回调
func (a *AttrTrigger) UpdateTuple(attr string, val interface{}, old interface{}) {

	index := a.object.AttrIndex(attr)
	if index == -1 {
		panic("attr tuple not found " + attr)
	}

	if a.flag[index]&FLAG_ALTER == 0 { //没有回调
		return
	}

	if trigger, has := a.attrTrigger[attr]; has {
		if a.flag[index]&FLAG_ALTER_RUNNING == 0 {
			a.flag[index] |= FLAG_ALTER_RUNNING
			trigger.Invoke(a.object, attr, val, old)
			a.flag[index] &= ^byte(FLAG_ALTER_RUNNING)
		}
	}
}

// table变动时的回调函数，由witness回调
func (a *AttrTrigger) UpdateTable(table string, op_type, row, col int) {
	index := a.object.AttrIndex(table)
	if index == -1 {
		panic("attr table not found " + table)
	}

	if a.flag[index]&FLAG_ALTER == 0 { //没有回调
		return
	}

	if trigger, has := a.tableTrigger[table]; has {
		if a.flag[index]&FLAG_ALTER_RUNNING == 0 {
			a.flag[index] |= FLAG_ALTER_RUNNING
			trigger.Invoke(a.object, table, op_type, row, col)
			a.flag[index] &= ^byte(FLAG_ALTER_RUNNING)
		}
	}
}
