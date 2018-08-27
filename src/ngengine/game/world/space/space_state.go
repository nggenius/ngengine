package space

import "ngengine/common/fsm"

const (
	NONE            = iota
	ETIMER          // 1秒钟的定时器
	EREGION_RESP    // region 响应
	EREGION_CREATED // region 创建成功
)

const (
	SIDLE   = "idle"
	SCREATE = "create_region"
)

func initState(s *SpaceManage) *fsm.FSM {
	fsm := fsm.NewFSM()
	fsm.Register(SIDLE, newIdle(s))
	fsm.Register(SCREATE, newCreateRegion(s))
	fsm.Start(SIDLE)
	return fsm
}
