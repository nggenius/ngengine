package session

import (
	"ngengine/core/rpc"
	"ngengine/game/gameobject/entity/inner"
	"ngengine/module/store"
	"ngengine/protocol"
	"ngengine/share"
)

type Account struct {
	ctx *SessionModule
}

func NewAccount(ctx *SessionModule) *Account {
	a := &Account{}
	a.ctx = ctx
	return a
}

func (a *Account) RegisterCallback(s rpc.Servicer) {
	s.RegisterCallback("Logged", a.Logged)
}

// login服务调用
func (a *Account) Logged(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	var id uint64
	var account string
	m.Read(&id)
	m.Read(&account)
	token := a.ctx.cache.Put(account)
	return 0, protocol.ReplyMessage(protocol.TINY, id, token)
}

// 请求玩家信息
func (a *Account) requestRoleInfo(session *Session) error {
	if err := a.ctx.requestRoleInfo.Find(
		session.Mailbox.String(),
		"inner.Role",
		map[string]interface{}{
			"Account=?": session.Account,
		},
		0, 0, a.OnRoleInfo); err != nil {
		session.Error(share.S2C_ERR_SERVICE_INVALID)
		return err
	}
	return nil
}

// 收到玩家信息
func (a *Account) OnRoleInfo(msg *protocol.Message) {
	var roles []*inner.Role
	errcode, err, tag := store.ParseGetReply(msg, &roles)
	if err != nil {
		a.ctx.core.LogErr(err)
		return
	}
	mailbox, err1 := rpc.NewMailboxFromStr(tag)
	if err1 != nil {
		a.ctx.core.LogErr(err1)
		return
	}

	session := a.ctx.FindSession(mailbox.Id())
	if session == nil {
		a.ctx.core.LogErr("session not found", mailbox.Id())
		return
	}

	session.Dispatch(ROLE_INFO, [2]interface{}{errcode, roles})
}
