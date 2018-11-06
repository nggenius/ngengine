package share

const (
	ERR_REPLY_SUCCEED        = iota // 成功
	ERR_TIME_OUT                    // 超时
	ERR_ARGS_ERROR                  // 参数错误
	ERR_SYSTEM_ERROR                // 系统错误
	ERR_RPC_FAILED                  // rpc错误
	ERR_REPLY_FAILED                // 失败
	ERR_CREATE_TIMEOUT              // 创建角色超时
	ERR_CHOOSE_ROLE                 // 选择角色出错
	ERR_CHOOSE_TIMEOUT              // 选择角色超时
	ERR_OBJECT_NOT_FOUND            // 对象没有找到
	ERR_OBJECT_RPC_CALL             // rpc 执行错误
	ERR_OBJECT_RPC_NOT_MATCH        // rpc 对象不匹配
	ERR_REGION_NOT_FOUND            // 场景未找到
	ERR_ENTER_REGION_FAILED         // 进入场景失败
)

// 存储错误码
const (
	ERR_STORE_NONE           = 11000 + iota
	ERR_STORE_SQL            // sql 错误
	ERR_STORE_NOROW          // 没有查找到记录
	ERR_STORE_ROLE_INDEX     // 索引错误
	ERR_STORE_ROLE_NOT_FOUND // 玩家没找到
	ERR_STORE_ERROR          // 其它错误
)

const (
	ERR_REGION_NONE          = 12000 + iota
	ERR_REGION_CREATE_FAILED // 创建场景失败
)

const (
	ERR_OBJECT_NONT      = 13000 + iota
	ERR_OBJECT_REPLICATE // 传送对象失败
	ERR_ORIGIN_NOT_FOUND // 对象没有找到
)
