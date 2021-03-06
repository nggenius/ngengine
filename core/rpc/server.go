package rpc

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/nggenius/ngengine/logger"
	"github.com/nggenius/ngengine/protocol"
	"github.com/nggenius/ngengine/share"
	"github.com/nggenius/ngengine/utils"
	"io"
	"strings"
	"sync"
	"time"
)

var (
	ErrTooLong    = errors.New("message is to long")
	rpccallCache  = make(chan *RpcCall, 256)
	headerCache   = make(chan *Header, 256)
	responseCache = make(chan *Response, 32)
)

type CB func(src Mailbox, dest Mailbox, msg *protocol.Message) (int32, *protocol.Message)

type RpcRegister interface {
	RegisterCallback(Servicer)
}

type Servicer interface {
	RegisterCallback(string, CB)
	ThreadPush(*RpcCall) bool
}

func NewRpcCall() *RpcCall {
	var call *RpcCall
	select {
	case call = <-rpccallCache:
	default:
		call = &RpcCall{}
	}
	return call
}

func NewHeader() *Header {
	var header *Header
	select {
	case header = <-headerCache:
	default:
		header = &Header{}
	}
	return header
}

func NewResponse() *Response {
	var resp *Response
	select {
	case resp = <-responseCache:
	default:
		resp = &Response{}
	}
	return resp
}

type Header struct {
	ServiceMethod string // format: "Service.Method"
	Seq           uint64 // sequence number chosen by client
	Src           uint64
	Dest          uint64
}

func (h *Header) Free() {
	h.ServiceMethod = ""
	h.Seq = 0
	h.Src = 0
	h.Dest = 0
	select {
	case headerCache <- h:
	default:
	}
}

type Response struct {
	Seq     uint64
	Errcode int32
	Reply   *protocol.Message
	cb      ReplyCB
}

func (r *Response) Free() {
	if r.Reply != nil {
		r.Reply.Free()
		r.Reply = nil
		r.cb = nil
	}
	select {
	case responseCache <- r:
	default:
	}
}

type RpcCall struct {
	session uint64
	srv     *service
	svr     *Server
	header  *Header
	message *protocol.Message
	reply   *protocol.Message
	errcode int32
	method  CB
	cb      ReplyCB
	cbparam interface{}
}

func (call *RpcCall) Error() bool {
	return call.errcode != 0
}

func (call *RpcCall) Call() error {
	sender := NewMailboxFromUid(call.header.Src)
	dest := NewMailboxFromUid(call.header.Dest)
	call.errcode, call.reply = call.method(sender, dest, call.message)
	return nil
}

func (call *RpcCall) GetSrc() Mailbox {
	mb := NewMailboxFromUid(call.header.Src)
	return mb
}

func (call *RpcCall) GetMethod() string {
	return call.header.ServiceMethod
}

func (call *RpcCall) IsThreadWork() bool {
	return call.srv.ThreadPush(call)
}

func (call *RpcCall) Free() error {
	call.message.Free()
	call.header.Free()
	call.header = nil
	call.message = nil
	call.srv = nil
	call.cb = nil
	call.cbparam = nil
	call.session = 0
	call.errcode = 0
	if call.reply != nil {
		call.reply.Free()
		call.reply = nil
	}
	select {
	case rpccallCache <- call:
	default:
	}
	return nil
}

func (call *RpcCall) Done() error {
	if call.header.Seq != 0 || call.cb != nil {
		return call.svr.sendResponse(call)
	}
	return nil
}

type service struct {
	svr    *Server
	rcvr   interface{}
	method map[string]CB
}

func (srv *service) RegisterCallback(name string, cb CB) {
	srv.method[name] = cb
}

func (srv *service) ThreadPush(call *RpcCall) bool {
	if t, ok := srv.rcvr.(ThreadHandler); ok {
		return t.NewJob(call)
	}
	return false
}

type Session struct {
	rwc       io.ReadWriteCloser
	codec     ServerCodec
	sendQueue chan *Response
	quit      bool
}

type Server struct {
	mutex      sync.RWMutex
	serial     uint64
	serviceMap map[string]*service
	sessions   map[uint64]*Session
	sendQueue  chan *Response
	ch         chan *RpcCall
	log        *logger.Log
}

func (server *Server) getCall(servicemethod string, src, dest Mailbox, cb ReplyCB, args ...interface{}) (*RpcCall, error) {
	var msg *protocol.Message
	if len(args) > 0 {
		msg = protocol.NewMessage(share.MAX_BUF_LEN)
		ar := utils.NewStoreArchiver(msg.Body)
		for i := 0; i < len(args); i++ {
			ar.Put(args[i])
		}
		msg.Body = msg.Body[:ar.Len()]
	} else {
		msg = protocol.NewMessage(1)
	}

	call := NewRpcCall()
	call.header = NewHeader()
	call.message = msg
	call.header.Seq = 0
	call.header.Src = src.Uid()
	call.header.Dest = dest.Uid()
	call.header.ServiceMethod = servicemethod
	call.cb = cb
	dot := strings.LastIndex(call.header.ServiceMethod, ".")
	if dot < 0 {
		err := fmt.Errorf("rpc: service/method request ill-formed: %s", call.header.ServiceMethod)
		server.log.LogErr(err)
		call.Free()
		return nil, err
	}
	serviceName := call.header.ServiceMethod[:dot]
	methodName := call.header.ServiceMethod[dot+1:]

	call.srv = server.serviceMap[serviceName]
	if call.srv == nil {
		err := fmt.Errorf("rpc: can't find service %s", call.header.ServiceMethod)
		server.log.LogErr(err)
		call.Free()
		return nil, err
	}

	call.method = call.srv.method[methodName]
	if call.method == nil {
		err := fmt.Errorf("rpc: can't find method %s", call.header.ServiceMethod)
		server.log.LogErr(err)
		call.Free()
		return nil, err
	}

	return call, nil
}

