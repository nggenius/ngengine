package extension

import (
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
		return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, err.Error(), tag)
	}
	err = m.Read(&player)
	if err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_ARGS_ERROR, err.Error(), tag)
	}

	session := r.store.Sql().Session()
	defer session.Close()

	count, err := session.Count(inner.Role{Index: role.Index, Account: role.Account})
	if err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_STORE_SQL, err.Error(), tag)
	}

	if count != 0 {
		return protocol.ReplyError(protocol.TINY, share.ERR_STORE_ROLE_INDEX, "index error", tag)
	}

	session.Begin()
	_, err = session.Insert(&role)
	if err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_STORE_ERROR, err.Error(), tag)
	}
	_, err = session.Insert(&player)
	if err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_STORE_ERROR, err.Error(), tag)
	}

	session.Commit()
	return protocol.Reply(protocol.TINY, tag)
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
		return protocol.ReplyError(protocol.TINY, share.ERR_STORE_ERROR, err.Error(), tag)
	}

	if count == 0 {
		return protocol.ReplyError(protocol.TINY, share.ERR_STORE_ROLE_NOT_FOUND, "player not found", tag)
	}

	session.Begin()

	// 备份数据
	_, err = session.Exec("insert into player_bak select * from player where id=?", roleid)
	if err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_STORE_ERROR, err.Error(), tag)
	}
	// 删除
	_, err = session.Delete(&role)
	if err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_STORE_ERROR, err.Error(), tag)
	}

	_, err = session.Delete(&player)
	if err != nil {
		return protocol.ReplyError(protocol.TINY, share.ERR_STORE_ERROR, err.Error(), tag)
	}
	session.Commit()
	return protocol.Reply(protocol.TINY, tag)
}
