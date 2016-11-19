package chatroom

import (
	"../../antgo"
	"../../antgo/multinet"
	"../../antgo/protocol"
	"bytes"
	"strings"
	"time"
)

var Timeout time.Duration

var Handlers map[string]func(conn *antgo.Conn, data string, worker *Worker)

func init() {
	Timeout = 10 * time.Second
}

func NewListenDialer(Ltype string, Transport string, IP string, Port int) antgo.ListenDialer {
	var buffer bytes.Buffer
	buffer.WriteString("new")
	buffer.WriteString(strings.ToLower(Ltype))
	buffer.WriteString("listendialer")
	// buffer.Reset()
	return multinet.Indexes[buffer.String()](Transport, IP, Port)
}

func NewProtocol(Ptype string, listendialer antgo.ListenDialer) antgo.Protocol {
	var buffer bytes.Buffer
	buffer.WriteString("new")
	buffer.WriteString(strings.ToLower(Ptype))
	buffer.WriteString("protocol")
	// buffer.Reset()
	return protocol.Indexes[buffer.String()](listendialer, Ptype)
}
