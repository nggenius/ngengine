package main

import (
	"flag"
	"fmt"

	"github.com/mysll/toolkit"
	"github.com/nggenius/ngengine/core"
	"github.com/nggenius/ngengine/core/rpc"
	"github.com/nggenius/ngengine/core/service"
	"github.com/nggenius/ngengine/protocol"
)

type MathSrv struct {
	service.BaseService
	math *MathModule
}

func (s *MathSrv) Prepare(core service.CoreAPI) error {
	s.CoreAPI = core
	s.math = New()
	return nil
}

func (s *MathSrv) Init(opt *service.CoreOption) error {
	s.AddModule(s.math)
	return nil
}

func (s *MathSrv) Start() error {
	s.BaseService.Start()
	return nil
}

type MathModule struct {
	service.Module
}

func New() *MathModule {
	m := new(MathModule)
	return m
}

func (m *MathModule) Name() string {
	return "Math"
}

func (m *MathModule) Init() bool {
	m.Core.RegisterRemote("Math", new(Math))
	return true
}

// Start 模块启动
func (m *MathModule) Start() {
}

// Shut 模块关闭
func (m *MathModule) Shut() {
}

// OnUpdate 模块Update
func (m *MathModule) OnUpdate(t *service.Time) {
	m.Module.Update(t)
}

// OnMessage 模块消息
func (m *MathModule) OnMessage(id int, args ...interface{}) {
}

type Math struct {
}

func NewMath() *Math {
	s := new(Math)
	return s
}

func (s *Math) RegisterCallback(srv rpc.Servicer) {
	srv.RegisterCallback("Add", s.Add)
	srv.RegisterCallback("Sub", s.Sub)
	srv.RegisterCallback("Mul", s.Mul)
	srv.RegisterCallback("Div", s.Div)
}

func (s *Math) Add(src rpc.Mailbox, dest rpc.Mailbox, msg *protocol.Message) (int32, *protocol.Message) {
	var i, j int
	protocol.ParseArgs(msg, &i, &j)
	return protocol.Reply(protocol.TINY, i+j)
}

func (s *Math) Sub(src rpc.Mailbox, dest rpc.Mailbox, msg *protocol.Message) (int32, *protocol.Message) {
	var i, j int
	protocol.ParseArgs(msg, &i, &j)
	return protocol.Reply(protocol.TINY, i-j)
}

func (s *Math) Mul(src rpc.Mailbox, dest rpc.Mailbox, msg *protocol.Message) (int32, *protocol.Message) {
	var i, j int
	protocol.ParseArgs(msg, &i, &j)
	return protocol.Reply(protocol.TINY, i*j)
}

func (s *Math) Div(src rpc.Mailbox, dest rpc.Mailbox, msg *protocol.Message) (int32, *protocol.Message) {
	var i, j int
	protocol.ParseArgs(msg, &i, &j)
	return protocol.Reply(protocol.TINY, i/j)
}

var args = `{
	"ServId":%d,
	"ServType": "math",
	"AdminAddr":"127.0.0.1",
	"AdminPort":12500,
	"ServName": "math_1",
	"ServAddr": "127.0.0.1",
	"ServPort": 0,
	"Expose": false,
	"OuterAddr":"",
	"HostAddr": "0.0.0.0",
	"HostPort": 0,
	"LogFile":"log/math.log",
	"Args": {}
}`

var (
	id = flag.Int("i", 2, "-i 4")
)

func main() {
	flag.Parse()

	core.RegisterService("math", new(MathSrv))
	core.CreateService("math", fmt.Sprintf(args, *id))
	core.RunAllService()
	toolkit.WaitForQuit()
	core.CloseAllService()
	core.Wait()
}
