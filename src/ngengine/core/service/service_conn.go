package service

import (
	"fmt"
	"net"
	"ngengine/core/rpc"
	"ngengine/logger"
)

// 服务信息
type Srv struct {
	SrvInfo
	mb        *rpc.Mailbox
	conn      net.Conn    // 网络连接
	client    *rpc.Client // rpc client
	connected bool        // 是否已经成功连接
	l         *logger.Log // 日志
}

func (s Srv) Mailbox() *rpc.Mailbox {
	return s.mb
}

// 格式化输出
func (s Srv) String() string {
	return fmt.Sprintf("Service{Id:%d,Name:%s,Type:%s,Status:%d,Addr:%s,Port:%d,Outer addr:%s,Outer Port:%d,Load:%d}",
		s.Id, s.Name, s.Type, s.Status, s.Addr, s.Port, s.OuterAddr, s.OuterPort, s.Load)
}

// 建立连接
func (s *Srv) Connect() error {
	if s.connected {
		return nil
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", s.Addr, s.Port))
	if err != nil {
		return err
	}
	s.conn = conn
	s.client = rpc.NewClient(s.conn, s.l)
	s.connected = true

	return nil
}

// 关闭连接
func (s *Srv) Close() {
	if s.connected {
		if s.client != nil {
			s.client.Close()
			s.client = nil
		}
		if s.conn != nil {
			s.conn.Close()
			s.conn = nil
		}

		s.connected = false
	}
}

func (s *Srv) Process() {
	if s.connected && s.client != nil {
		s.client.Process()
	}
}

// 远程调用
func (s *Srv) Call(src rpc.Mailbox, method string, args ...interface{}) error {
	if !s.connected {
		if err := s.Connect(); err != nil {
			return err
		}
	}

	s.l.LogInfo("call ", src, "/", method)
	err := s.client.Call(rpc.GetServiceMethod(method), src, args...)
	if err != nil {
		s.Close()
	}

	return err
}

// 带返回函数的调用
func (s *Srv) Callback(src rpc.Mailbox, method string, cb rpc.ReplyCB, args ...interface{}) error {
	if !s.connected {
		if err := s.Connect(); err != nil {
			return err
		}
	}

	s.l.LogInfo("call ", src, "/", method)

	err := s.client.CallBack(rpc.GetServiceMethod(method), src, cb, args...)
	if err != nil {
		s.Close()
	}

	return err
}

func (s *Srv) Handle(src rpc.Mailbox, method string, args ...interface{}) error {
	if !s.connected {
		if err := s.Connect(); err != nil {
			return err
		}
	}

	s.l.LogInfo("client call ", src, "/", method)
	err := s.client.Call(rpc.GetHandleMethod(method), src, args...)
	if err != nil {
		s.Close()
	}

	return err
}
