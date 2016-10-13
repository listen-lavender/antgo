package listener

import (
	"fmt"
	"net"
	"strconv"
	"bytes"
	"time"
)

type TCPListener struct {
	addr *net.TCPAddr
	listener *net.TCPListener
}

func NewTCPListener(netType string, ip string, port int) *TCPListener {
	addr, err := net.ResolveTCPAddr(netType, ip + ":" + strconv.Itoa(port))
	listener, err := net.ListenTCP(netType, addr)
	fmt.Println(err)
	fmt.Println(ip)
	fmt.Println(port)

	return &TCPListener{
		addr:     addr,
		listener: listener,
	}
}

func (p TCPListener) ReadPacket(conn net.Conn, endTag []byte) (<-chan string) {
	fullBuf := bytes.NewBuffer([]byte{})
	msg := make(chan string)

	for {
		data := make([]byte, 1024)

		readLengh, err := conn.Read(data)

		if err != nil { //EOF, or worse
			fmt.Println(err)
		}

		if readLengh > 0 {
			fullBuf.Write(data[:readLengh])

			index := bytes.Index(fullBuf.Bytes(), endTag)
			if index > -1 {
				command := fullBuf.Next(index)
				fullBuf.Next(2)
				//fmt.Println(string(command))
				msg <- string(command)
			}
		}
	}
	return msg
}

func (p TCPListener) SetDeadline(t time.Time) (err error){
	p.listener.SetDeadline(t)
	return nil
}

func (p TCPListener) Accept()(net.Conn, error){
	return p.listener.AcceptTCP()
}

func (p TCPListener) Close(){
	p.listener.Close()
}
