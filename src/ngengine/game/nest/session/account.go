package session

import (
	"ngengine/core/rpc"
	"ngengine/game/gameobject"
	"ngengine/game/gameobject/entity"
	"ngengine/game/gameobject/entity/inner"
	"ngengine/game/store/extension"
	"ngengine/module/store"
	"ngengine/protocol"
	"ngengine/protocol/proto/c2s"
	"ngengine/share"
	"ngengine/utils"
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
	return protocol.Reply(protocol.TINY, id, token)
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
func (a *Account) OnRoleInfo(e *rpc.Error, ar *utils.LoadArchive) {
	var roles []*inner.Role
	err, tag := store.ParseGetReply(e, ar, &roles)
	if err != nil && err.ErrCode == share.ERR_ARGS_ERROR {
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

	var errcode int32
	if err != nil {
		errcode = err.ErrCode
	}

	session.Dispatch(EROLEINFO, [2]interface{}{errcode, roles})
}

func (a *Account) CreateRole(session *Session, args c2s.CreateRole) error {

	player := entity.Create(a.ctx.mainEntity)
	player.SetAttr("Name", args.Name)
	player.SetId(a.ctx.core.GenerateGUID())

	role := inner.Role{}
	role.Account = session.Account
	role.CreateTime = time.Now()
	role.Index = int8(args.Index)
	role.RoleName = args.Name
	role.Id = player.DBId()

	if err := a.ctx.store.Custom(
		session.Mailbox.String(),
		a.OnCreateRole,
		"Store.CreateRole",
		&role,
		player.Archive()); err != nil {
		session.Error(share.S2C_ERR_SERVICE_INVALID)
		return err
	}
	return nil
}

func (a *Account) OnCreateRole(e *rpc.Error, ar *utils.LoadArchive) {
	err, tag := extension.ParseCreateRole(e, ar)
	if err != nil && e.ErrCode == share.ERR_ARGS_ERROR {
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

	var errcode int32
	if err != nil {
		errcode = err.ErrCode
	}

	session.Dispatch(ECREATED, errcode)
}

func (a *Account) ChooseRole(session *Session, args c2s.ChooseRole) error {

	if err := a.ctx.store.Get(
		session.Mailbox.String(),
		a.ctx.mainEntity,
		map[string]interface{}{
			"id=?": args.RoleID,
		},
		a.OnChooseRole); err != nil {
		session.Error(share.S2C_ERR_SERVICE_INVALID)
		return err
	}

	return nil
}

func (a *Account) OnChooseRole(e *rpc.Error, ar *utils.LoadArchive) {
	inst, err := a.ctx.factory.Create(a.ctx.mainEntity)
	if err != nil {
		a.ctx.core.LogFatal("entity create failed")
		return
	}

	gameobject, ok := inst.(gameobject.GameObject)
	if !ok {
		a.ctx.factory.Destroy(inst)
		a.ctx.core.LogFatal("entity is not gameobject")
		return
	}

	player := gameobject.Spirit()
	if player == nil {
		a.ctx.factory.Destroy(inst)
		a.ctx.core.LogFatal("spirit is nil")
		return
	}

	err1, tag := store.ParseGetReply(e, ar, player.Archive())

	if err1 != nil && err1.ErrCode == share.ERR_ARGS_ERROR {
		a.ctx.factory.Destroy(inst)
		a.ctx.core.LogErr(err1)
		return
	}

	mailbox, err2 := rpc.NewMailboxFromStr(tag)
	if err2 != nil {
		a.ctx.factory.Destroy(inst)
		a.ctx.core.LogErr(err2)
		return
	}

	session := a.ctx.FindSession(mailbox.Id())
	if session == nil {
		a.ctx.factory.Destroy(inst)
		a.ctx.core.LogErr("session not found", mailbox.Id())
		return
	}

	var errcode int32
	if err1 != nil {
		errcode = err1.ErrCode
	}

	session.Dispatch(ECHOOSED, [2]interface{}{errcode, gameobject})
}

func (a *Account) DeleteRole(session *Session, args c2s.DeleteRole) error {

	err := a.ctx.store.Custom(
		session.Mailbox.String(),
		a.OnDeleteRole,
		"Store.DeleteRole",
		args.RoleId)

	if err != nil {
		return err
	}

	return nil
}

func (a *Account) OnDeleteRole(e *rpc.Error, ar *utils.LoadArchive) {
	err, tag := extension.ParseDeleteRole(e, ar)
	if err != nil && err.ErrCode == share.ERR_ARGS_ERROR {
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

	var errcode int32
	if err != nil {
		errcode = err.ErrCode
	}
	session.Dispatch(EDELETED, errcode)
}
