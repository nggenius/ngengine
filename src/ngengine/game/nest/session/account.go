package session

import (
	"fmt"
	"ngengine/core/rpc"
	"ngengine/game/gameobject"
	"ngengine/game/gameobject/entity"
	"ngengine/game/gameobject/entity/inner"
	"ngengine/module/store"
	"ngengine/protocol"
	"ngengine/protocol/proto/c2s"
	"ngengine/share"
	"ngengine/utils"
	"time"
)

type LandInfo interface {
	SetLandScene(landscene int64)
	SetLandPosXYZOrient(x float64, y float64, z float64, orient float64)
	LandScene() int64
	LandPosXYZOrient() (x float64, y float64, z float64, orient float64)
}

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
	var account string
	m.Read(&account)
	token := a.ctx.cache.Put(account)
	return protocol.Reply(protocol.TINY, token)
}

// 请求玩家信息
func (a *Account) requestRoleInfo(session *Session) error {
	if err := a.ctx.store.Find(
		session.Mailbox,
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
func (a *Account) OnRoleInfo(p interface{}, e *rpc.Error, ar *utils.LoadArchive) {
	mailbox := p.(*rpc.Mailbox)

	session := a.ctx.FindSession(mailbox.Id())
	if session == nil {
		a.ctx.Core.LogErr("session not found", mailbox.Id())
		return
	}

	var roles []*inner.Role
	err := store.ParseGetReply(e, ar, &roles)

	if err != nil {
		session.Dispatch(EROLEINFO, [2]interface{}{err.ErrCode(), roles})
		return
	}
	session.Dispatch(EROLEINFO, [2]interface{}{rpc.OK, roles})
}

func (a *Account) CreateRole(session *Session, args c2s.CreateRole) error {

	player := entity.Create(a.ctx.mainEntity)
	player.SetAttr("Name", args.Name)
	player.SetId(a.ctx.Core.GenerateGUID())
	li, ok := player.(LandInfo)
	if !ok {
		return fmt.Errorf("player not implement LandInfo")
	}

	li.SetLandScene(1)
	li.SetLandPosXYZOrient(0, 0, 0, 0)

	role := inner.Role{}
	role.Account = session.Account
	role.CreateTime = time.Now()
	role.Index = int8(args.Index)
	role.RoleName = args.Name
	role.Id = player.DBId()

	if err := a.ctx.store.Custom(
		session.Mailbox,
		a.OnCreateRole,
		"Store.CreateRole",
		&role,
		player.Archive()); err != nil {
		session.Error(share.S2C_ERR_SERVICE_INVALID)
		return err
	}
	return nil
}

func (a *Account) OnCreateRole(p interface{}, e *rpc.Error, ar *utils.LoadArchive) {

	mailbox := p.(*rpc.Mailbox)

	session := a.ctx.FindSession(mailbox.Id())
	if session == nil {
		a.ctx.Core.LogErr("session not found", mailbox.Id())
		return
	}

	if e != nil {
		session.Dispatch(ECREATED, e.ErrCode())
		return
	}

	session.Dispatch(ECREATED, rpc.OK)

}

func (a *Account) ChooseRole(session *Session, args c2s.ChooseRole) error {

	if err := a.ctx.store.Get(
		session.Mailbox,
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

func (a *Account) OnChooseRole(p interface{}, e *rpc.Error, ar *utils.LoadArchive) {

	mailbox := p.(*rpc.Mailbox)

	session := a.ctx.FindSession(mailbox.Id())
	if session == nil {
		a.ctx.Core.LogErr("session not found", mailbox.Id())
		return
	}

	inst, err := a.ctx.factory.Create(a.ctx.role)
	if err != nil {
		a.ctx.Core.LogFatal("entity create failed")
		return
	}

	gameobject, ok := inst.(gameobject.GameObject)
	if !ok {
		a.ctx.factory.Destroy(inst)
		a.ctx.Core.LogFatal("entity is not gameobject")
		return
	}

	gameobject.SetCap(0)

	player := gameobject.Spirit()
	if player == nil {
		a.ctx.factory.Destroy(inst)
		a.ctx.Core.LogFatal("spirit is nil")
		return
	}

	err1 := store.ParseGetReply(e, ar, player.Archive())

	if err1 != nil {
		a.ctx.factory.Destroy(inst)
		a.ctx.Core.LogErr(err1)
		session.Dispatch(ECHOOSED, [2]interface{}{err1.ErrCode(), nil})
		return
	}

	session.Dispatch(ECHOOSED, [2]interface{}{rpc.OK, gameobject})
}

func (a *Account) DeleteRole(session *Session, args c2s.DeleteRole) error {

	err := a.ctx.store.Custom(
		session.Mailbox,
		a.OnDeleteRole,
		"Store.DeleteRole",
		args.RoleId)

	if err != nil {
		return err
	}

	return nil
}

func (a *Account) OnDeleteRole(p interface{}, e *rpc.Error, ar *utils.LoadArchive) {
	mailbox := p.(*rpc.Mailbox)

	session := a.ctx.FindSession(mailbox.Id())
	if session == nil {
		a.ctx.Core.LogErr("session not found", mailbox.Id())
		return
	}

	if e != nil {
		session.Dispatch(EDELETED, e.ErrCode())
		return
	}

	session.Dispatch(EDELETED, rpc.OK)

}

func (a *Account) FindRegion(session *Session, id int64, fx, fy, fz float64) error {
	srv := a.ctx.Core.LookupRandServiceByType("world")
	if srv == nil {
		return fmt.Errorf("world not found")
	}

	return a.ctx.Core.MailtoAndCallback(nil, srv.Mailbox(), "Space.FindRegion", a.OnFindRegion, session.Mailbox, id, fx, fy, fz)
}

func (a *Account) OnFindRegion(param interface{}, replyerr *rpc.Error, ar *utils.LoadArchive) {
	mb := param.(*rpc.Mailbox)
	session := a.ctx.FindSession(mb.Id())
	if session == nil {
		a.ctx.Core.LogErr("session not found,", mb)
		return
	}

	if replyerr != nil && protocol.CheckRpcError(replyerr) {
		session.Dispatch(EFREGION, [2]interface{}{replyerr.Code, rpc.NullMailbox})
		return
	}

	var w rpc.Mailbox
	if err := ar.Read(&w); err != nil {
		panic(err)
	}

	session.Dispatch(EFREGION, [2]interface{}{rpc.OK, w})

}

func (p *Account) EnterRegion(s *Session, r rpc.Mailbox) error {
	data, err := p.ctx.factory.Encode(s.gameobject.Spirit().ObjId())
	if err != nil {
		return err
	}
	return p.ctx.Core.MailtoAndCallback(nil, &r, "GameScene.AddPlayer", p.OnEnterRegion, s.id, data)
}

func (p *Account) OnEnterRegion(param interface{}, replyerr *rpc.Error, ar *utils.LoadArchive) {
	session := p.ctx.FindSession(param.(uint64))
	if session == nil {
		p.ctx.Core.LogErr("session not found,", param)
		return
	}

	if replyerr != nil {
		p.ctx.Core.LogErr("enter region failed ", replyerr)
		return
	}
	p.ctx.Core.LogInfo("enter session ", param, replyerr)
	session.Dispatch(EONLINE, nil)
}
