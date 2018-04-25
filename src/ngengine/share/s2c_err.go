package share

const (
	S2C_ERR_SUCCEED         = 0           // 正常
	S2C_ERR_SERVICE_INVALID = 1000 + iota // 服务不可用
	S2C_ERR_NAME_PASS                     // 帐户密码错误
)
