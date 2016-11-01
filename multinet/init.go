package multinet

import (
	"../../antgo"
)

var Indexes map[string]func(netType string, ip string, port int) antgo.ListenSpeaker

func init() {
	Indexes = make(map[string]func(netType string, ip string, port int) antgo.ListenSpeaker)
	Indexes["newtcplistenspeaker"] = NewTCPListenSpeaker
	Indexes["newudplistenspeaker"] = NewUDPListenSpeaker
}
