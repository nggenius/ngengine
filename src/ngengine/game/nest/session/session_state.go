package session

import (
	"ngengine/common/fsm"
)

const (
	NONE      = iota
	ETIMER    // 1秒钟的定时器
	EBREAK    // 客户端断开连接
	ELOGIN    // 客户端登录
	EROLEINFO // 角色列表
	ECREATE   // 创建角色
	ECREATED  // 创建完成
	ECHOOSE   // 选择角色
	ECHOOSED  // 选择角色成功
	EDELETE   // 删除角色
	EDELETED  // 删除成功
	ESTORED   // 存档完成
	EONLINE   // 进入场景
	EFREGION  // 查找场景
	ESWREGION // 切换场景
)

const (
	SIDLE    = "idle"
	SLOGGED  = "logged"
	SCREATE  = "create"
	SCHOOSE  = "choose"
	SDELETE  = "delete"
	SONLINE  = "online"
	SOFFLINE = "offline"
	SLEAVING = "leaving"
)

func initState(s *Session) *fsm.FSM {
	fsm := fsm.NewFSM()
	fsm.Register(SIDLE, &idlestate{owner: s})
	fsm.Register(SLOGGED, &logged{owner: s})
	fsm.Register(SCREATE, &createrole{owner: s})
	fsm.Register(SCHOOSE, &chooserole{owner: s})
	fsm.Register(SDELETE, &deleting{owner: s})
	fsm.Register(SONLINE, &online{owner: s})
	fsm.Register(SOFFLINE, &offline{owner: s})
	fsm.Register(SLEAVING, &leaving{owner: s})
	fsm.Start(SIDLE)
	return fsm
}