func (server *Server) Call(servicemethod string, src, dest Mailbox, args ...interface{}) error {
	call, err := server.getCall(servicemethod, src, dest, nil, args...)
	if call != nil {
		server.ch <- call
	}
	return err
}

func (server *Server) CallBack(servicemethod string, src, dest Mailbox, cb ReplyCB, cbparam interface{}, args ...interface{}) error {
	call, err := server.getCall(servicemethod, src, dest, cb, args...)
	if call != nil {
		call.cbparam = cbparam
		server.ch <- call
	}
	return err
}

func (server *Server) ServeConn(conn io.ReadWriteCloser, maxlen uint16) {
	codec := &byteServerCodec{
		rwc:    conn,
		encBuf: bufio.NewWriter(conn),
	}
	server.ServeCodec(codec, maxlen)
}

func (server *Server) ServeCodec(codec ServerCodec, maxlen uint16) {
	var serial uint64
	server.mutex.Lock()
	serial = server.serial
	server.serial++
	session := &Session{rwc: codec.GetConn(), codec: codec, sendQueue: make(chan *Response, 32)}
	server.sessions[serial] = session
	server.mutex.Unlock()
	go session.send(server.log)
	server.log.LogInfo("start new rpc server{serial:", serial, "}")
	for {
		msg, err := codec.ReadRequest(maxlen)
		if err != nil {
			if err != io.EOF &&
				!strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host") &&
				!strings.Contains(err.Error(), "use of closed network connection") {
				server.log.LogErr("rpc err:", err)
			} else {
				server.log.LogInfo("service client closed")
			}
			break
		}

		call, err := server.createCall(msg)
		if call != nil {
			if err != nil {
				call.errcode, call.reply = protocol.ReplyError(protocol.TINY, share.ERR_RPC_FAILED, err.Error())
			}
			call.svr = server
			call.session = serial
			server.ch <- call
			continue
		}
	}

	session.quit = true
	codec.Close()
	server.log.LogInfo("rpc server{serial:", serial, "} closed")
	server.mutex.Lock()
	delete(server.sessions, serial)
	server.mutex.Unlock()
}

func (session *Session) send(log *logger.Log) {
	for {
		select {
		case resp := <-session.sendQueue:

			if resp.Seq == 0 {
				resp.Free()
				log.LogFatal("resp.seq is zero")
			}

			err := session.codec.WriteResponse(resp.Seq, resp.Errcode, resp.Reply)
			if err != nil {
				resp.Free()
				log.LogErr(err)
				return
			}
			log.LogDebug("response call, seq:", resp.Seq)
			resp.Free()
		default:
			if session.quit {
				return
			}
			time.Sleep(time.Millisecond)
		}
	}
}

func (server *Server) sendResponse(call *RpcCall) error {
	if call.cb != nil { //本地回调
		if call.reply == nil && call.errcode != 0 {
			call.reply = protocol.NewMessage(1)
		}
		w := utils.NewStoreArchiver(call.reply.Header)
		w.Put(int8(1))
		w.Put(call.errcode)
		call.reply.Header = call.reply.Header[:w.Len()]
		errcode, ar := protocol.ParseReply(call.reply)
		var err *Error
		if errcode != 0 {
			err = NewError(errcode, "")
			errstr, e := ar.GetString()
			if e != nil {
				panic(e)
			}
			err.Err = errstr
		}
		call.cb(call.cbparam, err, ar)
		return nil
	}

	server.mutex.RLock()
	session := server.sessions[call.session]
	server.mutex.RUnlock()
	if session == nil {
		return fmt.Errorf("session not found %d", call.session)
	}

	resp := NewResponse()
	resp.Seq = call.header.Seq
	resp.Errcode = call.errcode
	resp.cb = call.cb
	if call.reply != nil {
		resp.Reply = call.reply.Dup()
	}
	session.sendQueue <- resp
	return nil
}

