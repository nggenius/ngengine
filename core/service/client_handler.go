package service

import (
	"errors"
	"io"
	"net"
	"github.com/nggenius/ngengine/common/event"
	"github.com/nggenius/ngengine/core/rpc"
	"github.com/nggenius/ngengine/protocol"
	"github.com/nggenius/ngengine/share"
)

var (
	ERRNOTSUPPORT = errors.New("not support")
)

type ClientCodec struct {
	client   *Client
	cachebuf []byte
}

// ReadRequest 解码客户端消息
func (c *ClientCodec) ReadRequest(maxrc uint16) (*protocol.Message, error) {
	for {
		id, data, err := protocol.ReadPkg(c.client.conn.Reader, c.cachebuf)
		if err != nil {
			return nil, err
		}

		switch id {
		case protocol.C2S_PING:
			break
		case protocol.C2S_RPC:
			msg := protocol.NewMessage(len(data))
			ar := protocol.NewHeadWriter(msg)
			ar.Put(uint64(0))
			ar.Put(c.client.Mailbox.Uid())
			ar.Put(uint64(0))
			ar.Put(rpc.GetHandleMethod("C2SHelper.Call"))
			msg.Header = msg.Header[:ar.Len()]
			msg.Body = append(msg.Body, data...)
			return msg, nil
		}
	}
}

// WriteResponse 发送rcp应答，不支持
func (c *ClientCodec) WriteResponse(seq uint64, errcode int32, body *protocol.Message) (err error) {
	return ERRNOTSUPPORT
}

// Close 关闭连接
func (c *ClientCodec) Close() error {
	c.client.Close()
	return nil
}

// GetConn 获取连接
func (c *ClientCodec) GetConn() io.ReadWriteCloser {
	return c.client.conn
}

// Mailbox 获取邮箱地址
func (c *ClientCodec) Mailbox() *rpc.Mailbox {
	return &c.client.Mailbox
}

type ClientHandler struct {
	ctx *context
}

func (c *ClientHandler) Handle(conn net.Conn) {
	if c.ctx.Core.closeState != CS_NONE { //服务已经关闭
		conn.Close()
		return
	}

	id := c.ctx.Core.clientDB.AddClient(conn)
	if id == 0 {
		conn.Close()
		return
	}

	mb := rpc.NewSessionMailbox(c.ctx.Core.Id, id)
	c.ctx.Core.Emitter.Fire(share.EVENT_USER_CONNECT, event.EventArgs{"id": id}, true)
	client := c.ctx.Core.clientDB.FindClient(id)
	client.Mailbox = mb
	// loop 只处理发送到客户端数据
	client.IOLoop()
	codec := &ClientCodec{}
	codec.client = client
	codec.cachebuf = make([]byte, share.MAX_BUF_LEN)
	// 将客户端消息对接进rpc，这里是一个阻塞操作
	c.ctx.Core.rpcSvr.ServeCodec(codec, share.MAX_BUF_LEN)
	c.ctx.Core.Emitter.Fire(share.EVENT_USER_LOST, event.EventArgs{"id": id}, true)
}
