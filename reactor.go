package antgo

type Reactor interface {
	OnConnect(c *Conn) string
	OnMessage(c *Conn, pt Packet) bool
	OnClose(c *Conn)
}
