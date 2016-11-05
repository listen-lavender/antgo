package multinet

import (
	"../../antgo"
	"bytes"
	"fmt"
	"net"
	"strconv"
	"time"
)

type TCPListenDialer struct {
	netType  string
	addr     *net.TCPAddr
	listener *net.TCPListener
}

func NewTCPListenDialer(netType string, ip string, port int) antgo.ListenDialer {
	addr, err := net.ResolveTCPAddr(netType, ip+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &TCPListenDialer{
		netType:  netType,
		addr:     addr,
		listener: nil,
	}
}

func (p *TCPListenDialer) ReadPacket(netConn net.Conn, endTag []byte) <-chan string {
	fullBuf := bytes.NewBuffer([]byte{})
	msg := make(chan string)

	for {
		data := make([]byte, 1024)

		readLengh, err := netConn.Read(data)

		if err != nil { //EOF, or worse
			fmt.Println(err)
			return nil
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

func (p *TCPListenDialer) SetDeadline(t time.Time) (err error) {
	p.listener.SetDeadline(t)
	return nil
}

func (p *TCPListenDialer) Listen() {
	listener, err := net.ListenTCP(p.netType, p.addr)
	p.listener = listener
	if err != nil {
		fmt.Println(err)
	}
}

func (p *TCPListenDialer) Accept() (net.Conn, error) {
	return p.listener.AcceptTCP()
}

func (p *TCPListenDialer) Dial() (net.Conn, error) {
	return net.DialTCP(p.netType, nil, p.addr)
}

func (p *TCPListenDialer) Close() {
	p.listener.Close()
}
