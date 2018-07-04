package session

import (
	"ngengine/common/fsm"
)

const (
	NONE      = iota
	TIMER     // 1秒钟的定时器
	BREAK     // 客户端断开连接
	LOGIN     // 客户端登录
	ROLE_INFO // 角色列表
	CREATE    // 创建角色
	CREATED   // 创建完成
	CHOOSE    // 选择角色
	CHOOSED   // 选择角色成功
)

const (
	SIDLE   = "idle"
	SLOGGED = "logged"
	SCREATE = "create"
	SCHOOSE = "choose"
	SONLINE = "online"
)

func initState(s *Session) *fsm.FSM {
	fsm := fsm.NewFSM()
	fsm.Register(SIDLE, &idlestate{owner: s})
	fsm.Register(SLOGGED, &logged{owner: s})
	fsm.Register(SCREATE, &createrole{owner: s})
	fsm.Register(SCHOOSE, &chooserole{owner: s})
	fsm.Register(SONLINE, &online{owner: s})
	fsm.Start(SIDLE)
	return fsm
}
