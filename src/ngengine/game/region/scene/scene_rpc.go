package scene

import (
	"ngengine/core/rpc"
	"ngengine/game/gameobject"
	"ngengine/protocol"
	"ngengine/share"
)

func (s *GameScene) RegisterCallback(svr rpc.Servicer) {
	svr.RegisterCallback("AddPlayer", s.AddPlayer)
}

func (s *GameScene) AddPlayer(src rpc.Mailbox, dest rpc.Mailbox, msg *protocol.Message) (int32, *protocol.Message) {
	s.Core().LogDebug("add player")
	var data []byte
	err := protocol.ParseArgs(msg, &data)
	if err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, err.Error())
	}

	obj, err := s.factory.Decode(data)
	if err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_ENTER_REGION_FAILED, err.Error())
	}

	s.addPlayer(obj.(gameobject.GameObject))

	s.Core().LogDebug("add player succeed")
	return 0, nil
}
