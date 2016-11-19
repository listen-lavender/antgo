package main

import (
	"../../antgo/chatroom"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	worker := chatroom.NewWorker("tcp4", "127.0.0.1", 2300, "tcp")
	worker.Run()
}
