package object

import (
	"ngengine/core/rpc"
	"ngengine/protocol"
	"ngengine/share"
	"reflect"
)

type ObjectService struct {
	methods map[string]*reflect.Method
}

func (o *ObjectService) Method(fname string) *reflect.Method {
	if f, ok := o.methods[fname]; ok {
		return f
	}
	return nil
}

type ObjectRegister struct {
	r   *ObjectRouter
	typ reflect.Type
	svr *ObjectService
}

func NewObjectRegister(r *ObjectRouter, t reflect.Type) *ObjectRegister {
	s := new(ObjectRegister)
	s.r = r
	s.typ = t
	svr := new(ObjectService)
	svr.methods = make(map[string]*reflect.Method)
	s.svr = svr
	return s
}

func (o *ObjectRegister) RegisterCallback(fun string, _ rpc.CB) {
	f, has := o.typ.MethodByName(fun)
	if !has {
		o.r.ctx.core.LogErr("func not found,", fun)
	}

	o.svr.methods[fun] = &f
}

func (o *ObjectRegister) ThreadPush(call *rpc.RpcCall) bool {
	return false
}

type ObjectRouter struct {
	ctx      *ObjectModule
	services map[string]*ObjectService
}

func NewObjectRouter(o *ObjectModule) *ObjectRouter {
	s := new(ObjectRouter)
	s.ctx = o
	s.services = make(map[string]*ObjectService)
	return s
}

func (s *ObjectRouter) Register(name string, object interface{}) {
	if r, ok := object.(rpc.RpcRegister); ok {
		t := reflect.TypeOf(object)
		reg := NewObjectRegister(s, t)
		sname := t.Elem().Name()
		r.RegisterCallback(reg)
		s.services[sname] = reg.svr
	}

}

func (s *ObjectRouter) RegisterCallback(srv rpc.Servicer) {
	srv.RegisterCallback("ToObject", s.ToObject)
}

func (s *ObjectRouter) ToObject(src rpc.Mailbox, dest rpc.Mailbox, msg *protocol.Message) (int32, *protocol.Message) {
	var method string
	var data []byte
	err := protocol.ParseArgs(msg, &method, &data)
	if err != nil {
		return protocol.ReplyError(protocol.DEF, share.ERR_ARGS_ERROR, err.Error())
	}

	obj, err := s.ctx.FindObject(dest)
	if err != nil {
		return protocol.ReplyError(protocol.DEF, share.ERR_OBJECT_NOT_FOUND, err.Error())
	}

	args := protocol.NewMessage(len(data))
	args.Body = append(args.Body, data...)
	defer args.Free()
	inter := reflect.ValueOf(obj)
	typ := reflect.Indirect(inter).Type().Name()
	if svr, ok := s.services[typ]; ok {
		f := svr.Method(method)
		if f != nil {
			ret := f.Func.Call([]reflect.Value{inter, reflect.ValueOf(src), reflect.ValueOf(dest), reflect.ValueOf(args)})
			return ret[0].Interface().(int32), ret[1].Interface().(*protocol.Message)
		}
	}
	return protocol.ReplyError(protocol.DEF, share.ERR_OBJECT_RPC_CALL, "")
}
