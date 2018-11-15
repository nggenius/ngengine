package protocol

const (
	S2C_ERROR = 10000 + iota // 系统错误
	S2C_RPC                  // 远程调用
)

const (
	C2S_PING = 20000 + iota // 心跳
	C2S_RPC                 // 远程调用
)

type S2CMsg struct {
	Sender string
	To     uint64
	Method string
	Data   []byte
}

type S2CBrocast struct {
	Sender string
	To     []uint64
	Method string
	Data   []byte
}

type S2SBrocast struct {
	To     []uint64
	Method string
	Args   interface{}
}
