package service

import (
	"errors"
	"ngengine/core/rpc"
	"ngengine/protocol"
	"ngengine/share"
)

var (
	ErrClientNotFound = errors.New("client not found")
	ErrAppNotFound    = errors.New("rpc call app not found")
	packagesize       = 1400 //MTU大小
)

//客户端向服务器的远程调用辅助工具
type C2SHelper struct {
	owner *Core
}

func (t *C2SHelper) RegisterCallback(s rpc.Servicer) {
	s.RegisterCallback("Call", t.Call)
}

//处理客户端的调用
func (ch *C2SHelper) Call(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	node, sm, data, err := ch.owner.Proto.DecodeRpcMessage(msg)
	if err != nil {
		ch.owner.LogErr(err)
		return share.ERR_ARGS_ERROR, nil
	}

	if node == "." {
		err = ch.owner.rpcSvr.Call(rpc.GetHandleMethod(sm), sender, rpc.NullMailbox, data)
	} else {
		srv := ch.owner.dns.LookupByName(node)
		if srv == nil {
			ch.owner.LogErr(ErrAppNotFound)
			return share.ERR_SYSTEM_ERROR, nil
		}

		err = srv.Handle(sender, sm, data)
	}

	ch.owner.LogDebug("client call ", node, "/", sm)
	if err != nil {
		ch.owner.LogErr(err)
	}
	return 0, nil
}
