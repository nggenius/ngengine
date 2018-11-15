package ngadmin

import (
	"fmt"
	"net"
	"github.com/nggenius/ngengine/logger"
	"github.com/nggenius/ngengine/protocol"

	"github.com/mysll/toolkit"
)

type NGAdmin struct {
	*logger.Log
	opts        *Options
	tcpListener net.Listener
	waitGroup   toolkit.WaitGroupWrapper
	DB          *ServiceDB
	slave       *Slave
	quit        bool
}

// New 创建一个新的admin,提供配置属性
func New(opts *Options) *NGAdmin {
	if opts.LogFile == "" {
		opts.LogFile = "admin.log"
	}

	if opts.Master {
		opts.Host = opts.LocalAddr
	}

	admin := &NGAdmin{
		Log:  logger.New(opts.LogFile, opts.LogLevel),
		opts: opts,
	}

	return admin
}

// Main 主函数
func (n *NGAdmin) Main() {
	ctx := &Context{n}
	n.DB = NewServiceDB(ctx)

	if n.opts.Master {
		tcpListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", n.opts.LocalAddr, n.opts.Port))
		if err != nil {
			n.LogFatalf("FATAL: listen (%s:%d) failed - %s", n.opts.LocalAddr, n.opts.Port, err)
		}
		n.tcpListener = tcpListener
		tcpServer := &tcpServer{ctx: ctx}
		n.waitGroup.Wrap(func() {
			protocol.TCPServer(tcpListener, tcpServer, n.Log)
		})

		// 启动其他服务器
		if !n.opts.DebugMode {
			n.StartService()
		}
	} else {
		n.slave = NewSlave(ctx)
		n.waitGroup.Wrap(func() {
			n.slave.KeepConnect(n.opts.Host, n.opts.Port)
		})
	}

	n.LogInfo("admin is ready, Debug Mode:", n.opts.DebugMode)
}

// Exit 退出函数
func (n *NGAdmin) Exit() {

	n.shutdown()
	// 等待所有服务退出
	<-n.DB.Done()

	if n.tcpListener != nil {
		n.tcpListener.Close()
	}

	if n.slave != nil {
		n.slave.Close()
	}

	n.waitGroup.Wait()
	if n.Log != nil {
		n.CloseLog()
	}
}

// shutdown 正常退出
func (n *NGAdmin) shutdown() {
	n.quit = true
	// 关闭所有服务
	n.DB.CloseAll()
}
