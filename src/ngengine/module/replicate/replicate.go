package replicate

import (
	"ngengine/core/rpc"
	"ngengine/module/object"
)

type Replicate struct {
	ctx *ReplicateModule
}

func NewReplicate(ctx *ReplicateModule) *Replicate {
	r := &Replicate{}
	r.ctx = ctx
	return r
}

func (r *Replicate) ObjectCreate(self, sender rpc.Mailbox, args ...interface{}) int {
	o, err := r.ctx.objectmodule.FindObject(self)
	if err != nil {
		r.ctx.Core.LogErr(err)
	}

	if obj, ok := o.(object.Object); ok {
		t := newtrigger(r.ctx)
		obj.AddAttrObserver("replicate", t)
		obj.AddTableObserver("replicate", t)
	}
	r.ctx.Core.LogInfo("object created")
	return 0
}

func (r *Replicate) ObjectDestroy(self, sender rpc.Mailbox, args ...interface{}) int {
	o, err := r.ctx.objectmodule.FindObject(self)
	if err != nil {
		r.ctx.Core.LogErr(err)
	}

	if obj, ok := o.(object.Object); ok {
		obj.RemoveAttrObserver("replicate")
		obj.RemoveTableObserver("replicate")
	}

	r.ctx.Core.LogInfo("object destroy")
	return 0
}
