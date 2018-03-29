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
	core   service.CoreApi
	mode   int
	client *StoreClient
	store  *Store
}

func New() *StoreModule {
	m := &StoreModule{}
	return m
}

func (m *StoreModule) Name() string {
	return "Store"
}

func (m *StoreModule) Init(core service.CoreApi) bool {
	m.core = core
	switch m.mode {
	case STORE_CLIENT:
		m.core.Service().AddListener(share.EVENT_READY, m.client.OnDatabaseReady)
	case STORE_SERVER:
		core.RegisterRemote("store", m.store)
	default:
		return false
	}

	return true
}

// Shut 模块关闭
func (m *StoreModule) Shut() {
	switch m.mode {
	case STORE_CLIENT:
		m.core.Service().RemoveListener(share.EVENT_READY, m.client.OnDatabaseReady)
	}
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

func (m *StoreModule) Client() *StoreClient {
	return m.client
}
