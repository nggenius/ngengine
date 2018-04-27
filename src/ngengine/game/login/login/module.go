package login

import (
	"container/list"
	"ngengine/common/event"
	"ngengine/core/rpc"
	"ngengine/core/service"
	"ngengine/module/store"
	"ngengine/module/timer"
	"ngengine/share"
	"time"
)

type LoginModule struct {
	service.Module
	core        service.CoreAPI
	account     *Account
	storeClient *store.StoreClient
	timer       *timer.TimerModule
	lastTime    time.Time
	sessions    map[uint64]*Session
	deleted     *list.List
	db          *rpc.Mailbox
}

func New() *LoginModule {
	l := &LoginModule{}
	l.account = &Account{ctx: l}
	l.sessions = make(map[uint64]*Session)
	l.deleted = list.New()
	return l
}

func (l *LoginModule) Name() string {
	return "Login"
}

func (l *LoginModule) Init(core service.CoreAPI) bool {
	store := core.Module("Store").(*store.StoreModule)
	if store == nil {
		core.LogFatal("need Store module")
		return false
	}
	l.core = core
	l.storeClient = store.Client()
	l.core.Service().AddListener(share.EVENT_READY, l.OnDatabaseReady)
	l.core.Service().AddListener(share.EVENT_USER_CONNECT, l.OnConnected)
	l.core.Service().AddListener(share.EVENT_USER_LOST, l.OnDisconnected)
	l.core.RegisterHandler("Account", l.account)
	l.lastTime = time.Now()
	return true
}

// Shut 模块关闭
func (l *LoginModule) Shut() {
	l.core.Service().RemoveListener(share.EVENT_READY, l.OnDatabaseReady)
	l.core.Service().RemoveListener(share.EVENT_USER_CONNECT, l.OnConnected)
	l.core.Service().RemoveListener(share.EVENT_USER_LOST, l.OnDisconnected)
}

func (l *LoginModule) OnUpdate(t *service.Time) {
	if time.Now().Sub(l.lastTime).Seconds() > 1.0 {
		l.lastTime = time.Now()
		for _, c := range l.sessions {
			if !c.delete {
				c.Dispatch(TIMER, nil)
			}
		}
	}

	for ele := l.deleted.Front(); ele != nil; {
		next := ele.Next()
		delete(l.sessions, ele.Value.(uint64))
		l.core.LogDebug("session delete,", ele.Value.(uint64))
		l.deleted.Remove(ele)
		ele = next
	}
}

func (l *LoginModule) OnConnected(evt string, args ...interface{}) {
	arg := args[0].(event.EventArgs)
	id := arg["id"].(uint64)
	c := NewSession(id, l)
	l.sessions[id] = c
	l.core.LogDebug("new session,", id)
}

func (l *LoginModule) OnDisconnected(evt string, args ...interface{}) {
	arg := args[0].(event.EventArgs)
	id := arg["id"].(uint64)
	if c, ok := l.sessions[id]; ok {
		c.Dispatch(BREAK, nil)
	}
}

func (l *LoginModule) FindSession(id uint64) *Session {
	if c, ok := l.sessions[id]; ok {
		return c
	}
	return nil
}

func (l *LoginModule) OnDatabaseReady(evt string, args ...interface{}) {
	srv := l.core.LookupOneServiceByType("database")
	if srv == nil {
		l.db = nil
		return
	}

	mb := rpc.GetServiceMailbox(srv.Id)
	l.db = &mb
}
