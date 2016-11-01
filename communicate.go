package antgo

import (
	"net"
	"time"
)

type ListenSpeaker interface {
	ReadPacket(netConn net.Conn, endTag []byte) <-chan string
	SetDeadline(t time.Time) (err error)
	Accept() (net.Conn, error)
	Dial() (net.Conn, error)
	Close()
}
