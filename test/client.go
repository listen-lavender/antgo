package main

import (
	"../../antgo/protocol"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8989")
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	echoProtocol := &protocol.EchoProtocol{}

	// ping <--> pong
	for i := 0; i < 3; i++ {
		// write
		conn.Write(protocol.NewEchoPacket([]byte("hello"), false).Serialize())

		// read
		q := echoProtocol.ReadPacket(conn)
		for p := range q {
			echoPacket := p.(*protocol.EchoPacket)
			fmt.Printf("Server reply:[%v] [%v]\n", echoPacket.GetLength(), string(echoPacket.GetBody()))
			time.Sleep(2 * time.Second)
		}
	}

	conn.Close()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
