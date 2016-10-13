package main

import (
	"../../antgo"
	"../../antgo/listener"
	"../../antgo/protocol"
	"../../antgo/reactor"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	echoListener := listener.NewTCPListener("echo4", "127.0.0.1", 2300)

	// creates a server
	config := &antgo.Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
	}
	echoReactor := reactor.EchoReactor{}
	echoProtocol := protocol.NewEchoProtocol(echoListener)
	srv := antgo.NewServer(config, echoReactor, echoProtocol)

	go srv.Start(time.Second)
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	srv.Stop()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
