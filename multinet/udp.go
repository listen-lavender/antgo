package multinet

import (
	"../../antgo"
	"fmt"
	"net"
	"strconv"
	"time"
)

type UDPListenSpeaker struct {
	netType  string
	addr     *net.UDPAddr
	listener *net.UDPConn
}

func NewUDPListenSpeaker(netType string, ip string, port int) antgo.ListenSpeaker {
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
	return &UDPListenSpeaker{
		netType:  netType,
		addr:     addr,
		listener: listener,
	}
}

func (p UDPListenSpeaker) ReadPacket(netConn net.Conn, endTag []byte) <-chan string {
	msg := make(chan string)

	for {
		data := make([]byte, 1024)

		readLengh, err := netConn.Read(data)

		if err != nil { //EOF, or worse
			return nil
		}

		if readLengh > 0 {
			msg <- string(data)
		}
	}
	return msg
}

func (p UDPListenSpeaker) SetDeadline(t time.Time) (err error) {
	p.listener.SetDeadline(t)
	return nil
}

func (p UDPListenSpeaker) Accept() (net.Conn, error) {
	return p.listener, nil
}

func (p UDPListenSpeaker) Dial() (net.Conn, error) {
	return net.DialUDP(p.netType, nil, p.addr)
}

func (p UDPListenSpeaker) Close() {
	p.listener.Close()
}
