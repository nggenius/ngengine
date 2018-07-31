package protocol

import (
	. "ngengine/share"
	"time"
)

const (
	HB_INTERVAL = time.Second * 10
)

const (
	S2A_REGISTER   = 1 + iota // 注册服务
	S2A_UNREGISTER            // 注销服务
	S2A_WATCH                 // 获取服务信息
	S2A_HEARTBEAT             // 心跳
	S2A_LOAD                  // 更新负载
	S2A_READY                 // 服务准备好了
)

const (
	A2S_SERVICES      = 100 + iota // 服务信息
	A2S_LOAD                       // 同步负载信息
	A2S_SERVICE_READY              // 服务就绪
	A2S_ALL_READY                  // 准备好
	A2S_SERVICE_CLOSE              // 关闭服务
)

const (
	ST_SYNC = 1 + iota // 同步
	ST_ADD             // 增加
	ST_DEL             // 删除
)

type ServiceInfo struct {
	Id        ServiceId // 服务ID
	Name      string    // 服务名称
	Type      string    // 服务类型
	Status    int8      // 状态
	Addr      string    // ip地址
	Port      int       // 端口号
	OuterAddr string    // 外网地址
	OuterPort int       // 外网端口
	Load      int32     // 负载情况
}

type LoadInfo struct {
	Id   ServiceId // 服务ID
	Load int32     // 负载情况
}

type Register struct {
	AdminId int // admin id
	Service ServiceInfo
}

type Watch struct {
	WatchType []string // 需要获取的服务类型，[0]获取所有
}

type Services struct {
	OpType  int8
	All     bool
	Service []ServiceInfo
}

type SeverClosing struct {
	ID        uint16 // 关闭的服务器id
	SeverName string //关闭的服务器名字
}

type ServiceReady struct {
	Id ServiceId // 服务ID
}
