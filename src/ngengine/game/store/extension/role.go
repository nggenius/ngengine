package extension

import (
	"errors"
	"ngengine/core/rpc"
	"ngengine/core/service"
	"ngengine/game/gameobject/entity"
	"ngengine/game/gameobject/entity/inner"
	"ngengine/module/store"
	"ngengine/protocol"
	"ngengine/share"
)

type Role struct {
	store *store.StoreModule
	ctx   service.CoreAPI
}

func NewRole(core service.CoreAPI, s *store.StoreModule) *Role {
	r := &Role{}
	r.ctx = core
	r.store = s
	return r
}

func (r *Role) RegisterCallback(svr rpc.Servicer) {
	svr.RegisterCallback("CreateRole", r.CreateRole)
	svr.RegisterCallback("DeleteRole", r.DeleteRole)
}

func ParseCreateRole(reply *protocol.Message) (errcode int32, err error, tag string) {
	errcode, ar := protocol.ParseReply(reply)
	tag, err = ar.ReadString()
	if err != nil {
		return
	}
	if int(errcode) == share.ERR_TIME_OUT {
		err = rpc.ErrTimeout
		return
	}
	if errcode != 0 {
		errstr, _ := ar.ReadString()
		err = errors.New(errstr)
		return
	}

	return
}

func (r *Role) CreateRole(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	tag, err := m.ReadString()
	if err != nil {
		r.ctx.LogFatal("read tag failed, ", err)
		return 0, nil
	}
	var role inner.Role
	var player entity.PlayerArchive
	err = m.Read(&role)
	if err != nil {
		return share.ERR_ARGS_ERROR, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}
	err = m.Read(&player)
	if err != nil {
		return share.ERR_ARGS_ERROR, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}

	session := r.store.Sql().Session()
	defer session.Close()

	count, err := session.Count(inner.Role{Index: role.Index, Account: role.Account})
	if err != nil {
		return share.ERR_STORE_SQL, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}

	if count != 0 {
		return share.ERR_STORE_ROLE_INDEX, protocol.ReplyMessage(protocol.TINY, tag, "index error")
	}

	session.Begin()
	_, err = session.Insert(&role)
	if err != nil {
		return share.ERR_STORE_ERROR, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}
	_, err = session.Insert(&player)
	if err != nil {
		return share.ERR_STORE_ERROR, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}

	session.Commit()
	return share.ERR_REPLY_SUCCEED, protocol.ReplyMessage(protocol.TINY, tag)
}

func ParseDeleteRole(reply *protocol.Message) (errcode int32, err error, tag string) {
	errcode, ar := protocol.ParseReply(reply)
	tag, err = ar.ReadString()
	if err != nil {
		return
	}
	if int(errcode) == share.ERR_TIME_OUT {
		err = rpc.ErrTimeout
		return
	}
	if errcode != 0 {
		errstr, _ := ar.ReadString()
		err = errors.New(errstr)
		return
	}

	return
}

func (r *Role) DeleteRole(sender, _ rpc.Mailbox, msg *protocol.Message) (errcode int32, reply *protocol.Message) {
	m := protocol.NewMessageReader(msg)
	tag, err := m.ReadString()
	if err != nil {
		r.ctx.LogFatal("read tag failed, ", err)
		return 0, nil
	}
	roleid, err := m.ReadInt64()
	if err != nil {
		r.ctx.LogFatal("read roleid failed, ", err)
		return 0, nil
	}

	session := r.store.Sql().Session()
	defer session.Close()

	var role inner.Role
	var player entity.PlayerArchive
	role.Id = roleid
	player.Id = roleid

	count, err := session.Count(player)
	if err != nil {
		return share.ERR_STORE_ERROR, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}

	if count == 0 {
		return share.ERR_STORE_ROLE_NOT_FOUND, protocol.ReplyMessage(protocol.TINY, tag, "player not found")
	}

	session.Begin()

	// 备份数据
	_, err = session.Exec("insert into player_bak select * from player where id=?", roleid)
	if err != nil {
		return share.ERR_STORE_ERROR, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}
	// 删除
	_, err = session.Delete(&role)
	if err != nil {
		return share.ERR_STORE_ERROR, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}

	_, err = session.Delete(&player)
	if err != nil {
		return share.ERR_STORE_ERROR, protocol.ReplyMessage(protocol.TINY, tag, err.Error())
	}
	session.Commit()
	return share.ERR_REPLY_SUCCEED, protocol.ReplyMessage(protocol.TINY, tag)
}
