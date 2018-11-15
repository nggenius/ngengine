package protocol

import (
	"net"
	"github.com/nggenius/ngengine/logger"
	"runtime"
	"strings"
)

type TCPHandler interface {
	Handle(net.Conn)
}

func TCPServer(listener net.Listener, handler TCPHandler, l *logger.Log) {
	l.LogInfof("TCP: listening on %s", listener.Addr())
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				l.LogWarnf("NOTICE: temporary Accept() failure - %s", err)
				runtime.Gosched()
				continue
			}
			if !strings.Contains(err.Error(), "use of closed network connection") {
				l.LogErrf("ERROR: listener.Accept() - %s", err)
			}
			break
		}
		go handler.Handle(clientConn)
	}

	l.LogInfof("TCP: closing %s", listener.Addr())
}
