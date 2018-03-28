package replicate

import "ngengine/module/object"

const (
	FLAG_DIRTY = 0x10
)

type trigger struct {
	ctx    *ReplicateModule
	object object.Object
	flag   []byte
}

// 构造函数
func newtrigger(ctx *ReplicateModule) *trigger {
	o := &trigger{}
	o.ctx = ctx
	return o
}

// 初始化，被witness回调
func (t *trigger) Init(o object.Object) {
	if t.object != o { // 只初始化一次
		attrs := o.AllAttr()
		t.object = o
		t.flag = make([]byte, len(attrs))
		for i, a := range attrs {
			t.flag[i] = byte(o.Expose(a))
		}
	}
}

// 属性变动时的回调函数，由witness回调
func (t *trigger) UpdateAttr(attr string, val interface{}, old interface{}) {
	index := t.object.AttrIndex(attr)
	if index == -1 {
		panic("attr not found " + attr)
	}

	if t.flag[index]&object.EXPOSE_ALL == 0 { //没有回调
		return
	}

	t.ctx.core.LogInfo("replicate attr ", attr)
}

// tuple属性变动时的回调函数，由witness回调
func (t *trigger) UpdateTuple(attr string, val interface{}, old interface{}) {

	index := t.object.AttrIndex(attr)
	if index == -1 {
		panic("attr tuple not found " + attr)
	}

	if t.flag[index]&object.EXPOSE_ALL == 0 { //没有回调
		return
	}

	t.ctx.core.LogInfo("replicate tuple ", attr)
}

// table变动时的回调函数，由witness回调
func (t *trigger) UpdateTable(table string, op_type, row, col int) {
	index := t.object.AttrIndex(table)
	if index == -1 {
		panic("attr table not found " + table)
	}

	if t.flag[index]&object.EXPOSE_ALL == 0 { //没有回调
		return
	}

	t.ctx.core.LogInfo("replicate table ", table)
}
