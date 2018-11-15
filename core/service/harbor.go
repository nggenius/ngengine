package service

import (
	"fmt"
	"net"
	"time"
)

// Harbor是一个连接点，提供三个服务，
// 1:保持与admin的连接，获取别的服务器的信息
// 2:维护本地服务的连接，供其它服务远程调用
// 3:提供外网服务
type Harbor struct {
	ctx             *context
	adminConn       net.Conn
	adminAddr       string // admin地址
	adminPort       int    // admin端口
	serviceListener net.Listener
	serviceAddr     string // 服务地址
	servicePort     int    // 服务端口
	clientListener  net.Listener
	outerAddr       string // 外网地址
	clientAddr      string // 供客户端连接的监听地址
	clientPort      int    // 供客户端连接的端口
	quit            bool
	watchs          []string
	protocol        *HarborProtocol
}

// NewHarbor 创建一个新的Harbor
func NewHarbor(ctx *context) *Harbor {
	h := &Harbor{
		ctx:  ctx,
		quit: false,
	}
	return h
}

// 设置admin的地址
func (h *Harbor) SetAdmin(addr string, port int) {
	h.adminAddr = addr
	h.adminPort = port
}

// KeepConnect 保持与admin的连接
func (h *Harbor) KeepConnect() {
	retrydelay := time.Second
	retrytimes := 0
	for !h.quit {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", h.adminAddr, h.adminPort))
		if err != nil {
			h.ctx.Core.LogErrf("connect to admin(%s:%d) failed", h.adminAddr, h.adminPort)
			time.Sleep(retrydelay)
			if retrytimes > 60 { //第一分钟每秒连接一次
				retrydelay *= 2
				if retrydelay >= time.Minute {
					retrydelay = time.Minute
				}
			}
			retrytimes++
			h.ctx.Core.LogInfof("reconnect admin, times: %d", retrytimes)
			continue
		}
		retrydelay = time.Second
		retrytimes = 0
		h.adminConn = conn
		hprot := &HarborProtocol{
			ctx:    h.ctx,
			watchs: h.watchs,
		}

		h.protocol = hprot
		hprot.IOLoop(h.adminConn)
		h.protocol = nil

	}
}

// Connected 连接状态
func (h *Harbor) Connected() bool {
	return h.protocol != nil && h.protocol.connected
}

// Watch 想要观察的其它服务，可以按类型监听，也可以监听所有的服务.传入["all"]监听所有服务
func (h *Harbor) Watch(watchs []string) {
	h.watchs = watchs
}

// Serv 启动本地服务监听
func (h *Harbor) Serv(addr string, port int) error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return err
	}

	h.serviceListener = l
	h.serviceAddr = addr
	h.servicePort = port
	if port == 0 {
		h.servicePort = l.Addr().(*net.TCPAddr).Port
	}

	h.ctx.Core.LogInfo("service listen:", l.Addr())
	return nil
}

// Expose 启动客户端的监听
func (h *Harbor) Expose(outer string, addr string, port int) error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return err
	}

	h.clientListener = l
	h.outerAddr = outer
	h.clientAddr = addr
	h.clientPort = port
	if port == 0 {
		h.clientPort = l.Addr().(*net.TCPAddr).Port
	}
	h.ctx.Core.LogInfo("expose listen:", l.Addr())
	return nil
}

// Close 关闭服务
func (h *Harbor) Close() {
	h.quit = true
	if h.clientListener != nil {
		h.clientListener.Close()
	}

	if h.serviceListener != nil {
		h.serviceListener.Close()
	}
	if h.adminConn != nil {
		h.adminConn.Close()
	}
}
