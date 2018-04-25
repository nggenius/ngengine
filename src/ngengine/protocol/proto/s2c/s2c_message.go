package s2c

type Rpc struct {
	Sender        string
	Servicemethod string
	Data          []byte
}

type Error struct {
	ErrCode int32
}
