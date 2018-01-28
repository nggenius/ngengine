package share

const (
	ST_ADMIN = 1 + iota
)

const (
	MAX_BUF_LEN = 1024 * 16      // 消息缓冲区大小16k
	SESSION_MAX = 0x7FFFFFFFFFFF // session最大值
)

const (
	EVENT_READY        = "ready"        // 别的服务就绪, args:{id:ServiceId}
	EVENT_LOST         = "lost"         // 丢失服务, args:{id:ServiceId}
	EVENT_USER_CONNECT = "user_connect" // 玩家连接, args:{id:uint64}
	EVENT_USER_LOST    = "user_lost"    // 玩家断开连接, args:{id:uint64}
)
