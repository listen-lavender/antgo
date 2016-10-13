package main

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "runtime"
    "syscall"
    "time"
    "../../antgo"
    "../../antgo/protocol"
    "../../antgo/listener"
    "../../antgo/reactor"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())

    // creates a tcp listener
    // tcpAddr, err := net.ResolveTCPAddr("tcp4", ":23")
    // checkError(err)
    // listener, err := net.ListenTCP("tcp", tcpAddr)
    // checkError(err)
    tcpListener := listener.NewTCPListener("tcp4", "127.0.0.1", 23)


    // creates a server
    config := &antgo.Config{
        PacketSendChanLimit:    20,
        PacketReceiveChanLimit: 20,
    }
    tcpReactor := reactor.TCPReactor{}
    tcpProtocol := protocol.NewTCPProtocol(tcpListener)
    srv := antgo.NewServer(config, tcpReactor, tcpProtocol)

    // starts service
    go srv.Start(time.Second)
    // fmt.Println("listening:", listener.Addr())

    // catchs system signal
    chSig := make(chan os.Signal)
    signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
    fmt.Println("Signal: ", <-chSig)

    // stops service
    srv.Stop()
}

func checkError(err error) {
    if err != nil {
        log.Fatal(err)
    }
}
