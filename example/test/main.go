package main

import (
	"time"

	"github.com/mysll/toolkit"
	"github.com/nggenius/ngengine/core"
	"github.com/nggenius/ngengine/core/rpc"
	"github.com/nggenius/ngengine/core/service"
	"github.com/nggenius/ngengine/share"
	"github.com/nggenius/ngengine/utils"
)

type Client struct {
	service.BaseService
	m *ClientModule
}

func (s *Client) Prepare(core service.CoreAPI) error {
	s.CoreAPI = core
	s.m = New()
	return nil
}

func (s *Client) Init(opt *service.CoreOption) error {
	s.AddModule(s.m)
	return nil
}

func (s *Client) Start() error {
	s.BaseService.Start()
	return nil
}

type ClientModule struct {
	service.Module
}

func New() *ClientModule {
	m := new(ClientModule)
	return m
}

func (m *ClientModule) Name() string {
	return "Client"
}

func (m *ClientModule) Init() bool {
	m.Core.Service().AddListener(share.EVENT_MUST_SERVICE_READY, m.OnReady)
	return true
}

func (m *ClientModule) OnReady(event string, args ...interface{}) {
	m.AddPeriod(time.Second)
	m.AddCallback(time.Second, m.Echo)
	m.Core.LogInfo("OnReady")
}

func (m *ClientModule) Echo(t time.Duration) {
	srv := m.Core.LookupRandServiceByType("echo")
	if srv != nil {
		m.Core.MailtoAndCallback(nil, srv.Mailbox(), "Echo.Echo", m.OnEcho, nil, "hello world")
	}

	srv = m.Core.LookupRandServiceByType("math")
	if srv != nil {
		m.Core.MailtoAndCallback(nil, srv.Mailbox(), "Math.Add", m.OnResult, nil, toolkit.RandRange(1, 10000), toolkit.RandRange(1, 10000))
	}
}

func (m *ClientModule) OnEcho(param interface{}, replyerr *rpc.Error, ar *utils.LoadArchive) {
	if replyerr != nil {
		m.Core.LogErr(replyerr)
		return
	}

	result, _ := ar.ReadString()

	m.Core.LogInfo("echo:", result)
}

func (m *ClientModule) OnResult(param interface{}, replyerr *rpc.Error, ar *utils.LoadArchive) {
	if replyerr != nil {
		m.Core.LogErr(replyerr)
		return
	}

	result, _ := ar.ReadInt64()

	m.Core.LogInfo("result:", result)
}

var start = `{
	"ServId":3,
	"ServType": "client",
	"AdminAddr":"127.0.0.1",
	"AdminPort":12500,
	"ServName": "client_1",
	"ServAddr": "127.0.0.1",
	"ServPort": 0,
	"Expose": false,
	"OuterAddr":"",
	"HostAddr": "0.0.0.0",
	"HostPort": 0,
	"LogFile":"log/client.log",
	"Args": {}
}`

func main() {
	core.RegisterService("client", new(Client))
	core.CreateService("client", start)
	core.RunAllService()
	toolkit.WaitForQuit()
	core.CloseAllService()
	core.Wait()
}
