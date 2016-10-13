package listener

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type UDPListener struct {
	addr     *net.UDPAddr
	listener *net.UDPConn
}

func NewUDPListener(netType string, ip string, port int) *UDPListener {
	addr, err := net.ResolveUDPAddr(netType, ip+":"+strconv.Itoa(port))
	if err != nil { //EOF, or worse
		fmt.Println(err)
		return nil
	}
	listener, err := net.ListenUDP(netType, addr)
	if err != nil { //EOF, or worse
		fmt.Println(err)
		return nil
	}
	fmt.Println(ip)
	fmt.Println(port)
	return &UDPListener{
		addr:     addr,
		listener: listener,
	}
}

func (p UDPListener) ReadUdpPacket(conn net.Conn, endTag []byte) <-chan string {
	msg := make(chan string)

	for {
		data := make([]byte, 1024)

		readLengh, err := conn.Read(data)

		if err != nil { //EOF, or worse
			return nil
		}

		if readLengh > 0 {
			msg <- string(data)
		}
	}
	return msg
}

func (p UDPListener) SetDeadline(t time.Time) (err error) {
	p.listener.SetDeadline(t)
	return nil
}

func (p UDPListener) Accept() (net.Conn, error) {
	return p.listener, nil
}

func (p UDPListener) Close() {
	p.listener.Close()
}
