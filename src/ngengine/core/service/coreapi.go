package service

import (
	"errors"
	"fmt"
	"ngengine/core/rpc"
	"ngengine/logger"
	"ngengine/protocol"
	"ngengine/share"
	"ngengine/utils"
	"time"
)

var (
	magic_time, _ = time.Parse("2006-01-02 15:04:05", "2018-01-01 00:00:00")
)

type SrvInfo struct {
	Id        share.ServiceId // 服务ID
	Name      string          // 服务名称
	Type      string          // 服务类型
	Status    int8            // 状态
	Addr      string          // ip地址
	Port      int             // 端口号
	OuterAddr string          // 外网地址
	OuterPort int             // 外网端口
	Load      int32           // 负载
}

// Core接口
type CoreAPI interface {
	// 获取当前服务的goroutine id
	GID() int64
	// 获取当前service
	Service() Service
	// 获取当前服务的mailbox
	Mailbox() rpc.Mailbox
	// 返回配置选项
	Option() *CoreOption
	// 关闭服务
	Shut()
	// SendReady 发送Ready消息
	SendReady()
	// 关注其它服务，"all" 关注全部服务
	Watch(...string)
	// 返回服务相关的时间
	Time() Time
	// 发起远程调用
	Mailto(src *rpc.Mailbox, dest *rpc.Mailbox, method string, args ...interface{}) error
	// 发起远程调用并调用回调函数
	MailtoAndCallback(src *rpc.Mailbox, dest *rpc.Mailbox, method string, cb rpc.ReplyCB, cbparam interface{}, args ...interface{}) error
	// 通过服务ID获取服务信息
	LookupService(id share.ServiceId) *Srv
	// 获取一个特定类型的服务
	LookupOneServiceByType(typ string) *Srv
	// 获取一个负载最小的特定类型的服务
	LookupMinLoadByType(typ string) *Srv
	// 随机获取一个特定类型的服务
	LookupRandServiceByType(typ string) *Srv
	// 获取所有服务信息
	LookupAllServiceByType(typ string) []*Srv
	// 通过服务名获取服务信息
	LookupServiceByName(name string) *Srv
	// 日志函数
	LogDebug(v ...interface{})
	// 日志函数
	LogInfo(v ...interface{})
	// 日志函数
	LogWarn(v ...interface{})
	// 日志函数
	LogErr(v ...interface{})
	// 日志函数
	LogFatal(v ...interface{})
	// 日志函数
	LogDebugf(format string, v ...interface{})
	// 日志函数
	LogInfof(format string, v ...interface{})
	// 日志函数
	LogWarnf(format string, v ...interface{})
	// 日志函数
	LogErrf(format string, v ...interface{})
	// 日志函数
	LogFatalf(format string, v ...interface{})
	// 注册C2S模块
	RegisterHandler(name string, handler interface{})
	// 注册S2S模块
	RegisterRemote(name string, remote interface{})
	// 增加模块
	AddModule(m ModuleHandler) error
	// 获取模块
	MustModule(name string) interface{}
	// 调用模块
	Call(module string, id int, args ...interface{}) error
	// 获取log指针
	Logger() *logger.Log
	// 消息协议解码
	ParseProto(msg *protocol.Message, obj interface{}) error
	// 更新负载信息
	UpdateLoad(load int32)
	// 获取负载信息
	Load() int32
	// 断开client连接
	Break(session uint64)
	// 生成GUID
	GenerateGUID() int64
}

// 获取当前服务的goroutine id
func (c *Core) GID() int64 {
	return c.gid
}

// 获取当前service
func (c *Core) Service() Service {
	return c.service
}

// 获取当前服务的mailbox
func (c *Core) Mailbox() rpc.Mailbox {
	return c.mailbox
}

// 返回配置选项
func (c *Core) Option() *CoreOption {
	return c.opts
}

// 关闭服务
func (c *Core) Shut() {
	if c.closeState >= CS_QUIT {
		return
	}

	// 踢出所有客户端连接
	if c.clientDB != nil {
		c.clientDB.CloseAll()
	}

	c.closeState = CS_QUIT

	// 给admin发关闭消息
	s := &protocol.SeverClosing{}
	s.ID = uint16(c.Id)
	s.SeverName = c.opts.ServName
	c.harbor.protocol.WriteProtocol(protocol.S2A_UNREGISTER, s)

	// 关闭harbor
	if c.harbor != nil {
		c.harbor.Close()
	}

	c.notifyDone()
}

