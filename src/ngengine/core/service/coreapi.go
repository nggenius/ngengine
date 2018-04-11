package service

import (
	"errors"
	"fmt"
	"ngengine/core/rpc"
	"ngengine/logger"
	"ngengine/protocol"
	"ngengine/share"
	"time"
)

type SrvInfo struct {
	Id     share.ServiceId // 服务ID
	Name   string          // 服务名称
	Type   string          // 服务类型
	Status int             // 状态
	Addr   string          // ip地址
	Port   int             // 端口号
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
	// 关注其它服务，"all" 关注全部服务
	Watch(...string)
	// 返回服务相关的时间
	Time() Time
	// 发起远程调用
	Mailto(src *rpc.Mailbox, dest *rpc.Mailbox, method string, args ...interface{}) error
	// 发起远程调用并调用回调函数
	MailtoAndCallback(src *rpc.Mailbox, dest *rpc.Mailbox, method string, cb rpc.ReplyCB, args ...interface{}) error
	// 查看服务信息
	LookupService(id share.ServiceId) *SrvInfo
	// 查看服务信息
	LookupOneServiceByType(typ string) *SrvInfo
	// 查看服务信息
	LookupAllServiceByType(typ string) []*SrvInfo
	// 查看服务信息
	LookupServiceByName(name string) *SrvInfo
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
	Module(name string) interface{}
	// 调用模块
	Call(module string, id int, args ...interface{}) error
	// 获取log指针
	Logger() *logger.Log
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
	// 关闭harbor
	if c.harbor != nil {
		c.harbor.Close()
	}
	c.Wait()

	// 关闭所有的模块
	for n, m := range c.modules.modules {
		m.Shut()
		c.LogInfo("module '", n, "' is shut")
	}

	// 给主循环足够的时间进行收尾处理
	time.Sleep(time.Second)
	c.closeState = CS_SHUT
	<-c.quitCh

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

	if !dest.IsClient() { // 判断是否是客户端的消息
		if dest.Sid == c.mailbox.Sid { // 本地调用
			return c.rpcSvr.Call(rpc.GetServiceMethod(method), *src, args...)
		}
		srv := c.dns.LookupByMailbox(*dest)
		if srv == nil {
			return errors.New("service not found")
		}
		return srv.Call(*src, method, args...)
	}

	if len(args) == 0 {
		return errors.New("args is zero")
	}

	return c.ClientCall(src, dest, method, args[0])
}

// 发起远程调用并调用回调函数
func (c *Core) MailtoAndCallback(src *rpc.Mailbox, dest *rpc.Mailbox, method string, cb rpc.ReplyCB, args ...interface{}) error {
	if dest == nil {
		return errors.New("dest is nil")
	}

	if src == nil {
		src = &c.mailbox
	}

	if dest.IsClient() { // 客户端的调用不支持回调
		return fmt.Errorf("client not support callback")
	}

	if dest.Sid == c.mailbox.Sid { // 本地调用
		return c.rpcSvr.CallBack(rpc.GetServiceMethod(method), *src, cb, args...)
	}

	srv := c.dns.LookupByMailbox(*dest)
	if srv == nil {
		return errors.New("service not found")
	}

	return srv.Callback(*src, method, cb, args...)
}

// 向客记端发送消息
func (c *Core) ClientCall(src *rpc.Mailbox, dest *rpc.Mailbox, method string, args interface{}) error {

	var err error
	var pb protocol.S2CMsg
	pb.Sender = c.opts.ServName
	pb.To = dest.Id
	pb.Method = method

	if pb.Data, err = c.Proto.CreateRpcMessage(c.opts.ServName, method, args); err != nil {
		c.LogErr(err)
		return err
	}

	if src == nil {
		src = &c.mailbox
	}

	if dest.Sid == c.mailbox.Sid {
		msg := protocol.NewProtoMessage()
		msg.Put(pb)
		msg.Flush()
		c.s2chelper.Call(*src, msg.GetMessage())
		msg.Free()
		return nil
	}

	srv := c.dns.LookupByMailbox(*dest)
	if srv == nil {
		return errors.New("service not found")
	}

	err = srv.Call(*src, "S2CHelper.Call", pb)
	if err == rpc.ErrShutdown {
		srv.Close()
	}

	return err
}

// 查找服务
func (c *Core) LookupService(id share.ServiceId) *SrvInfo {
	c.dns.RLock()
	defer c.dns.RUnlock()
	s := c.dns.Lookup(id)
	if s == nil {
		return nil
	}

	return &SrvInfo{
		Id:     s.Id,
		Name:   s.Name,
		Type:   s.Type,
		Status: s.Status,
		Addr:   s.Addr,
		Port:   s.Port,
	}

}

// 查找服务
func (c *Core) LookupOneServiceByType(typ string) *SrvInfo {
	c.dns.RLock()
	defer c.dns.RUnlock()
	s := c.dns.LookupByType(typ)
	if s == nil || len(s) == 0 {
		return nil
	}

	return &SrvInfo{
		Id:     s[0].Id,
		Name:   s[0].Name,
		Type:   s[0].Type,
		Status: s[0].Status,
		Addr:   s[0].Addr,
		Port:   s[0].Port,
	}
}

// 查找服务
func (c *Core) LookupAllServiceByType(typ string) []*SrvInfo {
	c.dns.RLock()
	defer c.dns.RUnlock()
	sa := c.dns.LookupByType(typ)
	if sa == nil || len(sa) == 0 {
		return nil
	}

	result := make([]*SrvInfo, 0, len(sa))

	for _, s := range sa {
		result = append(result, &SrvInfo{
			Id:     s.Id,
			Name:   s.Name,
			Type:   s.Type,
			Status: s.Status,
			Addr:   s.Addr,
			Port:   s.Port,
		})
	}

	return result
}

// 查找服务
func (c *Core) LookupServiceByName(name string) *SrvInfo {
	c.dns.RLock()
	defer c.dns.RUnlock()
	s := c.dns.LookupByName(name)
	if s == nil {
		return nil
	}

	return &SrvInfo{
		Id:     s.Id,
		Name:   s.Name,
		Type:   s.Type,
		Status: s.Status,
		Addr:   s.Addr,
		Port:   s.Port,
	}
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
func (c *Core) Module(module string) interface{} {
	return c.modules.Module(module)
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
