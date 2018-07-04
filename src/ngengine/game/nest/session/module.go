package session

import (
	"container/list"
	"ngengine/common/event"
	"ngengine/core/service"
	"ngengine/module/object"
	"ngengine/module/store"
	"ngengine/share"
	"time"
)

// 登录Session模块
// 登录信息存储在这里
// 模块提供功能：
// 		proxy:网关功能，对应的客户端消息在这里进行中转处理
//		session:登录管理，角色管理
//		存储客户端对应的entity数据
type SessionModule struct {
	service.Module
	core       service.CoreAPI
	store      *store.StoreClient
	factory    *object.ObjectModule
	account    *Account
	proxy      *proxy
	sessions   SessionDB  // session管理器
	deleted    *list.List // 标志为删除的session
	lastTime   time.Time  // 最后一次更新时间
	cache      cache      // 缓存的口令
	mainEntity string     // 主实体
}

func New() *SessionModule {
	l := &SessionModule{}
	l.account = NewAccount(l)
	l.proxy = NewProxy(l)
	l.cache = make(cache)
	l.sessions = make(SessionDB)
	l.deleted = list.New()
	return l
}

func (s *SessionModule) Name() string {
	return "Session"
}

func (s *SessionModule) Init(core service.CoreAPI) bool {

	opt := core.Option()
	s.mainEntity = opt.Args.String("MainEntity")

	store := core.Module("Store").(*store.StoreModule)
	if store == nil {
		core.LogFatal("need Store module")
		return false
	}
	s.core = core
	factory := core.Module("Object").(*object.ObjectModule)
	if factory == nil {
		core.LogFatal("need object module")
		return false
	}
	s.factory = factory
	s.store = store.Client()
	s.core.RegisterRemote("Account", s.account)
	s.core.RegisterHandler("Self", s.proxy)
	s.core.Service().AddListener(share.EVENT_USER_CONNECT, s.OnConnected)
	s.core.Service().AddListener(share.EVENT_USER_LOST, s.OnDisconnected)
	s.lastTime = time.Now()
	return true
}

// Shut 模块关闭
func (s *SessionModule) Shut() {
	s.core.Service().RemoveListener(share.EVENT_USER_CONNECT, s.OnConnected)
	s.core.Service().RemoveListener(share.EVENT_USER_LOST, s.OnDisconnected)
}

func (s *SessionModule) OnUpdate(t *service.Time) {
	if time.Now().Sub(s.lastTime).Seconds() > 1.0 {
		s.lastTime = time.Now()
		s.cache.Check()
		for _, s := range s.sessions {
			if !s.delete {
				s.Dispatch(TIMER, nil)
			}
		}
	}

	// 清理删除对象
	for ele := s.deleted.Front(); ele != nil; {
		next := ele.Next()
		delete(s.sessions, ele.Value.(uint64))
		s.core.LogInfo("remove session, ", ele.Value.(uint64))
		s.core.UpdateLoad(s.core.Load() - 1)
		s.deleted.Remove(ele)
		ele = next
	}
}

// 新客户端连接
func (s *SessionModule) OnConnected(evt string, args ...interface{}) {
	arg := args[0].(event.EventArgs)
	id := arg["id"].(uint64)
	ns := NewSession(id, s)
	s.core.LogInfo("new session, ", id)
	s.sessions[id] = ns
	s.core.UpdateLoad(s.core.Load() + 1)
}

// 客户端断线
func (s *SessionModule) OnDisconnected(evt string, args ...interface{}) {
	arg := args[0].(event.EventArgs)
	id := arg["id"].(uint64)
	if s, ok := s.sessions[id]; ok {
		s.Dispatch(BREAK, nil)
	}
}

// 查找session
func (s *SessionModule) FindSession(id uint64) *Session {
	if session, ok := s.sessions[id]; ok {
		return session
	}
	return nil
}
