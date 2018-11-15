package protocol

import "net"

var (
	MagicV1 = []byte("  V1")
)

const (
	MAX_ADMIN_MESSAGE_SIZE = 0x1000
)

type Protocol interface {
	IOLoop(conn net.Conn) error
}
