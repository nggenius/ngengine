package protocol

import (
	. "ngengine/share"
	"time"
)

const (
	HB_INTERVAL = time.Second * 10
)

const (
	S2A_REGISTER   = 1 + iota //注册服务
	S2A_UNREGISTER            //注销服务
	S2A_WATCH                 //获取服务信息
	S2A_HEARTBEAT             //心跳

	A2S_SERVICES = 100 + iota //服务信息
)

const (
	ST_SYNC = 1 + iota //同步
	ST_ADD             //增加
	ST_DEL             //删除
)

type ServiceInfo struct {
	Id     ServiceId //服务ID
	Name   string    //服务名称
	Type   string    //服务类型
	Status int       //状态
	Addr   string    //ip地址
	Port   int       //端口号
}

type Register struct {
	Service ServiceInfo
}

type Watch struct {
	WatchType []string //需要获取的服务类型，[0]获取所有
}

type Services struct {
	OpType  int8
	All     bool
	Service []ServiceInfo
}
