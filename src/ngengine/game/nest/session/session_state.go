package session

import (
	"ngengine/common/fsm"
	"ngengine/game/gameobject/entity/inner"
	"ngengine/protocol/proto/c2s"
	"ngengine/share"
)

const (
	NONE      = iota
	TIMER     // 1秒钟的定时器
	BREAK     // 客户端断开连接
	LOGIN     // 客户端登录
	ROLE_INFO // 角色列表
	CREATE    // 创建角色
	CREATED   // 创建完成
)

type idlestate struct {
	fsm.Default
	owner *Session
	Idle  int32
}

func (s *idlestate) Handle(event int, param interface{}) string {
	switch event {
	case LOGIN:
		token := param.(string)
		if s.owner.ValidToken(token) {
			if !s.owner.QueryRoleInfo() {
				s.owner.Error(-1)
				return ""
			}
			return "logged"
		}
		// 验证失败直接踢下线
		s.owner.Break()
		return ""
	case TIMER:
		s.Idle++
		if s.Idle > 60 {
			s.owner.Break()
			return ""
		}
	case BREAK:
		s.owner.DestroySelf()
		return fsm.STOP
	default:
		s.owner.ctx.core.LogWarnf("idle state receive error event(%d)", event)
	}
	return ""
}

type logged struct {
	fsm.Default
	owner *Session
	Idle  int32
}

func (s *logged) Handle(event int, param interface{}) string {
	switch event {
	case TIMER:
		s.Idle++
		if s.Idle > 60 {
			s.owner.Break()
			return ""
		}
	case BREAK:
		s.owner.DestroySelf()
		return fsm.STOP
	case ROLE_INFO:
		args := param.([2]interface{})
		errcode := args[0].(int32)
		roles := args[1].([]*inner.Role)
		if errcode != 0 {
			s.owner.Error(errcode)
			return ""
		}

		s.owner.OnRoleInfo(roles)
		s.Idle = 0
	case CREATE:
		args := param.(c2s.CreateRole)
		s.owner.CreateRole(args)
		return "createrole"
	default:
		s.owner.ctx.core.LogWarnf("logged state receive error event(%d)", event)
	}
	return ""
}

type createrole struct {
	fsm.Default
	owner *Session
	Idle  int32
}

func (c *createrole) Handle(event int, param interface{}) string {
	switch event {
	case CREATED:
		errcode := param.(int32)
		if errcode != 0 {
			c.owner.Error(errcode)
			return "logged"
		}
		if !c.owner.QueryRoleInfo() {
			c.owner.Error(-1)
		}
		return "logged"
	case TIMER:
		c.Idle++
		if c.Idle > 60 {
			c.owner.Error(share.ERR_CREATE_TIMEOUT)
			return "logged"
		}
	default:
		c.owner.ctx.core.LogWarnf("create role state receive error event(%d)", event)
	}
	return ""
}
