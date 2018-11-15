package service

import (
	"bufio"
	"net"
	"github.com/nggenius/ngengine/protocol"
	"github.com/nggenius/ngengine/share"
	"strconv"
)

// 和admin的连接结点
type HarborConn struct {
	net.Conn
	quit    bool
	closeCh chan struct{}
	ctx     *context
	// reading/writing interfaces
	Reader    *bufio.Reader
	Writer    *bufio.Writer
	sendqueue chan *protocol.Message
	Addr      string
	Port      int
	lenBuf    [4]byte
	lenSlice  []byte
}

// 新的连接
func NewHarborConn(conn net.Conn, ctx *context) *HarborConn {

	addr, port, _ := net.SplitHostPort(conn.RemoteAddr().String())
	p, _ := strconv.ParseInt(port, 10, 32)

	c := &HarborConn{
		ctx:       ctx,
		quit:      false,
		closeCh:   make(chan struct{}),
		Conn:      conn,
		Reader:    bufio.NewReaderSize(conn, share.MAX_BUF_LEN),
		Writer:    bufio.NewWriterSize(conn, share.MAX_BUF_LEN),
		sendqueue: make(chan *protocol.Message, 32),
		Addr:      addr,
		Port:      int(p),
	}
	c.lenSlice = c.lenBuf[:]
	return c
}

// SendMessage 发送消息
// 注意：调用方主动调用msg.Free()，SendMessage会调用msg.Dup()，这样msg会正常进入消息缓冲池
func (c *HarborConn) SendMessage(msg *protocol.Message) bool {
	if c.quit {
		return false
	}
	msg.Dup()
	c.sendqueue <- msg //消息太多的情况可能会阻塞
	return true
}

// 关闭连接
func (c *HarborConn) CloseConn() {
	if !c.quit {
		c.quit = true
		c.Close()
		close(c.closeCh)
	}
}
