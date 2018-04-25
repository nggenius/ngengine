package login

import (
	"ngengine/core/rpc"
	"ngengine/game/gameobject/entity/inner"
	"ngengine/module/store"
	"ngengine/protocol"
	"ngengine/protocol/proto/c2s"
	"ngengine/protocol/proto/s2c"
	"ngengine/share"
)

type Account struct {
	ctx *LoginModule
	db  *rpc.Mailbox
}

func (a *Account) RegisterCallback(s rpc.Servicer) {
	s.RegisterCallback("Login", a.Login)
}

func (a *Account) Login(mailbox rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	login := &c2s.Login{}
	if err := a.ctx.core.ParseProto(msg, login); err != nil {
		a.ctx.core.LogErr("login parse error,", err)
		return 0, nil
	}

	result := s2c.Error{}
	result.ErrCode = share.S2C_ERR_SUCCEED

	if a.db == nil {
		result.ErrCode = share.S2C_ERR_SERVICE_INVALID
		a.ctx.core.Mailto(nil, &mailbox, "login.Error", result)
		return 0, nil
	}

	if err := a.ctx.storeClient.Get(
		mailbox.String(),
		"inner.Account",
		map[string]interface{}{
			"Account=?":  login.Name,
			"Password=?": login.Pass,
		},
		a.LoginResult); err != nil {
		result.ErrCode = share.S2C_ERR_SERVICE_INVALID
		a.ctx.core.Mailto(nil, &mailbox, "login.Error", result)
		return 0, nil
	}
	return 0, nil
}

func (a *Account) LoginResult(reply *protocol.Message) {
	result := s2c.Error{}
	result.ErrCode = share.S2C_ERR_SUCCEED
	accinfo := &inner.Account{}
	errcode, err, tag := store.ParseGetReply(reply, accinfo)
	mb, err := rpc.NewMailboxFromStr(tag)
	if err != nil {
		a.ctx.core.LogErr(err)
		return
	}

	if err != nil {
		result.ErrCode = errcode
		a.ctx.core.Mailto(nil, &mb, "login.Error", result)
		return
	}

	if accinfo.Id != 0 {
		if err := a.ctx.storeClient.Find(tag, "inner.Role", map[string]interface{}{"Account=?": accinfo.Account}, 0, 0, a.RoleInfo); err != nil {
			result.ErrCode = share.S2C_ERR_SERVICE_INVALID
			a.ctx.core.Mailto(nil, &mb, "login.Error", result)
			a.ctx.core.LogErr(err)
			return
		}
		return
	}
	result.ErrCode = share.S2C_ERR_NAME_PASS
	a.ctx.core.Mailto(nil, &mb, "login.Error", result)
}

func (a *Account) RoleInfo(reply *protocol.Message) {
	result := s2c.Error{}
	result.ErrCode = share.S2C_ERR_SUCCEED
	var roleinfo []*inner.Role
	errcode, err, tag := store.ParseGetReply(reply, &roleinfo)
	mb, err := rpc.NewMailboxFromStr(tag)
	if err != nil {
		a.ctx.core.LogErr(err)
		return
	}

	if err != nil {
		result.ErrCode = errcode
		a.ctx.core.Mailto(nil, &mb, "login.Error", result)
		return
	}

	if err != nil {
		result.ErrCode = errcode
		a.ctx.core.Mailto(nil, &mb, "login.Error", result)
		return
	}
}

func (a *Account) OnDatabaseReady(evt string, args ...interface{}) {
	srv := a.ctx.core.LookupOneServiceByType("database")
	if srv == nil {
		a.db = nil
		return
	}

	mb := rpc.GetServiceMailbox(srv.Id)
	a.db = &mb
}
