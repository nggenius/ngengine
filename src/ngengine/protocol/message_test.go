package protocol

import (
	"fmt"
	"testing"
)

func TestNewMessage(t *testing.T) {
	msg := NewMessage(16)
	fmt.Println(len(msg.Header), cap(msg.Header))
	msg.Header = msg.Header[:8]
	fmt.Println(len(msg.Header), cap(msg.Header))
}
