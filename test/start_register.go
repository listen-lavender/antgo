package main

import (
	// "../../antgo"
	"../../antgo/chatroom"
	// "../../antgo/multinet"
	// "../../antgo/protocol"
	// "../../antgo/reactor"
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
	// tcpListenSpeaker := listenspeaker.NewTCPListenSpeaker("tcp4", "127.0.0.1", 2300)

	// // creates a server
	// config := &antgo.Config{
	//     PacketSendChanLimit:    20,
	//     PacketReceiveChanLimit: 20,
	// }
	// tcpReactor := reactor.TCPReactor{}
	// tcpProtocol := protocol.NewTCPProtocol(tcpListenSpeaker)
	// srv := antgo.NewServer(config, tcpReactor, tcpProtocol)
	register := chatroom.NewRegister("tcp4", "127.0.0.1", 2300, "tcp", "tcp", 20, 20)

	// go srv.Start(time.Second)
	go register.Listen(time.Second)
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	register.Stop()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
