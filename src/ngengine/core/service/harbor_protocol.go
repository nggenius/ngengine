package service

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"net"
	"ngengine/protocol"
	"strings"
	"time"
)

// 与admin之间的通讯协议
type HarborProtocol struct {
	ctx       *context
	conn      *HarborConn
	connected bool
	heartbeat [4]byte
	watchs    []string
	lastload  int32
}

// 网络主循环
func (h *HarborProtocol) IOLoop(conn net.Conn) {
	binary.LittleEndian.PutUint16(h.heartbeat[:2], 2)
	binary.LittleEndian.PutUint16(h.heartbeat[2:], protocol.S2A_HEARTBEAT)
	h.conn = NewHarborConn(conn, h.ctx)

	_, err := h.Write(protocol.MagicV1)
	if err != nil {
		h.ctx.Core.LogErr("send magic failed")
		return
	}
	// 连接成功后，第一条消息就是注册自己
	if err := h.Register(); err != nil {
		h.ctx.Core.LogErr(err)
		return
	}
	h.lastload = h.ctx.Core.load
	h.connected = true
	// 发送自己关心的服务
	if err := h.Watch(); err != nil {
		h.ctx.Core.LogErr(err)
		return
	}
	// 启动发送协程
	go h.messagePump()
	// 读主循环
	for {
		id, msg, err := h.Read()
		if err != nil {
			if err != io.EOF &&
				!strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host") &&
				!strings.Contains(err.Error(), "use of closed network connection") {
				h.ctx.Core.LogErr(err)
			}
			h.Close()
			break
		}
		h.Exec(id, msg)
		msg.Free()
	}
	h.connected = false
	h.ctx.Core.LogInfo("lost admin connect")
}

// 处理消息
func (h *HarborProtocol) Exec(msgid uint16, msg *protocol.Message) {
	switch msgid {
	case protocol.A2S_SERVICES:
		{
			var srvs protocol.Services
			if err := json.Unmarshal(msg.Body, &srvs); err != nil {
				h.ctx.Core.LogErr(err)
				break
			}
			h.ctx.Core.dns.Update(srvs)
		}
	case protocol.A2S_LOAD:
		{
			var load protocol.LoadInfo
			if err := json.Unmarshal(msg.Body, &load); err != nil {
				h.ctx.Core.LogErr(err)
				break
			}
			h.ctx.Core.dns.UpdateLoad(load)
		}
	case protocol.S2A_UNREGISTER:
		{
			h.ctx.Core.Close()
		}
	}
}

// 消息发送队列处理
func (h *HarborProtocol) messagePump() {
	tick := time.NewTicker(protocol.HB_INTERVAL)
loop:
	for {
		select {
		case m := <-h.conn.sendqueue:
			h.conn.SetWriteDeadline(time.Now().Add(protocol.HB_INTERVAL * 2))
			_, err := h.conn.Writer.Write(m.Body)
			m.Free()
			if err != nil {
				h.Close()
				break
			}

			err = h.conn.Writer.Flush()
			if err != nil {
				h.Close()
				break
			}
		case <-tick.C: //
			h.conn.SetWriteDeadline(time.Now().Add(protocol.HB_INTERVAL * 2))
			_, err := h.conn.Writer.Write(h.heartbeat[:])
			if err != nil {
				h.Close()
				break
			}
			err = h.conn.Writer.Flush()
			if err != nil {
				h.Close()
				break
			}
			// 更新服务负载
			if h.lastload != h.ctx.Core.load {
				h.lastload = h.ctx.Core.load
				h.UpdateLoad()
			}
			//h.ctx.Core.LogDebug("send heartbeat")
		case <-h.conn.closeCh:
			break loop
		}
	}

exit:
	for {
		select {
		case m := <-h.conn.sendqueue:
			m.Free()
		default:
			break exit
		}
	}

	tick.Stop()
	h.ctx.Core.LogInfo("send loop quit")
}

// 关闭连接
func (h *HarborProtocol) Close() {
	if h.conn != nil {
		h.conn.CloseConn()
	}
}

// 更新负载信息
func (h *HarborProtocol) UpdateLoad() error {
	load := &protocol.LoadInfo{
		Id:   h.ctx.Core.Id,
		Load: h.ctx.Core.load,
	}
	_, err := h.WriteProtocol(protocol.S2A_LOAD, load)
	return err
}

// 注册服务协议
func (h *HarborProtocol) Register() error {
	opts := h.ctx.Core.opts
	r := &protocol.Register{}
	r.Service.Id = opts.ServId
	r.Service.Name = opts.ServName
	r.Service.Type = opts.ServType
	r.Service.Addr = h.ctx.Core.harbor.serviceAddr
	r.Service.Port = h.ctx.Core.harbor.servicePort
	r.Service.OuterAddr = h.ctx.Core.harbor.outerAddr
	r.Service.OuterPort = h.ctx.Core.harbor.clientPort
	r.Service.Status = 0
	r.Service.Load = h.ctx.Core.load
	_, err := h.WriteProtocol(protocol.S2A_REGISTER, r)
	return err
}

// 关闭服务协议
func (h *HarborProtocol) Watch() error {
	r := protocol.Watch{}
	r.WatchType = h.watchs
	_, err := h.WriteProtocol(protocol.S2A_WATCH, r)
	return err
}

// 按协议发送数据
func (h *HarborProtocol) WriteProtocol(msgid uint16, msg interface{}) (int, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return 0, err
	}
	size := len(data)
	m := protocol.NewMessage(size + 4)
	buff := bytes.NewBuffer(m.Body)
	binary.Write(buff, binary.LittleEndian, uint16(size+2))
	binary.Write(buff, binary.LittleEndian, msgid)
	buff.Write(data)
	m.Body = m.Body[:buff.Len()]
	if h.conn == nil || !h.conn.SendMessage(m) {
		m.Free()
		return 0, errors.New("send protocol failed")
	}
	h.ctx.Core.LogDebug("send message, ", msgid, msg)
	return len(m.Body), nil
}

// 写入原始数据
func (h *HarborProtocol) Write(data []byte) (int, error) {
	if len(data) > protocol.MAX_ADMIN_MESSAGE_SIZE {
		return 0, errors.New("msg too long")
	}

	m := protocol.NewMessage(len(data))
	m.Body = append(m.Body, data...)
	if !h.conn.SendMessage(m) {
		m.Free()
		return 0, errors.New("send data failed")
	}

	return len(data), nil
}

// 读取消息
func (h *HarborProtocol) Read() (uint16, *protocol.Message, error) {
	var size, msgid uint16
	if _, err := io.ReadFull(h.conn.Reader, h.conn.lenSlice); err != nil {
		return 0, nil, err
	}

	size = binary.LittleEndian.Uint16(h.conn.lenSlice[:2])
	if size > protocol.MAX_ADMIN_MESSAGE_SIZE {
		return 0, nil, errors.New("message size exceed")
	}

	msgid = binary.LittleEndian.Uint16(h.conn.lenSlice[2:])

	msg := protocol.NewMessage(int(size) - 2)
	msg.Body = msg.Body[:int(size)-2]
	if _, err := io.ReadFull(h.conn.Reader, msg.Body); err != nil {
		msg.Free()
		return 0, nil, err
	}
	return msgid, msg, nil
}