// SendReady 发送Ready消息
func (c *Core) SendReady() {
	if c.IsReady {
		return
	}
	c.IsReady = true
	if c.harbor.protocol != nil && c.harbor.protocol.connected {
		c.harbor.protocol.WriteProtocol(protocol.S2A_READY, nil)
	}
}

// 关注其它服务，"all" 关注全部服务
func (c *Core) Watch(w ...string) {
	c.watchs = c.watchs[:0]
	c.watchs = append(c.watchs, w...)
}

// 返回服务相关的时间
func (c *Core) Time() Time {
	return *c.time
}

// 发起远程调用
func (c *Core) Mailto(src *rpc.Mailbox, dest *rpc.Mailbox, method string, args ...interface{}) error {
	if dest == nil {
		return errors.New("dest is nil")
	}

	if src == nil {
		src = &c.mailbox
	}

	// 对象
	if dest.IsObject() {
		return c.ObjectCall(src, dest, method, args...)
	}

	if !dest.IsClient() { // 判断是否是客户端的消息
		if dest.ServiceId() == c.mailbox.ServiceId() { // 本地调用
			return c.rpcSvr.Call(rpc.GetServiceMethod(method), *src, *dest, args...)
		}
		srv := c.dns.LookupByMailbox(*dest)
		if srv == nil {
			return errors.New("service not found")
		}
		return srv.Call(*src, *dest, method, args...)
	}

	if len(args) == 0 {
		return errors.New("args is zero")
	}

	if err := c.ClientCall(src, dest, method, args[0]); err != nil {
		c.LogErr(err)
		return err
	}

	return nil
}

// 发起远程调用并调用回调函数
func (c *Core) MailtoAndCallback(src *rpc.Mailbox, dest *rpc.Mailbox, method string, cb rpc.ReplyCB, cbparam interface{}, args ...interface{}) error {
	if dest == nil {
		return errors.New("dest is nil")
	}

	if src == nil {
		src = &c.mailbox
	}

	// 对象
	if dest.IsObject() {
		return c.ObjectCallback(src, dest, method, cb, cbparam, args...)
	}

	if dest.IsClient() { // 客户端的调用不支持回调
		return fmt.Errorf("client not support callback")
	}

	if dest.ServiceId() == c.mailbox.ServiceId() { // 本地调用
		return c.rpcSvr.CallBack(rpc.GetServiceMethod(method), *src, *dest, cb, cbparam, args...)
	}

	srv := c.dns.LookupByMailbox(*dest)
	if srv == nil {
		return errors.New("service not found")
	}

	if err := srv.Callback(*src, *dest, method, cb, cbparam, args...); err != nil {
		c.LogErr(err)
		return err
	}

	return nil
}

// ObjectCall 向对象发送消息
func (c *Core) ObjectCall(src *rpc.Mailbox, dest *rpc.Mailbox, method string, args ...interface{}) (err error) {
	msg := protocol.NewMessage(share.MAX_BUF_LEN)
	defer msg.Free()
	ar := utils.NewStoreArchiver(msg.Body)
	for i := 0; i < len(args); i++ {
		err = ar.Put(args[i])
		if err != nil {
			return
		}
	}
	msg.Body = msg.Body[:ar.Len()]

	c.LogInfo(src, " call ", dest, "/", method)
	// 本地
	if dest.ServiceId() == c.mailbox.ServiceId() {
		err = c.rpcSvr.Call(rpc.GetServiceMethod(share.ROUTER_TO_OBJECT), *src, *dest, method, msg.Body)
		return
	}

	srv := c.dns.LookupByMailbox(*dest)
	if srv == nil {
		return errors.New("service not found")
	}

	err = srv.Call(*src, *dest, share.ROUTER_TO_OBJECT, method, msg.Body)
	return
}

// ObjectCallback 向对象发送消息
func (c *Core) ObjectCallback(src *rpc.Mailbox, dest *rpc.Mailbox, method string, cb rpc.ReplyCB, cbparam interface{}, args ...interface{}) (err error) {
	msg := protocol.NewMessage(share.MAX_BUF_LEN)
	defer msg.Free()
	ar := utils.NewStoreArchiver(msg.Body)
	for i := 0; i < len(args); i++ {
		err = ar.Put(args[i])
		if err != nil {
			return
		}
	}

	msg.Body = msg.Body[:ar.Len()]

	c.LogInfo(src, " call ", dest, "/", method)
	// 本地
	if dest.ServiceId() == c.mailbox.ServiceId() {
		err = c.rpcSvr.CallBack(rpc.GetServiceMethod(share.ROUTER_TO_OBJECT), *src, *dest, cb, method, msg.Body)
		return
	}

	srv := c.dns.LookupByMailbox(*dest)
	if srv == nil {
		return errors.New("service not found")
	}

	err = srv.Callback(*src, *dest, share.ROUTER_TO_OBJECT, cb, cbparam, method, msg.Body)
	return
}

