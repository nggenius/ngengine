package s2c

type Rpc struct {
	Sender        string
	Servicemethod string
	Data          []byte
}

type Error struct {
	ErrCode int32
}

type NestInfo struct {
	Addr  string
	Port  int32
	Token string
}

type Role struct {
	Index int8
	Name  string
}

type RoleInfo struct {
	Roles []Role
}
