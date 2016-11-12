package antgo

import (
	"net"
	"time"
)

type ListenDialer interface {
	ReadPacket(netConn net.Conn, endTag []byte) []byte
	SetDeadline(t time.Time) (err error)
	Listen()
	Accept() (net.Conn, error)
	Dial() (net.Conn, error)
	Close()
}
