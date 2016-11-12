package multinet

import (
	"../../antgo"
	"fmt"
	"net"
	"strconv"
	"time"
)

type UDPListenDialer struct {
	netType  string
	addr     *net.UDPAddr
	listener *net.UDPConn
}

func NewUDPListenDialer(netType string, ip string, port int) antgo.ListenDialer {
	addr, err := net.ResolveUDPAddr(netType, ip+":"+strconv.Itoa(port))
	if err != nil { //EOF, or worse
		fmt.Println(err)
		return nil
	}

	return &UDPListenDialer{
		netType:  netType,
		addr:     addr,
		listener: nil,
	}
}

func (p UDPListenDialer) ReadPacket(netConn net.Conn, endTag []byte) []byte {
	// msg := make(chan string)

	for {
		data := make([]byte, 1024)

		readLengh, err := netConn.Read(data)

		if err != nil { //EOF, or worse
			return nil
		}

		if readLengh > 0 {
			return data
		}
	}
	// return msg
}

func (p *UDPListenDialer) SetDeadline(t time.Time) (err error) {
	p.listener.SetDeadline(t)
	return nil
}

func (p *UDPListenDialer) Listen() {
	listener, err := net.ListenUDP(p.netType, p.addr)
	p.listener = listener
	if err != nil {
		fmt.Println(err)
	}
}

func (p *UDPListenDialer) Accept() (net.Conn, error) {
	return p.listener, nil
}

func (p *UDPListenDialer) Dial() (net.Conn, error) {
	return net.DialUDP(p.netType, nil, p.addr)
}

func (p *UDPListenDialer) Close() {
	p.listener.Close()
}
