package main

import (
	"../../antgo/chatroom"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	gateway := chatroom.NewGateway("tcp4", "127.0.0.1", 2400, "tcp",
		"tcp4", "127.0.0.1", []int{2301, 2302, 2303, 2304}, "tcp",
		"tcp4", "127.0.0.1", 2300, "tcp")
	gateway.Run()
}
