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
	tcpListener := listener.NewTCPListener("tcp4", "127.0.0.1", 2300)

	// creates a server
	config := &antgo.Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
	}
	tcpReactor := reactor.TCPReactor{}
	tcpProtocol := protocol.NewTCPProtocol(tcpListener)
	srv := antgo.NewServer(config, tcpReactor, tcpProtocol)

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
