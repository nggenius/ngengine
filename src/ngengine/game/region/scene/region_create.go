package scene

import (
	"ngengine/core/rpc"
	"ngengine/protocol"
	"ngengine/share"
)

type RegionCreate struct {
	ctx *SceneModule
}

func NewRegionCreate(ctx *SceneModule) *RegionCreate {
	s := new(RegionCreate)
	s.ctx = ctx
	return s
}

func (s *RegionCreate) RegisterCallback(srv rpc.Servicer) {
	srv.RegisterCallback("Query", s.Query)
	srv.RegisterCallback("Create", s.Create)
}

func (s *RegionCreate) Query(src rpc.Mailbox, dest rpc.Mailbox, msg *protocol.Message) (int32, *protocol.Message) {
	return protocol.Reply(protocol.TINY, s.ctx.Core.Mailbox().ServiceId())
}

func (s *RegionCreate) Create(src rpc.Mailbox, dest rpc.Mailbox, msg *protocol.Message) (int32, *protocol.Message) {
	var r share.Region
	if err := protocol.ParseArgs(msg, &r); err != nil {
		s.ctx.Core.LogErr("parse args error")
		return 0, nil
	}
	mb, err := s.ctx.scenes.CreateScene(r)
	if err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_REGION_CREATE_FAILED, err.Error())
	}
	return protocol.Reply(protocol.TINY, mb)
}
