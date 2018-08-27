package login

import (
	"fmt"
	"ngengine/core/rpc"
	"ngengine/core/service"
	"ngengine/game/gameobject/entity/inner"
	"ngengine/module/store"
	"ngengine/protocol"
	"ngengine/protocol/proto/c2s"
	"ngengine/share"
	"ngengine/utils"
)

type Account struct {
	ctx *LoginModule
}

func (a *Account) RegisterCallback(s rpc.Servicer) {
	s.RegisterCallback("Login", a.Login)
}

// client:登录
func (a *Account) Login(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	login := &c2s.Login{}
	if err := a.ctx.Core.ParseProto(msg, login); err != nil {
		a.ctx.Core.LogErr("login parse error,", err)
		return 0, nil
	}

	session := a.ctx.FindSession(sender.Id())
	if session == nil {
		a.ctx.Core.LogErr("session not found", sender.Id())
		return 0, nil
	}

	session.SetMailbox(&sender)
	session.SetAccount(login.Name)
	session.Dispatch(LOGIN, login)
	return 0, nil
}

// 请求数据库帐号信息
func (a *Account) sendLogin(s *Session, login *c2s.Login) error {
	if a.ctx.db == nil {
		s.Error(share.S2C_ERR_SERVICE_INVALID)
		return fmt.Errorf("database is not ready")
	}

	if err := a.ctx.storeClient.Get(
		s.Mailbox,
		"inner.Account",
		map[string]interface{}{
			"Account=?":  login.Name,
			"Password=?": login.Pass,
		},
		a.OnLoginResult); err != nil {
		s.Error(share.S2C_ERR_SERVICE_INVALID)
		return err
	}

	return nil
}

// 帐号信息回调
func (a *Account) OnLoginResult(p interface{}, e *rpc.Error, ar *utils.LoadArchive) {
	mailbox := p.(*rpc.Mailbox)

	session := a.ctx.FindSession(mailbox.Id())
	if session == nil {
		a.ctx.Core.LogErr("session not found", mailbox.Id())
		return
	}

	accinfo := &inner.Account{}
	err := store.ParseGetReply(e, ar, accinfo)
	if err != nil {
		session.Dispatch(LOGIN_RESULT, [2]interface{}{err.ErrCode(), accinfo})
		return
	}

	session.Dispatch(LOGIN_RESULT, [2]interface{}{rpc.OK, accinfo})

}

func (a *Account) findNest(s *Session) *service.Srv {
	srv := a.ctx.Core.LookupMinLoadByType("nest")
	if srv == nil {
		s.Error(share.S2C_ERR_SERVICE_INVALID)
		return nil
	}

	err := a.ctx.Core.MailtoAndCallback(nil, srv.Mailbox(), "Account.Logged", a.OnNestLogged, s.id, s.Account)
	if err != nil {
		s.Error(share.S2C_ERR_SERVICE_INVALID)
		return nil
	}
	return srv
}

func (a *Account) OnNestLogged(p interface{}, e *rpc.Error, ar *utils.LoadArchive) {
	id := p.(uint64)
	session := a.ctx.FindSession(id)
	if session == nil {
		a.ctx.Core.LogErr("session not found", id)
		return
	}

	if e != nil && protocol.CheckRpcError(e) {
		session.Dispatch(NEST_RESULT, [2]interface{}{e.ErrCode(), ""})
		return
	}

	var token string
	ar.Read(&token)

	session.Dispatch(NEST_RESULT, [2]interface{}{rpc.OK, token})
}
