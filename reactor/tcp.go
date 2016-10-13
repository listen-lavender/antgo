package reactor

import (
	"../../antgo"
	"../../antgo/protocol"
	"fmt"
)

type TCPReactor struct {
}

func (p TCPReactor) OnConnect(c *antgo.Connection) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	fmt.Println("OnConnect:", addr)
	c.AsyncWritePacket(protocol.NewTCPPacket("unknow", []byte("Welcome to p TCP Server")), 0)
	return true
}

func (p TCPReactor) OnMessage(c *antgo.Connection, pt antgo.Packet) bool {
	ptt := pt.(*protocol.TCPPacket)
	cmd := ptt.GetData()
	cmdType := ptt.GetType()

	switch cmdType {
	case "echo":
		c.AsyncWritePacket(protocol.NewTCPPacket("echo", cmd), 0)
	case "login":
		c.AsyncWritePacket(protocol.NewTCPPacket("login", []byte(string(cmd)+" has login")), 0)
	case "quit":
		return false
	default:
		c.AsyncWritePacket(protocol.NewTCPPacket("unknow", []byte("unknow cmd")), 0)
	}

	return true
}

func (p TCPReactor) OnClose(c *antgo.Connection) {
	fmt.Println("OnClose:", c.GetExtraData())
}
