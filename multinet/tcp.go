package multinet

import (
	"../../antgo"
	"bytes"
	"fmt"
	"net"
	"strconv"
	"time"
)

type TCPListenSpeaker struct {
	netType  string
	addr     *net.TCPAddr
	listener *net.TCPListener
}

func NewTCPListenSpeaker(netType string, ip string, port int) antgo.ListenSpeaker {
	addr, err := net.ResolveTCPAddr(netType, ip+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	listener, err := net.ListenTCP(netType, addr)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &TCPListenSpeaker{
		netType:  netType,
		addr:     addr,
		listener: listener,
	}
}

func (p TCPListenSpeaker) ReadPacket(netConn net.Conn, endTag []byte) <-chan string {
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

func (p TCPListenSpeaker) SetDeadline(t time.Time) (err error) {
	p.listener.SetDeadline(t)
	return nil
}

func (p TCPListenSpeaker) Accept() (net.Conn, error) {
	return p.listener.AcceptTCP()
}

func (p TCPListenSpeaker) Dial() (net.Conn, error) {
	return net.DialTCP(p.netType, nil, p.addr)
}

func (p TCPListenSpeaker) Close() {
	p.listener.Close()
}
