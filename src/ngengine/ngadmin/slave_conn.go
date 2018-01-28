package ngadmin

import (
	"bufio"
	"net"
	"ngengine/protocol"
	"strconv"
)

type SlaveConn struct {
	net.Conn
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

func NewSlaveConn(conn net.Conn) *SlaveConn {

	addr, port, _ := net.SplitHostPort(conn.RemoteAddr().String())
	p, _ := strconv.ParseInt(port, 10, 32)

	c := &SlaveConn{
		quit:      false,
		closeCh:   make(chan struct{}),
		Conn:      conn,
		Reader:    bufio.NewReaderSize(conn, defaultBufferSize),
		Writer:    bufio.NewWriterSize(conn, defaultBufferSize),
		sendqueue: make(chan *protocol.Message, 32),
		Addr:      addr,
		Port:      int(p),
	}
	c.lenSlice = c.lenBuf[:]
	return c
}

//发送消息
func (c *SlaveConn) SendMessage(msg *protocol.Message) bool {
	if c.quit {
		return false
	}

	c.sendqueue <- msg //消息太多的情况可能会阻塞
	return true
}

//关闭连接
func (c *SlaveConn) CloseConn() {
	if !c.quit {
		c.quit = true
		c.Close()
		close(c.closeCh)
	}
}
