package chatroom

import (
	"../../antgo"
	"../../antgo/multinet"
	"../../antgo/protocol"
	"bytes"
	"strings"
)

var Handlers map[string]func(conn *antgo.Conn, data string, worker *Worker)

func NewListenSpeaker(Ltype string, Transport string, IP string, Port int) antgo.ListenSpeaker {
	var buffer bytes.Buffer
	buffer.WriteString("new")
	buffer.WriteString(strings.ToLower(Ltype))
	buffer.WriteString("listenspeaker")
	// buffer.Reset()
	return multinet.Indexes[buffer.String()](Transport, IP, Port)
}

func NewProtocol(Ptype string, listenspeaker antgo.ListenSpeaker) antgo.Protocol {
	var buffer bytes.Buffer
	buffer.WriteString("new")
	buffer.WriteString(strings.ToLower(Ptype))
	buffer.WriteString("protocol")
	// buffer.Reset()
	return protocol.Indexes[buffer.String()](listenspeaker)
}
