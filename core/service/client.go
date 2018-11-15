package service

import (
	"fmt"
	"net"
	"github.com/nggenius/ngengine/core/rpc"
	"github.com/nggenius/ngengine/protocol"
	"time"
)

type Client struct {
	Session uint64
	conn    *ClientConn
	ctx     *context
	Mailbox rpc.Mailbox
	quit    bool
}

func NewClient(id uint64, conn net.Conn, ctx *context) *Client {
	c := &Client{
		Session: id,
		conn:    NewClientConn(id, conn),
		ctx:     ctx,
		quit:    false,
	}

	return c
}

func (c *Client) String() string {
	return fmt.Sprintf("client{session:%d,addr:%s,port:%d}", c.Session, c.conn.Addr, c.conn.Port)
}

func (c *Client) Close() {
	if !c.quit {
		c.quit = true
		c.conn.Close()
	}
}

func (c *Client) Send(msg *protocol.Message) bool {
	return c.conn.SendMessage(msg)
}

func (c *Client) IOLoop() {
	go c.innerSend()
}

func (c *Client) innerSend() {
	flush := false
	for !c.quit {
		select {
		case msg := <-c.conn.sendqueue:
			c.conn.Writer.Write(msg.Body)
			c.ctx.Core.LogInfo("send message to ", c, " size: ", len(msg.Body))
			msg.Free()
			flush = true
			break
		default:
			if flush {
				c.conn.Writer.Flush()
				flush = false
			}
			time.Sleep(time.Millisecond)
			break
		}
	}

quit:
	for {
		select {
		case msg := <-c.conn.sendqueue:
			msg.Free()
		default:
			break quit
		}
	}

	c.ctx.Core.LogInfo("client quit io loop")
}
