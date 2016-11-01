package antgo

import (
	"net"
)

type Reactor interface {
	OnConnect(c *Conn) net.Addr
	OnMessage(c *Conn, pt Packet) bool
	OnClose(c *Conn)
}
