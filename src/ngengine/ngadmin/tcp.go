package ngadmin

import (
	"io"
	"net"
	"ngengine/protocol"
)

type tcpServer struct {
	ctx *Context
}

// Handle tcp连接处理函数
func (p *tcpServer) Handle(clientConn net.Conn) {
	if p.ctx.ngadmin.quit { // 已经退出了
		clientConn.Close()
		return
	}

	p.ctx.ngadmin.LogInfof("TCP: new client(%s)", clientConn.RemoteAddr())

	buf := make([]byte, 4)
	_, err := io.ReadFull(clientConn, buf)
	if err != nil {
		p.ctx.ngadmin.LogErrf("failed to read protocol version - %s", err)
		return
	}
	protocolMagic := string(buf)

	var prot protocol.Protocol
	switch protocolMagic {
	case string(protocol.MagicV1):
		prot = &AdminProtocol{ctx: p.ctx}
	default:
		clientConn.Close()
		p.ctx.ngadmin.LogErrf("client(%s) bad protocol magic '%s'",
			clientConn.RemoteAddr(), protocolMagic)
		return
	}

	err = prot.IOLoop(clientConn)
	if err != nil {
		p.ctx.ngadmin.LogErrf("client(%s) - %s", clientConn.RemoteAddr(), err)
	}

	p.ctx.ngadmin.LogInfof("TCP: client(%s) offline", clientConn.RemoteAddr())
	clientConn.Close()
}
