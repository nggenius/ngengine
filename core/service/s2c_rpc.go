package service

import (
	"github.com/nggenius/ngengine/core/rpc"
	"github.com/nggenius/ngengine/protocol"
	"github.com/nggenius/ngengine/share"
)

//服务器向客户端的远程调用
type S2CHelper struct {
	owner     *Core
	sendbuf   []byte
	cachedata map[int64]*protocol.Message
}

func (t *S2CHelper) RegisterCallback(s rpc.Servicer) {
	s.RegisterCallback("Call", t.Call)
}

func NewS2CHelper(core *Core) *S2CHelper {
	sc := &S2CHelper{}
	sc.owner = core
	sc.sendbuf = make([]byte, 0, share.MAX_BUF_LEN)
	sc.cachedata = make(map[int64]*protocol.Message)
	return sc
}

// Call 处理服务器向客户端的调用，对消息进行封装转成客户端的协议
func (s *S2CHelper) Call(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	request := &protocol.S2CMsg{}
	reader := protocol.NewMessageReader(msg)
	if err := reader.ReadObject(request); err != nil {
		s.owner.LogErr(err)
		return 0, nil
	}

	out, err := protocol.CreateMsg(s.sendbuf, request.Data, protocol.S2C_RPC)
	if err != nil {
		s.owner.LogErr(err)
		return 0, nil
	}

	err = s.call(sender, request.To, request.Method, out)
	if err != nil {
		s.owner.LogErr(err)
	}

	return 0, nil
}

func (s *S2CHelper) call(sender rpc.Mailbox, session uint64, method string, out []byte) error {
	client := s.owner.clientDB.FindClient(session)
	if client == nil {
		return ErrClientNotFound
	}

	s.owner.LogDebug("call ", client, " /", method)
	msg := protocol.NewMessage(len(out))
	msg.Body = append(msg.Body, out...)
	if !client.Send(msg) {
		msg.Free()
	}

	return nil
}

// Broadcast 处理服务器向客户端的广播
func (s *S2CHelper) Broadcast(sender rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	request := &protocol.S2CBrocast{}
	reader := protocol.NewMessageReader(msg)
	if err := reader.ReadObject(request); err != nil {
		s.owner.LogErr(err)
		return 0, nil
	}

	out, err := protocol.CreateMsg(s.sendbuf, request.Data, protocol.S2C_RPC)
	if err != nil {
		s.owner.LogErr(err)
		return 0, nil
	}

	for _, to := range request.To {
		if err = s.call(sender, to, request.Method, out); err != nil {
			s.owner.LogErr(to, err)
		}
	}
	return 0, nil
}
