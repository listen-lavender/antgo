package antgo

import (
	"net"
	"time"
)

type Listener interface {
	ReadPacket(conn net.Conn, endTag []byte) <-chan string
	SetDeadline(t time.Time) (err error)
	Accept() (net.Conn, error)
	Close()
}
