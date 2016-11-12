package reactor

import (
	"../../antgo"
	"../../antgo/protocol"
	"fmt"
	"net"
)

type TCPReactor struct {
}

func (p TCPReactor) OnConnect(c *antgo.Conn) net.Addr {
	addr := c.RemoteAddr()
	fmt.Println("OnConnect:", addr)
	// c.AsyncWritePacket(protocol.NewTCPPacket("prompt", []byte("Welcome to p TCP Server")), 0)
	return addr
}

func (p TCPReactor) OnMessage(c *antgo.Conn, pt antgo.Packet) bool {
	code := pt.Code()
	event := pt.Event()
	msg := pt.Msg()

	switch event {
	case "echo":
		c.AsyncWritePacket(protocol.NewTCPPacket(code, "echo", msg), 0)
	case "login":
		c.AsyncWritePacket(protocol.NewTCPPacket(code, "login", msg), 0)
	case "prompt":
		c.AsyncWritePacket(protocol.NewTCPPacket(code, "prompt", msg), 0)
	case "quit":
		return false
	default:
		c.AsyncWritePacket(protocol.NewTCPPacket(code, event, msg), 0)
	}

	return true
}

func (p TCPReactor) OnClose(c *antgo.Conn) {
	addr := c.RemoteAddr()
	fmt.Println("OnClose:", addr)
	c.Close()
}
