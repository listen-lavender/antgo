package main

import (
	"../../antgo/chatroom"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	register := chatroom.NewRegister("tcp4", "127.0.0.1", 2300, "tcp", "tcp")
	register.Run()
}
