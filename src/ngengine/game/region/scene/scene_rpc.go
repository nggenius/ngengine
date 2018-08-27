package scene

import (
	"ngengine/core/rpc"
	"ngengine/protocol"
)

func (s *GameScene) RegisterCallback(svr rpc.Servicer) {
	svr.RegisterCallback("AddPlayer", s.AddPlayer)
}

func (s *GameScene) AddPlayer(src rpc.Mailbox, dest rpc.Mailbox, msg *protocol.Message) (int32, *protocol.Message) {
	s.Core().LogDebug("add player")
	return 0, nil
}
