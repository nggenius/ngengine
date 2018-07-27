package ngadmin

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"ngengine/protocol"
	"strings"
	"time"
)

type AdminProtocol struct {
	ctx *Context
}

func (p *AdminProtocol) IOLoop(conn net.Conn) error {
	var zeroTime time.Time
	adminid, peer, err := p.Register(conn)
	if err != nil {
		return err
	}
	client := newClient(peer.ServId, conn, p.ctx)
	srv := NewServ(adminid, peer, client)
	if err := p.ctx.ngadmin.DB.AddService(peer.ServId, srv); err != nil {
		return err
	}
	//发送消息进程
	go p.messagePump(client)
	var size, msgid uint16
	for {
		if client.HeartbeatInterval > 0 {
			client.SetReadDeadline(time.Now().Add(client.HeartbeatInterval * 2))
		} else {
			client.SetReadDeadline(zeroTime)
		}

		_, err := io.ReadFull(client.Reader, client.lenSlice[:])
		if err != nil {
			if err != io.EOF &&
				!strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host") {
				p.ctx.ngadmin.LogErr("read err:", err)
			}
			break
		}
		size = binary.LittleEndian.Uint16(client.lenSlice[:2])

		if size < 2 || size > 0x1000 {
			p.ctx.ngadmin.LogErrf("message size %d exceed", size)
			break
		}

		msgid = binary.LittleEndian.Uint16(client.lenSlice[2:])
		size = size - 2
		if size > 0 {
			msg := protocol.NewMessage(int(size))
			msg.Body = msg.Body[:size]
			if n, err := io.ReadFull(client.Reader, msg.Body); err != nil || uint16(n) != size {
				if err != io.EOF &&
					!strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host") &&
					!strings.Contains(err.Error(), "use of closed network connection") {
					p.ctx.ngadmin.LogErr(err.Error())
				}
				break
			}
			p.Exec(srv, msgid, msg)
			msg.Free()
		} else {
			p.Exec(srv, msgid, nil)
		}

	}
	p.ctx.ngadmin.DB.RemoveService(peer.ServName, peer.ServId)
	return nil
}

func (p *AdminProtocol) Register(conn net.Conn) (int, *PeerInfo, error) {
	var buf [4]byte
	if _, err := io.ReadFull(conn, buf[:]); err != nil {
		return 0, nil, err
	}
	var size, msgid uint16
	size = binary.LittleEndian.Uint16(buf[:2])

	if size <= 2 || size > protocol.MAX_ADMIN_MESSAGE_SIZE {
		return 0, nil, fmt.Errorf("message size %d exceed", size)
	}

	msgid = binary.LittleEndian.Uint16(buf[2:])

	if msgid != protocol.S2A_REGISTER {
		return 0, nil, fmt.Errorf("first message must register")
	}

	size = size - 2
	msg := protocol.NewMessage(int(size))
	msg.Body = msg.Body[:size]
	if _, err := io.ReadFull(conn, msg.Body[:size]); err != nil {
		return 0, nil, err
	}

	var reg protocol.Register
	if err := json.Unmarshal(msg.Body, &reg); err != nil {
		return 0, nil, err
	}
	msg.Free()
	pi := &PeerInfo{
		ServId:     reg.Service.Id,
		ServName:   reg.Service.Name,
		ServType:   reg.Service.Type,
		Status:     reg.Service.Status,
		RemoteAddr: reg.Service.Addr,
		RemotePort: reg.Service.Port,
		OuterAddr:  reg.Service.OuterAddr,
		OuterPort:  reg.Service.OuterPort,
		Load:       reg.Service.Load,
	}

	return reg.AdminId, pi, nil
}

func (p *AdminProtocol) Exec(srv *ServiceInfo, msgid uint16, msg *protocol.Message) {
	switch msgid {
	case protocol.S2A_WATCH:
		{
			var watchs protocol.Watch
			if err := json.Unmarshal(msg.Body, &watchs); err != nil {
				p.ctx.ngadmin.LogErr(err)
				break
			}
			p.ctx.ngadmin.DB.Watch(srv.PeerInfo.ServId, watchs.WatchType)
		}
	case protocol.S2A_HEARTBEAT:
		//p.ctx.ngadmin.LogDebug("recv heartbeat")
	case protocol.S2A_LOAD:
		{
			var load protocol.LoadInfo
			if err := json.Unmarshal(msg.Body, &load); err != nil {
				p.ctx.ngadmin.LogErr(err)
				break
			}
			p.ctx.ngadmin.DB.UpdateLoad(load.Id, load.Load)
		}
	}
}

func (p *AdminProtocol) messagePump(client *Client) {
	for {
		select {
		case m := <-client.sendqueue:
			n, err := client.Writer.Write(m.Body)
			m.Free()
			if err != nil {
				p.ctx.ngadmin.LogErr("write message error", err)
			}
			if err := client.Writer.Flush(); err != nil {
				p.ctx.ngadmin.LogErr("flush message error")
			}
			p.ctx.ngadmin.LogInfof("send message to service(%d), size:%d", client.Id, n)
		case <-client.exitchan:
			goto exit
		}
	}

exit:
	for {
		select {
		case m := <-client.sendqueue: //清理发送队列
			m.Free()
		default:
			break exit
		}
	}
	p.ctx.ngadmin.LogInfo("client quit loop")
}
