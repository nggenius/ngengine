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
)

type Account struct {
	ctx *LoginModule
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

	session := a.ctx.FindSession(mailbox.Id)
	if session == nil {
		a.ctx.core.LogErr("session not found", mailbox.Id)
		return 0, nil
	}

	session.SetMailbox(&mailbox)
	session.SetAccount(login.Name)
	session.Dispatch(LOGIN, login)
	return 0, nil
}

func (a *Account) sendLogin(s *Session, login *c2s.Login) error {
	if a.ctx.db == nil {
		s.Error(share.S2C_ERR_SERVICE_INVALID)
		return fmt.Errorf("database is not ready")
	}

	if err := a.ctx.storeClient.Get(
		s.Mailbox.String(),
		"inner.Account",
		map[string]interface{}{
			"Account=?":  login.Name,
			"Password=?": login.Pass,
		},
		a.LoginResult); err != nil {
		s.Error(share.S2C_ERR_SERVICE_INVALID)
		return err
	}

	return nil
}

func (a *Account) LoginResult(reply *protocol.Message) {
	accinfo := &inner.Account{}
	errcode, err, tag := store.ParseGetReply(reply, accinfo)
	if err != nil {
		a.ctx.core.LogErr(err)
		return
	}
	mailbox, err1 := rpc.NewMailboxFromStr(tag)
	if err1 != nil {
		a.ctx.core.LogErr(err1)
		return
	}

	session := a.ctx.FindSession(mailbox.Id)
	if session == nil {
		a.ctx.core.LogErr("session not found", mailbox.Id)
		return
	}

	session.Dispatch(LOGIN_RESULT, [2]interface{}{errcode, accinfo})
}

func (a *Account) findNest(s *Session) *service.Srv {
	srv := a.ctx.core.LookupMinLoadByType("nest")
	if srv == nil {
		s.Error(share.S2C_ERR_SERVICE_INVALID)
		return nil
	}

	err := a.ctx.core.MailtoAndCallback(nil, srv.Mailbox(), "Account.Logged", a.OnNestLogged, s.id, s.Account)
	if err != nil {
		s.Error(share.S2C_ERR_SERVICE_INVALID)
		return nil
	}
	return srv
}

func (a *Account) OnNestLogged(msg *protocol.Message) {
	errcode, ar := protocol.ParseReply(msg)
	var id uint64
	var token string
	err := ar.Read(&id)
	if err != nil {
		a.ctx.core.LogErr("read id failed")
		return
	}
	session := a.ctx.FindSession(id)
	if session == nil {
		a.ctx.core.LogErr("session not found", id)
		return
	}

	ar.Read(&token)

	session.Dispatch(NEST_RESULT, [2]interface{}{errcode, token})
}
