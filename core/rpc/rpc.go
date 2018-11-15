package rpc

import (
	"fmt"
	"net"
	"github.com/nggenius/ngengine/logger"
	"runtime"
)

const (
	RPC_BUF_LEN = 0xFFFF
)

func GetServiceMethod(m string) string {
	return fmt.Sprintf("S2S%s", m)
}

func GetHandleMethod(m string) string {
	return fmt.Sprintf("C2S%s", m)
}

func CreateRpcService(service map[string]interface{}, handle map[string]interface{}, ch chan *RpcCall, log *logger.Log) (rpcsvr *Server, err error) {
	rpcsvr = NewServer(ch, log)
	for k, v := range service {
		err = rpcsvr.RegisterName(GetServiceMethod(k), v)
		if err != nil {
			return
		}
	}

	for k, v := range handle {
		err = rpcsvr.RegisterName(GetHandleMethod(k), v)
		if err != nil {
			return
		}
	}

	return
}

func CreateService(rs *Server, l net.Listener, log *logger.Log) {
	log.LogInfo("rpc start at:", l.Addr().String())

	for _, v := range rs.serviceMap {
		if t, ok := v.rcvr.(Threader); ok {
			t.Run(log)
		}
	}

	for {
		conn, err := l.Accept()
		if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
			log.LogWarn("TCP", "temporary Accept() failure - ", err.Error())
			runtime.Gosched()
			continue
		}
		if err != nil {
			log.LogInfo("rpc service quit")
			break
		}
		//启动服务
		log.LogInfo("new rpc client,", conn.RemoteAddr())
		go rs.ServeConn(conn, RPC_BUF_LEN)
	}

	for _, v := range rs.serviceMap {
		if t, ok := v.rcvr.(Threader); ok {
			t.Terminate()
			t.WaitDone()
		}
	}
}
