package service

import (
	"ngengine/core/rpc"
)

type RpcRegister struct {
	// 服务器内部调用(s2s)
	remotes map[string]interface{}
	// 客户端和服务器之间的调用(c2s)
	handlers map[string]interface{}
}

func NewRpcRegister() *RpcRegister {
	r := &RpcRegister{}
	r.remotes = make(map[string]interface{})
	r.handlers = make(map[string]interface{})
	return r
}

// 通过名称获取某个内部调用
func (r *RpcRegister) GetRemote(name string) interface{} {
	if k, ok := r.remotes[name]; ok {
		return k
	}

	return nil
}

// 通过名称获取某个客户端调用
func (r *RpcRegister) GetHandler(name string) interface{} {
	if k, ok := r.handlers[name]; ok {
		return k
	}

	return nil
}

// 获取所有Handler
func (r *RpcRegister) GetAllHandler() map[string]interface{} {
	return r.handlers
}

// 注册Remote
func (r *RpcRegister) RegisterRemote(name string, remote interface{}) {
	if remote == nil {
		panic("rpc: Register remote is nil")
	}
	if _, dup := r.remotes[name]; dup {
		panic("rpc: Register called twice for remote " + name)
	}
	r.remotes[name] = remote
}

// 注册Handler
func (r *RpcRegister) RegisterHandler(name string, handler interface{}) {
	if handler == nil {
		panic("rpc: Register handler is nil")
	}
	if _, dup := r.handlers[name]; dup {
		panic("rpc: Register called twice for handler " + name)
	}
	r.handlers[name] = handler
}

func (r *RpcRegister) createRpc(ch chan *rpc.RpcCall, ctx *context) *rpc.Server {
	rpc, err := rpc.CreateRpcService(r.remotes, r.handlers, ch, ctx.Core.Log)
	if err != nil {
		ctx.Core.LogFatal(err)
	}
	return rpc
}
