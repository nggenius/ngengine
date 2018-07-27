package service

import (
	"bufio"
	"errors"
	"net"
	"ngengine/protocol"
	"ngengine/share"
	"strconv"
)

// ClientConn 和admin的连接结点
type ClientConn struct {
	conn    net.Conn
	Id      uint64
	quit    bool
	closeCh chan struct{}
	// reading/writing interfaces
	Reader    *bufio.Reader
	Writer    *bufio.Writer
	sendqueue chan *protocol.Message
	Addr      string
	Port      int
	lenBuf    [4]byte
	lenSlice  []byte
}

// NewClientConn 新的连接
func NewClientConn(id uint64, conn net.Conn) *ClientConn {

	addr, port, _ := net.SplitHostPort(conn.RemoteAddr().String())
	p, _ := strconv.ParseInt(port, 10, 32)

	c := &ClientConn{
		Id:        id,
		quit:      false,
		closeCh:   make(chan struct{}),
		conn:      conn,
		Reader:    bufio.NewReaderSize(conn, share.MAX_BUF_LEN),
		Writer:    bufio.NewWriterSize(conn, share.MAX_BUF_LEN),
		sendqueue: make(chan *protocol.Message, 32),
		Addr:      addr,
		Port:      int(p),
	}
	c.lenSlice = c.lenBuf[:]
	return c
}

// SendMessage 异步发送消息
func (c *ClientConn) SendMessage(msg *protocol.Message) bool {
	if c.quit {
		return false
	}

	c.sendqueue <- msg //消息太多的情况可能会阻塞
	return true
}

// Read 读取消息
func (c *ClientConn) Read(p []byte) (n int, err error) {
	return c.Reader.Read(p)
}

// Write 写入消息
func (c *ClientConn) Write(p []byte) (n int, err error) {
	msg := protocol.NewMessage(len(p))
	msg.Body = append(msg.Body, p...)
	if c.SendMessage(msg) {
		return len(p), nil
	}

	return 0, errors.New("socket is closed")
}

// Close 关闭连接
func (c *ClientConn) Close() error {
	if !c.quit {
		c.quit = true
		c.conn.Close()
		close(c.closeCh)
	}
	return nil
}
