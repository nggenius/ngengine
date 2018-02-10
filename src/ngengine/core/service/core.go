package service

import (
	"fmt"
	"ngengine/common/event"
	"ngengine/core/rpc"
	"ngengine/logger"
	"ngengine/protocol"
	"ngengine/protocol/proto"
	"ngengine/share"
	"time"

	"github.com/mysll/toolkit"
	"github.com/petermattis/goid"
)

const (
	CS_NONE = iota
	CS_CLOSE
	CS_QUIT
	CS_SHUT
	WATCH_INTERVAL = time.Second * 10
)

// 服务的核心
type Core struct {
	*logger.Log
	toolkit.WaitGroupWrapper
	Id         share.ServiceId     // 服务ID
	gid        int64               // goroutine id
	opts       *CoreOption         // 配置项
	service    Service             // 逻辑服务
	closeState int                 // 关闭状态
	quitCh     chan struct{}       // 关闭chan
	startTime  time.Time           // 启动时间
	harbor     *Harbor             // 服务连接点
	watchs     []string            // 想要关注的服务
	time       *Time               // 服务器时间
	dns        *SrvDNS             // dns服务
	rpcCh      chan *rpc.RpcCall   // rpc 调用通道
	rpcSvr     *rpc.Server         // rpc 服务
	Emitter    *event.AsyncEvent   // 事件调度器
	Mailbox    rpc.Mailbox         // 服务的地址
	busy       bool                // 运行状态
	Proto      protocol.ProtoCodec // 消息编码解码器
	clientDB   *ClientDB           // 客户端管理
	rpcmgr     *RpcRegister        // rpc注册
	s2chelper  *S2CHelper          // 客户端调用工具
	modules    *modules
}

// 创建一个服务
func CreateService(s Service) *Core {
	sc := &Core{
		service:    s,
		closeState: CS_NONE,
		quitCh:     make(chan struct{}),
		watchs:     make([]string, 0, 8),
		Emitter:    event.NewAsyncEvent(),
		rpcmgr:     NewRpcRegister(),
		Proto:      &proto.JsonProto{},
		modules:    NewModules(),
	}

	sc.s2chelper = NewS2CHelper(sc)

	sc.rpcmgr.RegisterHandler("C2SHelper", &C2SHelper{sc})
	sc.rpcmgr.RegisterRemote("S2CHelper", sc.s2chelper)

	if err := sc.service.Prepare(sc); err != nil {
		panic(err)
	}

	return sc
}

// 初始化服务
func (c *Core) Init(args string) error {
	c.startTime = time.Now()
	opt, err := ParseOption(args)
	if err != nil {
		panic(err)
	}

	c.Id = opt.ServId
	c.Mailbox = rpc.GetServiceMailbox(opt.ServId)
	c.opts = opt
	if c.opts.LogFile == "" {
		c.opts.LogFile = fmt.Sprintf("%s_%d.log", c.startTime.Format("06_01_02_15_04_05"), toolkit.RandRange(100, 999))
	}
	c.Log = logger.New(c.opts.LogFile, c.opts.LogLevel)

	if c.opts.AdminAddr == "" || c.opts.AdminPort == 0 {
		c.LogFatalf("admin address is error, get (%s:%d)", c.opts.AdminAddr, c.opts.AdminPort)
	}

	// 调用服务的初始化
	if err := c.service.Init(c.opts); err != nil {
		c.LogFatal(err)
		return err
	}

	// 初始化模块
	em := make([]string, 0, 8)
	for n, m := range c.modules.modules {
		if !m.Init(c) {
			c.LogErr("module '", n, "' init failed")
			em = append(em, n)
			continue
		}

		c.LogErr("module '", n, "' init ok")
	}

	// 删除初始化失败的模块
	for _, n := range em {
		c.LogErr("module '", n, "' init failed, now is removed")
		delete(c.modules.modules, n)
	}

	return nil
}

// 启动服务
func (c *Core) Serv() {
	ctx := &context{c}
	c.dns = NewSrvDNS(ctx)
	if err := c.service.Start(); err != nil {
		c.LogFatal(err)
	}

	c.gid = goid.Get()
	c.LogInfo(c.opts.ServName, " start, goroutine id ", c.gid)
	harbor := NewHarbor(ctx)
	// 连接admin
	harbor.SetAdmin(c.opts.AdminAddr, c.opts.AdminPort)
	if err := harbor.Serv(c.opts.ServAddr, c.opts.ServPort); err != nil {
		c.LogFatal(err)
	}

	if c.watchs == nil || len(c.watchs) == 0 {
		c.watchs = []string{"all"}
	}

	harbor.Watch(c.watchs)
	// 创建rpc服务
	c.rpcCh = make(chan *rpc.RpcCall, c.opts.MaxRpcCall)
	c.rpcSvr = c.rpcmgr.createRpc(c.rpcCh, ctx)
	c.Wrap(func() { rpc.CreateService(c.rpcSvr, harbor.serviceListener, c.Log) })

	// 启动外部连接
	if c.opts.Expose {
		if err := harbor.Expose(c.opts.HostAddr, c.opts.HostPort); err != nil {
			c.LogFatal(err)
		}

		c.Wrap(func() {
			protocol.TCPServer(harbor.clientListener, &ClientHandler{ctx}, c.Log)
		})

		c.clientDB = NewClientDB(ctx)
	}

	c.harbor = harbor
	c.Wrap(func() { harbor.KeepConnect() })
	c.service.Ready()
	c.run()
}

// 停止服务
func (c *Core) Close() {
	if c.closeState != 0 {
		return
	}

	c.closeState = CS_CLOSE
	if !c.service.Close() { // 服务自己控制关闭时机
		return
	}

	c.Shut()
}
