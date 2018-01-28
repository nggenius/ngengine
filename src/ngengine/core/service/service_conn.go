package service

import (
	"fmt"
	"net"
	"ngengine/core/rpc"
	"ngengine/logger"
	. "ngengine/share"
)

// 服务信息
type Srv struct {
	Id        ServiceId   // 服务ID
	Name      string      // 服务名称
	Type      string      // 服务类型
	Status    int         // 状态
	Addr      string      // ip地址
	Port      int         // 端口号
	Conn      net.Conn    // 网络连接
	RpcClient *rpc.Client // rpc client
	Connected bool        // 是否已经成功连接
	l         *logger.Log // 日志
}

// 格式化输出
func (s Srv) String() string {
	return fmt.Sprintf("Service{Id:%d,Name:%s,Type:%s,Status:%d,Addr:%s,Port:%d}", s.Id, s.Name, s.Type, s.Status, s.Addr, s.Port)
}

// 建立连接
func (s *Srv) Connect() error {
	if s.Connected {
		return nil
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", s.Addr, s.Port))
	if err != nil {
		return err
	}
	s.Conn = conn
	s.RpcClient = rpc.NewClient(s.Conn, s.l)
	s.Connected = true

	return nil
}

// 关闭连接
func (s *Srv) Close() {
	if s.Connected {
		if s.RpcClient != nil {
			s.RpcClient.Close()
			s.RpcClient = nil
		}
		if s.Conn != nil {
			s.Conn.Close()
			s.Conn = nil
		}

		s.Connected = false
	}
}

func (s *Srv) Process() {
	if s.Connected && s.RpcClient != nil {
		s.RpcClient.Process()
	}
}

// 远程调用
func (s *Srv) Call(src rpc.Mailbox, method string, args ...interface{}) error {
	if !s.Connected {
		if err := s.Connect(); err != nil {
			return err
		}
	}

	s.l.LogInfo("call ", src, "/", method)
	err := s.RpcClient.Call(rpc.GetServiceMethod(method), src, args...)
	if err != nil {
		s.Close()
	}

	return err
}

// 带返回函数的调用
func (s *Srv) Callback(src rpc.Mailbox, method string, cb rpc.ReplyCB, args ...interface{}) error {
	if !s.Connected {
		if err := s.Connect(); err != nil {
			return err
		}
	}

	s.l.LogInfo("call ", src, "/", method)

	err := s.RpcClient.CallBack(rpc.GetServiceMethod(method), src, cb, args...)
	if err != nil {
		s.Close()
	}

	return err
}

func (s *Srv) Handle(src rpc.Mailbox, method string, args ...interface{}) error {
	if !s.Connected {
		if err := s.Connect(); err != nil {
			return err
		}
	}

	s.l.LogInfo("client call ", src, "/", method)
	err := s.RpcClient.Call(rpc.GetHandleMethod(method), src, args...)
	if err != nil {
		s.Close()
	}

	return err
}
