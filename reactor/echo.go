package reactor

import (
	"../../antgo"
	"../../antgo/protocol"
	"fmt"
	"time"
)

type EchoReactor struct {
}

func (p EchoReactor) OnConnect(c *antgo.Connection) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	fmt.Println("OnConnect:", addr)
	return true
}

func (p EchoReactor) OnMessage(c *antgo.Connection, pt antgo.Packet) bool {
	echoPacket := pt.(*protocol.EchoPacket)
	fmt.Printf("OnMessage:[%v] [%v]\n", echoPacket.GetLength(), string(echoPacket.GetBody()))
	c.AsyncWritePacket(protocol.NewEchoPacket(echoPacket.Serialize(), true), time.Second)
	return true
}

func (p EchoReactor) OnClose(c *antgo.Connection) {
	fmt.Println("OnClose:", c.GetExtraData())
}