func (server *Server) createCall(msg *protocol.Message) (*RpcCall, error) {
	call := NewRpcCall()
	call.header = NewHeader()
	call.message = msg
	ar := utils.NewLoadArchiver(msg.Header)
	var err error
	call.header.Seq, err = ar.GetUint64()
	if err != nil {
		server.log.LogErr(err)
		call.Free()
		return nil, err
	}
	call.header.Src, err = ar.GetUint64()
	if err != nil {
		server.log.LogErr(err)
		call.Free()
		return nil, err
	}
	call.header.Dest, err = ar.GetUint64()
	if err != nil {
		server.log.LogErr(err)
		call.Free()
		return nil, err
	}
	call.header.ServiceMethod, err = ar.GetString()
	if err != nil {
		server.log.LogErr(err)
		call.Free()
		return nil, err
	}
	dot := strings.LastIndex(call.header.ServiceMethod, ".")
	if dot < 0 {
		err := fmt.Errorf("rpc: service/method request ill-formed: %s", call.header.ServiceMethod)
		server.log.LogErr(err)
		return call, err
	}
	serviceName := call.header.ServiceMethod[:dot]
	methodName := call.header.ServiceMethod[dot+1:]

	call.srv = server.serviceMap[serviceName]
	if call.srv == nil {
		err := fmt.Errorf("rpc: can't find service %s", call.header.ServiceMethod)
		server.log.LogErr(err)
		return call, err
	}

	call.method = call.srv.method[methodName]
	if call.method == nil {
		err := fmt.Errorf("rpc: can't find method %s", call.header.ServiceMethod)
		server.log.LogErr(err)
		return call, err
	}
	return call, nil
}

func ReadMessage(rwc io.Reader, maxrx uint16) (*protocol.Message, error) {
	var sz uint16
	var headlen uint16
	var err error
	var msg *protocol.Message

	if err = binary.Read(rwc, binary.LittleEndian, &sz); err != nil {
		return nil, err
	}

	// Limit messages to the maximum receive value, if not
	// unlimited.  This avoids a potential denaial of service.
	if sz < 0 || (maxrx > 0 && sz > maxrx) {
		return nil, ErrTooLong
	}

	if err = binary.Read(rwc, binary.LittleEndian, &headlen); err != nil {
		return nil, err
	}
	bodylen := int(sz - headlen)
	msg = protocol.NewMessage(bodylen)
	msg.Header = msg.Header[0:headlen]
	if _, err = io.ReadFull(rwc, msg.Header); err != nil {
		msg.Free()
		return nil, err
	}

	if bodylen > 0 {
		msg.Body = msg.Body[0:bodylen]

		if _, err = io.ReadFull(rwc, msg.Body); err != nil {
			msg.Free()
			return nil, err
		}
	}

	return msg, nil
}

func (server *Server) RegisterName(name string, rcvr interface{}) error {
	if reg, ok := rcvr.(RpcRegister); ok {
		srv := &service{}
		srv.svr = server
		srv.rcvr = rcvr
		srv.method = make(map[string]CB, 16)
		reg.RegisterCallback(srv)
		server.serviceMap[name] = srv
		return nil
	}

	return fmt.Errorf("%s is not RpcRegister", name)
}

func (server *Server) GetRpcInfo(name string) []string {
	var ret []string
	if s, ok := server.serviceMap[name]; ok {
		for k, _ := range s.method {
			ret = append(ret, k)
		}
	}
	return ret
}

func NewServer(ch chan *RpcCall, l *logger.Log) *Server {
	s := &Server{}
	s.log = l
	s.serviceMap = make(map[string]*service)
	s.ch = ch
	s.sessions = make(map[uint64]*Session)
	return s
}

type byteServerCodec struct {
	rwc    io.ReadWriteCloser
	encBuf *bufio.Writer
	closed bool
}

func (c *byteServerCodec) ReadRequest(maxrc uint16) (*protocol.Message, error) {
	return ReadMessage(c.rwc, maxrc)
}

func (c *byteServerCodec) WriteResponse(seq uint64, errcode int32, body *protocol.Message) (err error) {
	if body == nil {
		body = protocol.NewMessage(1)
	}

	body.Header = body.Header[:0]
	w := utils.NewStoreArchiver(body.Header)
	w.Put(seq)
	if errcode != 0 {
		w.Put(int8(1))
		w.Put(errcode)
	} else {
		w.Put(int8(0))
	}
	body.Header = body.Header[:w.Len()]
	size := len(body.Header) + len(body.Body)
	if size > RPC_BUF_LEN {
		return fmt.Errorf("message size is too large")
	}

	binary.Write(c.encBuf, binary.LittleEndian, uint16(size))             //数据大小
	binary.Write(c.encBuf, binary.LittleEndian, uint16(len(body.Header))) //头部大小
	c.encBuf.Write(body.Header)
	if len(body.Body) > 0 {
		c.encBuf.Write(body.Body)
	}
	body.Header = body.Header[:0]
	return c.encBuf.Flush()
}

func (c *byteServerCodec) Close() error {
	if c.closed {
		// Only call c.rwc.Close once; otherwise the semantics are undefined.
		return nil
	}
	c.closed = true
	return c.rwc.Close()
}

func (c *byteServerCodec) GetConn() io.ReadWriteCloser {
	return c.rwc
}

type ServerCodec interface {
	ReadRequest(maxrc uint16) (*protocol.Message, error)
	WriteResponse(seq uint64, errcode int32, body *protocol.Message) (err error)
	GetConn() io.ReadWriteCloser
	Close() error
}
