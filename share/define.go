package share

type SessionId uint64

type ServiceId uint16

const (
	MAX_BUF_LEN           = 1024 * 16      // 消息缓冲区大小16k
	INNER_MESSAGE_BUF_LEN = 1024 * 1024    // 内部消息缓冲区大小1M
	SID_MAX               = 0xFFFF         // service id 最大值
	SESSION_MAX           = 0x7FFFFFFFFFFF // session最大值
	OBJECT_ID_MAX         = 0xFFFFFFFFFF   // objectid最大值
	OBJECT_TYPE_MAX       = 0x7FFF         // object type最大值
	OBJECT_MAX            = 0xFFFFFF       // 最大对象数量
)

// mailbox类型定义
const (
	MB_FLAG_APP    = iota // app
	MB_FLAG_CLIENT        // client
)

const (
	ST_ADMIN = 1 + iota
)

// 事件定义
const (
	EVENT_SERVICE_READY      = "svr_ready"       // 别的服务就绪, args:{id:ServiceId}
	EVENT_SERVICE_LOST       = "lost"            // 丢失服务, args:{id:ServiceId}
	EVENT_USER_CONNECT       = "user_connect"    // 玩家连接, args:{id:uint64}
	EVENT_USER_LOST          = "user_lost"       // 玩家断开连接, args:{id:uint64}
	EVENT_SHUTDOWN           = "svr_shutdown"    // 关闭系统
	EVENT_MUST_SERVICE_READY = "must_ready"      // 必须启动的服务器就绪
	EVENT_ADMIN_CONNECTED    = "admin_connected" // admin建立连接成功
)

const (
	ROUTER_TO_OBJECT = "ObjectRouter.ToObject" // 对象消息路由
)
