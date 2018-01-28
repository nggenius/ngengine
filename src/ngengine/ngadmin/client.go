package ngadmin

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"ngengine/protocol"
	"ngengine/share"
	"strconv"
	"time"
)

const defaultBufferSize = 16 * 1024

type Client struct {
	net.Conn
	Id       share.ServiceId
	quit     bool
	exitchan chan struct{}
	// reading/writing interfaces
	Reader            *bufio.Reader
	Writer            *bufio.Writer
	HeartbeatInterval time.Duration
	sendqueue         chan *protocol.Message
	Addr              string
	Port              int
	lenBuf            [4]byte
	lenSlice          []byte
}

func newClient(id share.ServiceId, conn net.Conn, ctx *Context) *Client {

	addr, port, _ := net.SplitHostPort(conn.RemoteAddr().String())
	p, _ := strconv.ParseInt(port, 10, 32)

	c := &Client{
		Id:                id,
		Conn:              conn,
		Reader:            bufio.NewReaderSize(conn, defaultBufferSize),
		Writer:            bufio.NewWriterSize(conn, defaultBufferSize),
		HeartbeatInterval: time.Duration(ctx.ngadmin.opts.HeartTimeout) * time.Second,
		sendqueue:         make(chan *protocol.Message, 32),
		exitchan:          make(chan struct{}),
		Addr:              addr,
		Port:              int(p),
	}
	c.lenSlice = c.lenBuf[:]
	return c
}

func (c *Client) SendMessage(msg *protocol.Message) bool {
	if c.quit {
		return false
	}

	c.sendqueue <- msg //消息太多的情况可能会阻塞
	return true
}

func (c *Client) SendProtocol(msgid uint16, msg interface{}) (bool, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return false, nil
	}
	size := len(data)
	m := protocol.NewMessage(size + 4)
	buff := bytes.NewBuffer(m.Body)
	binary.Write(buff, binary.LittleEndian, uint16(size+2))
	binary.Write(buff, binary.LittleEndian, msgid)
	buff.Write(data)
	m.Body = m.Body[:buff.Len()]
	if !c.SendMessage(m) {
		m.Free()
		return false, fmt.Errorf("send message failed, reason: client is quit")
	}
	return true, nil
}

func (c *Client) Quit() {
	if !c.quit {
		c.quit = true
		close(c.exitchan)
		c.Close()
	}
}
