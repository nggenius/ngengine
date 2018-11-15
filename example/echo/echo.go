package main

import (
	"github.com/mysll/toolkit"
	"github.com/nggenius/ngengine/core"
	"github.com/nggenius/ngengine/core/rpc"
	"github.com/nggenius/ngengine/core/service"
	"github.com/nggenius/ngengine/protocol"
)

type EchoSrv struct {
	service.BaseService
	echo *EchoModule
}

func (s *EchoSrv) Prepare(core service.CoreAPI) error {
	s.CoreAPI = core
	s.echo = New()
	return nil
}

func (s *EchoSrv) Init(opt *service.CoreOption) error {
	s.AddModule(s.echo)
	return nil
}

func (s *EchoSrv) Start() error {
	s.BaseService.Start()
	return nil
}

type EchoModule struct {
	service.Module
}

func New() *EchoModule {
	m := new(EchoModule)
	return m
}

func (m *EchoModule) Name() string {
	return "Echo"
}

func (m *EchoModule) Init() bool {
	m.Core.RegisterRemote("Echo", new(Echo))
	return true
}

// Start 模块启动
func (m *EchoModule) Start() {
}

// Shut 模块关闭
func (m *EchoModule) Shut() {
}

// OnUpdate 模块Update
func (m *EchoModule) OnUpdate(t *service.Time) {
	m.Module.Update(t)
}

// OnMessage 模块消息
func (m *EchoModule) OnMessage(id int, args ...interface{}) {
}

type Echo struct {
}

func NewEcho() *Echo {
	s := new(Echo)
	return s
}

func (s *Echo) RegisterCallback(srv rpc.Servicer) {
	srv.RegisterCallback("Echo", s.Echo)
}

func (s *Echo) Echo(src rpc.Mailbox, dest rpc.Mailbox, msg *protocol.Message) (int32, *protocol.Message) {
	var info string
	protocol.ParseArgs(msg, &info)
	return protocol.Reply(protocol.TINY, info)
}

var startnest = `{
	"ServId":1,
	"ServType": "echo",
	"AdminAddr":"127.0.0.1",
	"AdminPort":12500,
	"ServName": "echo_1",
	"ServAddr": "127.0.0.1",
	"ServPort": 0,
	"Expose": false,
	"OuterAddr":"",
	"HostAddr": "0.0.0.0",
	"HostPort": 0,
	"LogFile":"log/echo.log",
	"Args": {}
}`

func main() {
	core.RegisterService("echo", new(EchoSrv))
	core.CreateService("echo", startnest)
	core.RunAllService()
	toolkit.WaitForQuit()
	core.CloseAllService()
	core.Wait()
}
