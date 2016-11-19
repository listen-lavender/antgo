package protocol

import (
	"../../antgo"
)

var Indexes map[string]func(listendialer antgo.ListenDialer, protocoltype string) antgo.Protocol

func init() {
	Indexes = make(map[string]func(listendialer antgo.ListenDialer, protocoltype string) antgo.Protocol)
	// Indexes["NewEchoProtocol"] = NewEchoProtocol
	Indexes["newtcpprotocol"] = NewTCPProtocol
}
