package protocol

import (
	"../../antgo"
)

var Indexes map[string]func(listenspeaker antgo.ListenSpeaker) antgo.Protocol

func init() {
	Indexes = make(map[string]func(listenspeaker antgo.ListenSpeaker) antgo.Protocol)
	// Indexes["NewEchoProtocol"] = NewEchoProtocol
	Indexes["newtcpprotocol"] = NewTCPProtocol
}
