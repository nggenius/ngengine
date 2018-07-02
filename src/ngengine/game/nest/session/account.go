package session

import (
	"ngengine/core/rpc"
	"ngengine/game/gameobject/entity"
	"ngengine/game/gameobject/entity/inner"
	"ngengine/module/store"
	"ngengine/protocol"
	"ngengine/protocol/proto/c2s"
	"ngengine/share"
	"time"
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
	if err := a.ctx.store.Find(
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

func (a *Account) CreateRole(session *Session, args c2s.CreateRole) error {

	player := entity.NewPlayer()
	player.SetName(args.Name)
	player.Archive().SetId(a.ctx.core.GenerateGUID())

	role := inner.Role{}
	role.Account = session.Account
	role.CreateTime = time.Now()
	role.Index = int8(args.Index)
	role.RoleName = args.Name
	role.RoleId = player.Archive().DBId()

	if err := a.ctx.store.MultiInsert(
		session.Mailbox.String(),
		a.OnCreateRole,
		[]string{
			"entity.Player",
			"inner.Role",
		},
		[]interface{}{
			player.Archive(),
			&role,
		},
	); err != nil {
		session.Error(share.S2C_ERR_SERVICE_INVALID)
		return err
	}
	return nil
}

func (a *Account) OnCreateRole(msg *protocol.Message) {
	errcode, err, tag := store.ParseMultiInsertReply(msg)
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

	session.Dispatch(CREATED, errcode)
}

func (a *Account) ChooseRole(session *Session, args c2s.ChooseRole) error {

	if err := a.ctx.store.Get(
		session.Mailbox.String(),
		"entity.Player",
		map[string]interface{}{
			"id=?": args.RoleID,
		},
		a.OnChooseRole); err != nil {
		session.Error(share.S2C_ERR_SERVICE_INVALID)
		return err
	}

	return nil
}

func (a *Account) OnChooseRole(msg *protocol.Message) {
	player := entity.NewPlayer()
	errcode, err, tag := store.ParseGetReply(msg, player.Archive())
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

	session.Dispatch(CHOOSED, [2]interface{}{errcode, player})
}
