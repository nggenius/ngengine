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

// 登录模块
// 负责客户端的登录，提供用户名和密码的简单方式
// 登录流程:
// 	Client 提供帐号密码
//	通过数据库进行帐密验证
//  验证成功后，查找一个负载最小的nest。请求nest进行玩家登录
// 	nest返回token
// 	向客户端发送nest的ip,port,token
// 	客户端断开连接，重新建立与nest的连接，通过token进行验证
type LoginModule struct {
	service.Module
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

func (l *LoginModule) Init() bool {
	store := l.Core.Module("Store").(*store.StoreModule)
	if store == nil {
		l.Core.LogFatal("need Store module")
		return false
	}
	l.storeClient = store.Client()
	l.Core.Service().AddListener(share.EVENT_READY, l.OnDatabaseReady)
	l.Core.Service().AddListener(share.EVENT_USER_CONNECT, l.OnConnected)
	l.Core.Service().AddListener(share.EVENT_USER_LOST, l.OnDisconnected)
	l.Core.RegisterHandler("Account", l.account)
	l.lastTime = time.Now()
	return true
}

// Shut 模块关闭
func (l *LoginModule) Shut() {
	l.Core.Service().RemoveListener(share.EVENT_READY, l.OnDatabaseReady)
	l.Core.Service().RemoveListener(share.EVENT_USER_CONNECT, l.OnConnected)
	l.Core.Service().RemoveListener(share.EVENT_USER_LOST, l.OnDisconnected)
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

	// 清理删除对象
	for ele := l.deleted.Front(); ele != nil; {
		next := ele.Next()
		delete(l.sessions, ele.Value.(uint64))
		l.Core.LogDebug("session delete,", ele.Value.(uint64))
		l.deleted.Remove(ele)
		ele = next
	}
}

// 客户端连接回调
func (l *LoginModule) OnConnected(evt string, args ...interface{}) {
	arg := args[0].(event.EventArgs)
	id := arg["id"].(uint64)
	c := NewSession(id, l)
	l.sessions[id] = c
	l.Core.LogDebug("new session,", id)
}

// 客户端断线回调
func (l *LoginModule) OnDisconnected(evt string, args ...interface{}) {
	arg := args[0].(event.EventArgs)
	id := arg["id"].(uint64)
	if c, ok := l.sessions[id]; ok {
		c.Dispatch(BREAK, nil)
	}
}

// 查找连接的客户端信息
func (l *LoginModule) FindSession(id uint64) *Session {
	if c, ok := l.sessions[id]; ok {
		return c
	}
	return nil
}

// 服务变动回调
func (l *LoginModule) OnDatabaseReady(evt string, args ...interface{}) {
	srv := l.Core.LookupOneServiceByType("store")
	if srv == nil {
		l.db = nil
		return
	}

	mb := rpc.GetServiceMailbox(srv.Id)
	l.db = &mb
}
