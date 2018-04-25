package store

import (
	"ngengine/core/service"
	"ngengine/share"
)

const (
	STORE_CLIENT = iota + 1
	STORE_SERVER
)

type StoreModule struct {
	service.Module
	core     service.CoreAPI
	mode     int
	client   *StoreClient
	store    *Store
	register *Register
	sql      *Sql
}

func New() *StoreModule {
	m := &StoreModule{}
	m.register = newRegister()
	m.sql = newSql()
	return m
}

func (m *StoreModule) Name() string {
	return "Store"
}

// SetMode 设置工作模式
func (m *StoreModule) SetMode(mode int) {
	switch mode {
	case STORE_CLIENT:
		m.client = NewStoreClient(m)
	case STORE_SERVER:
		m.store = NewStore(m)
	default:
		panic("mode is illegal")
	}

	m.mode = mode
}

func (m *StoreModule) Init(core service.CoreAPI) bool {
	m.core = core

	switch m.mode {
	case STORE_CLIENT:
		m.core.Service().AddListener(share.EVENT_READY, m.client.OnDatabaseReady)
	case STORE_SERVER:
		m.sql.Init(core)
		core.RegisterRemote("Store", m.store)
	default:
		return false
	}

	return true
}

func (m *StoreModule) Start() {
	if m.mode == STORE_SERVER {
		err := m.register.Sync(m)
		if err != nil {
			panic(err)
		}
	}
}

// Shut 模块关闭
func (m *StoreModule) Shut() {
	switch m.mode {
	case STORE_CLIENT:
		m.core.Service().RemoveListener(share.EVENT_READY, m.client.OnDatabaseReady)
	}
}

func (m *StoreModule) Client() *StoreClient {
	return m.client
}

func (m *StoreModule) Register(name string, creater DataCreater) error {
	return m.register.Register(name, creater)
}
