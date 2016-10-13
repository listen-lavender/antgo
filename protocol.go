package antgo

import (
	"net"
)

type Packet interface {
	Serialize() []byte
	// GetType() string
	// GetData() []byte
}

type Protocol interface {
	ReadPacket(conn net.Conn) <-chan Packet
	GetListener() Listener
	// SetListener(listener Listener)
}
