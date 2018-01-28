package ngadmin

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"ngengine/protocol"
	"time"
)

type Slave struct {
	ctx       *Context
	host      string
	port      int
	connected bool
	conn      *SlaveConn
	quit      bool
	heartbeat [4]byte
}

func NewSlave(ctx *Context) *Slave {
	s := &Slave{
		ctx: ctx,
	}
	return s
}

func (s *Slave) KeepConnect(host string, port int) {
	s.host = host
	s.port = port
	for !s.quit {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			time.Sleep(time.Second)
		}
		s.connected = true
		s.IOLoop(conn)
		s.connected = false
	}
}

func (s *Slave) IOLoop(conn net.Conn) {
	binary.LittleEndian.PutUint16(s.heartbeat[:2], 2)
	binary.LittleEndian.PutUint16(s.heartbeat[2:], protocol.S2A_HEARTBEAT)
	s.conn = NewSlaveConn(conn)

	for {
		id, msg, err := s.Read()
		if err != nil {
			s.ctx.ngadmin.LogErr(err)
			s.Close()
			break
		}
		s.Exec(id, msg)
		msg.Free()
	}
}

func (s *Slave) messagePump() {
	tick := time.NewTicker(protocol.HB_INTERVAL)
	for {
		select {
		case m := <-s.conn.sendqueue:
			s.conn.SetWriteDeadline(time.Now().Add(protocol.HB_INTERVAL * 2))
			_, err := s.conn.Writer.Write(m.Body)
			m.Free()
			if err != nil {
				s.Close()
				break
			}

			err = s.conn.Writer.Flush()
			if err != nil {
				s.Close()
				break
			}
		case <-tick.C:
			s.conn.SetWriteDeadline(time.Now().Add(protocol.HB_INTERVAL * 2))
			_, err := s.conn.Writer.Write(s.heartbeat[:])
			if err != nil {
				s.Close()
				break
			}
			err = s.conn.Writer.Flush()
			if err != nil {
				s.Close()
				break
			}
		case <-s.conn.closeCh:
			goto exit
		}
	}

exit:
	for {
		select {
		case m := <-s.conn.sendqueue:
			m.Free()
		default:
			break exit
		}
	}
	s.ctx.ngadmin.LogInfo("send loop quit")
}

func (s *Slave) Exec(msgid uint16, msg *protocol.Message) {
	switch msgid {

	}
}

func (s *Slave) Read() (uint16, *protocol.Message, error) {
	var size, msgid uint16
	if _, err := io.ReadFull(s.conn.Reader, s.conn.lenSlice); err != nil {
		return 0, nil, err
	}

	size = binary.LittleEndian.Uint16(s.conn.lenSlice[:2])
	if size > protocol.MAX_ADMIN_MESSAGE_SIZE {
		return 0, nil, errors.New("message size exceed")
	}

	msgid = binary.LittleEndian.Uint16(s.conn.lenSlice[2:])

	msg := protocol.NewMessage(int(size) - 2)
	msg.Body = msg.Body[:int(size)-2]
	if _, err := io.ReadFull(s.conn.Reader, msg.Body); err != nil {
		msg.Free()
		return 0, nil, err
	}
	return msgid, msg, nil
}

func (s *Slave) Close() {
	if s.conn != nil {
		s.conn.CloseConn()
	}
}
