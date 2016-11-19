package antgo

import (
	"net"
)

type Packet interface {
	Serialize() []byte
	Code() int
	Event() string
	Msg() interface{}
}

type Protocol interface {
	ReadPacket(netConn net.Conn) Packet
	Deserialize(code int, event string, msg interface{}) Packet
	ListenDialer() ListenDialer
    Type() string
}
