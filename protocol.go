package antgo

import (
	"net"
)

type Packet interface {
	Serialize() []byte
	Event() string
	Msg() []byte
}

type Protocol interface {
	ReadPacket(netConn net.Conn) <-chan Packet
	ListenSpeaker() ListenSpeaker
}
