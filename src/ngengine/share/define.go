package share

type SessionId uint64

type ServiceId uint16

const (
	MAX_BUF_LEN     = 1024 * 16      // 消息缓冲区大小16k
	SID_MAX         = 0xFFFF         // service id 最大值
	SESSION_MAX     = 0x7FFFFFFFFFFF // session最大值
	OBJECT_ID_MAX   = 0xFFFFFFFFFF   // objectid最大值
	OBJECT_TYPE_MAX = 0x7F           // object type最大值
)

// mailbox类型定义
const (
	MB_FLAG_APP    = iota // app
	MB_FLAG_CLIENT        // client
)

// object类型定义
const (
	OBJECT_TYPE_NONE   = iota
	OBJECT_TYPE_OBJECT // 对象
	OBJECT_TYPE_GHOST  // 对象副本
	OBJECT_TYPE_SHARE  // 共享数据
)

const (
	ST_ADMIN = 1 + iota
)

// 事件定义
const (
	EVENT_READY        = "ready"        // 别的服务就绪, args:{id:ServiceId}
	EVENT_LOST         = "lost"         // 丢失服务, args:{id:ServiceId}
	EVENT_USER_CONNECT = "user_connect" // 玩家连接, args:{id:uint64}
	EVENT_USER_LOST    = "user_lost"    // 玩家断开连接, args:{id:uint64}
)

// MessageBox 结构体消息包
type MessageBox struct {
	Message interface{}
}
