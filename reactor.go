package antgo

type Reactor interface {
	// OnOpen(c *Connection) bool
	// OnBeat(c *Connection) bool
	OnConnect(c *Connection) bool
	OnMessage(c *Connection, pt Packet) bool
	OnClose(c *Connection)
}
