package multinet

import (
	"../../antgo"
)

var Indexes map[string]func(netType string, ip string, port int) antgo.ListenDialer

func init() {
	Indexes = make(map[string]func(netType string, ip string, port int) antgo.ListenDialer)
	Indexes["newtcplistendialer"] = NewTCPListenDialer
	Indexes["newudplistendialer"] = NewUDPListenDialer
}
