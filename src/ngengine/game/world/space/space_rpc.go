package space

import (
	"ngengine/core/rpc"
	"ngengine/protocol"
	"ngengine/share"
)

type Space struct {
	ctx *WorldSpaceModule
}

func NewSpace(ctx *WorldSpaceModule) *Space {
	r := &Space{}
	r.ctx = ctx
	return r
}

func (r *Space) RegisterCallback(s rpc.Servicer) {
	s.RegisterCallback("CreateRegion", r.CreateRegion)
	s.RegisterCallback("FindRegion", r.FindRegion)
}

func (r *Space) CreateRegion(src rpc.Mailbox, dest rpc.Mailbox, msg *protocol.Message) (int32, *protocol.Message) {
	return 0, nil
}

func (s *Space) FindRegion(src rpc.Mailbox, dest rpc.Mailbox, msg *protocol.Message) (int32, *protocol.Message) {
	var id int
	var fx, fy, fz float64
	err := protocol.ParseArgs(msg, &id, &fx, &fy, &fz)
	if err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, err.Error())
	}

	r := s.ctx.sm.FindRegionById(id)
	if r == nil {
		return protocol.Reply(protocol.TINY, rpc.NullMailbox)
	}

	return protocol.Reply(protocol.TINY, r.Dest)
}
