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
	store      *store.StoreClient
	factory    *object.ObjectModule
	account    *Account
	proxy      *proxy
	sessions   SessionDB  // session管理器
	deleted    *list.List // 标志为删除的session
	cache      cache      // 缓存的口令
	mainEntity string     // 主实体
	role       string     // 玩家类名
	ls         map[string]*event.EventListener
}

func New() *SessionModule {
	l := &SessionModule{}
	l.account = NewAccount(l)
	l.proxy = NewProxy(l)
	l.cache = make(cache)
	l.sessions = make(SessionDB)
	l.deleted = list.New()
	l.ls = make(map[string]*event.EventListener)
	return l
}

func (s *SessionModule) Name() string {
	return "Session"
}

func (s *SessionModule) Init() bool {
	opt := s.Core.Option()
	s.mainEntity = opt.Args.String("MainEntity")
	s.role = opt.Args.String("Role")
	store := s.Core.MustModule("Store").(*store.StoreModule)
	if store == nil {
		s.Core.LogFatal("need Store module")
		return false
	}
	factory := s.Core.MustModule("Object").(*object.ObjectModule)
	if factory == nil {
		s.Core.LogFatal("need object module")
		return false
	}
	s.factory = factory
	s.store = store.Client()
	s.Core.RegisterRemote("Account", s.account)
	s.Core.RegisterHandler("Self", s.proxy)
	s.ls[share.EVENT_USER_CONNECT] = s.Core.Service().AddListener(share.EVENT_USER_CONNECT, s.OnConnected)
	s.ls[share.EVENT_USER_LOST] = s.Core.Service().AddListener(share.EVENT_USER_LOST, s.OnDisconnected)
	s.AddPeriod(time.Second)
	s.AddCallback(time.Second, s.PerSecondCheck)
	return true
}

// Shut 模块关闭
func (s *SessionModule) Shut() {
	for k, v := range s.ls {
		s.Core.Service().RemoveListener(k, v)
	}
}

// PerSecondCheck 每分钟检查
func (s *SessionModule) PerSecondCheck(d time.Duration) {
	s.cache.Check()
	for _, s := range s.sessions {
		if !s.delete {
			s.Dispatch(ETIMER, nil)
		}
	}
}

func (s *SessionModule) OnUpdate(t *service.Time) {
	s.Module.OnUpdate(t)
	// 清理删除对象
	for ele := s.deleted.Front(); ele != nil; {
		next := ele.Next()
		delete(s.sessions, ele.Value.(uint64))
		s.Core.LogInfo("remove session, ", ele.Value.(uint64))
		s.Core.UpdateLoad(s.Core.Load() - 1)
		s.deleted.Remove(ele)
		ele = next
	}
}

// 新客户端连接
func (s *SessionModule) OnConnected(evt string, args ...interface{}) {
	arg := args[0].(event.EventArgs)
	id := arg["id"].(uint64)
	ns := NewSession(id, s)
	s.Core.LogInfo("new session, ", id)
	s.sessions[id] = ns
	s.Core.UpdateLoad(s.Core.Load() + 1)
}

// 客户端断线
func (s *SessionModule) OnDisconnected(evt string, args ...interface{}) {
	arg := args[0].(event.EventArgs)
	id := arg["id"].(uint64)
	if s, ok := s.sessions[id]; ok {
		s.Dispatch(EBREAK, nil)
	}
}

// 查找session
func (s *SessionModule) FindSession(id uint64) *Session {
	if session, ok := s.sessions[id]; ok {
		return session
	}
	return nil
}
