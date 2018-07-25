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

func (m *StoreModule) Sql() *Sql {
	return m.sql
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

// 扩充接口
func (m *StoreModule) Extend(name string, ext Extension) {
	if m.store != nil {
		m.store.AddExtension(name, ext)
		return
	}

	m.Core.LogErr("add extension failed")
}

func (m *StoreModule) Init() bool {
	switch m.mode {
	case STORE_CLIENT:
		m.Core.Service().AddListener(share.EVENT_READY, m.client.OnDatabaseReady)
	case STORE_SERVER:
		m.sql.Init(m.Core)
		m.Core.RegisterRemote("Store", m.store)
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
		m.Core.Service().RemoveListener(share.EVENT_READY, m.client.OnDatabaseReady)
	}
}

func (m *StoreModule) Client() *StoreClient {
	return m.client
}

func (m *StoreModule) Register(name string, creater DataCreater) error {
	return m.register.Register(name, creater)
}
