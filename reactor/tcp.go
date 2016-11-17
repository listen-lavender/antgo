package reactor

import (
	"../../antgo"
	"fmt"
)

type TCPReactor struct {
}

func (p TCPReactor) OnConnect(c *antgo.Conn) string {
	addr := c.RemoteAddr()
	fmt.Println("OnConnect:", addr)
	return addr
}

func (p TCPReactor) OnMessage(c *antgo.Conn, pt antgo.Packet) bool {
	code := pt.Code()
	event := pt.Event()
	msg := pt.Msg()

	fmt.Println("code:", code)
	fmt.Println("event:", event)
	fmt.Println("msg:", msg)

	return true
}

func (p TCPReactor) OnClose(c *antgo.Conn) {
	addr := c.RemoteAddr()
	fmt.Println("OnClose:", addr)
	c.Close()
}
