package protocol

const (
	SLAVE_REGISTER = 1 + iota // SLAVE注册
)

type SlaveRegister struct {
	Id   int    // slave id
	Addr string // listen addr
	Port int    // listen port
}
