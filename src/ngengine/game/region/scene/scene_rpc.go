package scene

import (
	"ngengine/core/rpc"
	"ngengine/protocol"
	"ngengine/share"
)

func (s *GameScene) RegisterCallback(svr rpc.Servicer) {
	svr.RegisterCallback("Test", s.Test)
}

//srv.RegisterCallback("FunctionName", s.FunctionName)
func (s *GameScene) Test(src rpc.Mailbox, dest rpc.Mailbox, msg *protocol.Message) (int32, *protocol.Message) {
	var info string
	err := protocol.ParseArgs(msg, &info)
	if err != nil {
		s.Core().LogDebug("scene test msg, err:", err.Error())
		return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, err.Error())
	}
	s.Core().LogDebug("scene recv msg:", info)
	return 0, nil
}