// 向客记端发送消息
func (c *Core) ClientCall(src *rpc.Mailbox, dest *rpc.Mailbox, method string, args interface{}) error {

	var err error
	var pb protocol.S2CMsg
	pb.Sender = c.opts.ServName
	pb.To = dest.Id()
	pb.Method = method

	if pb.Data, err = c.Proto.CreateRpcMessage(c.opts.ServName, method, args); err != nil {
		c.LogErr(err)
		return err
	}

	if src == nil {
		src = &c.mailbox
	}

	if dest.ServiceId() == c.mailbox.ServiceId() {
		msg := protocol.NewProtoMessage()
		msg.Put(pb)
		msg.Flush()
		c.s2chelper.Call(*src, rpc.NullMailbox, msg.GetMessage())
		msg.Free()
		return nil
	}

	srv := c.dns.LookupByMailbox(*dest)
	if srv == nil {
		return errors.New("service not found")
	}

	err = srv.Call(*src, rpc.NullMailbox, "S2CHelper.Call", pb)
	if err == rpc.ErrShutdown {
		srv.Close()
		c.LogErr(err)
	}

	return err
}

// 查找服务
func (c *Core) LookupService(id share.ServiceId) *Srv {
	return c.dns.Lookup(id)
}

func (c *Core) LookupMinLoadByType(typ string) *Srv {
	return c.dns.LookupMinLoadByType(typ)
}

// 获取某个类型的一个服务
func (c *Core) LookupOneServiceByType(typ string) *Srv {
	return c.dns.LookupOneByType(typ)
}

// 随机获取某个类型的一个服务
func (c *Core) LookupRandServiceByType(typ string) *Srv {
	return c.dns.LookupRandByType(typ)
}

// 查找服务
func (c *Core) LookupAllServiceByType(typ string) []*Srv {
	return c.dns.LookupByType(typ)
}

// 查找服务
func (c *Core) LookupServiceByName(name string) *Srv {
	return c.dns.LookupByName(name)
}

// 注册 c2s 模块
func (c *Core) RegisterHandler(name string, handler interface{}) {
	c.rpcmgr.RegisterHandler(name, handler)
}

// 注册 s2s 模块
func (c *Core) RegisterRemote(name string, remote interface{}) {
	c.rpcmgr.RegisterRemote(name, remote)
}

// 增加模块
func (c *Core) AddModule(m ModuleHandler) error {
	return c.modules.AddModule(m)
}

// 获取模块
func (c *Core) MustModule(module string) interface{} {
	m := c.modules.Module(module)
	if m == nil {
		panic(fmt.Errorf("must get module failed, %s", module))
	}
	return m
}

// 调用模块
func (c *Core) Call(module string, id int, args ...interface{}) error {
	m := c.modules.Module(module)
	if m == nil {
		return errors.New("module not found")
	}

	m.OnMessage(id, args...)
	return nil
}

// 获取log指针
func (c *Core) Logger() *logger.Log {
	return c.Log
}

// 消息协议解码
func (c *Core) ParseProto(msg *protocol.Message, obj interface{}) error {
	return c.Proto.DecodeMessage(msg, obj)
}

// 更新负载
func (c *Core) UpdateLoad(load int32) {
	c.load = load
}

// 负载
func (c *Core) Load() int32 {
	return c.load
}

// 断开client连接
func (c *Core) Break(session uint64) {
	c.clientDB.BreakClient(session)
}

// 生成GUID
// |63 48|47 16|15       4|3   0|
// |sid  |time |id(0~FFF) |ms   |
func (c *Core) GenerateGUID() int64 {
	c.uuid++
	dur := time.Now().Sub(magic_time).Seconds()
	ms := int64(dur*10) - int64(dur)*10
	if ms == 0 {
		ms = 1
	}
	return (int64(c.Id)&0xFFFF)<<48 |
		(int64(dur)&0xFFFFFFFF)<<16 |
		(int64(c.uuid%0xFFF)&0xFFF)<<4 |
		int64(0xF/ms)&0xF
}
