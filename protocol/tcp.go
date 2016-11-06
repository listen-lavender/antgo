package protocol

import (
	"../../antgo"
	"net"
	"strings"
)

var (
	endTag = []byte("\r\n") // TCP command's end tag
)

// Packet
type TCPPacket struct {
	event string
	msg   []byte
}

func NewTCPPacket(event string, msg []byte) *TCPPacket {
	return &TCPPacket{
		event: event,
		msg:   msg,
	}
}

func (p TCPPacket) Serialize() []byte {
	buf := p.msg
	buf = append(buf, endTag...)
	return buf
}

func (p TCPPacket) Event() string {
	return p.event
}

func (p TCPPacket) Msg() []byte {
	return p.msg
}

type TCPProtocol struct {
	// listendialer *multinet.TCPListenDialer
	listendialer antgo.ListenDialer
}

func NewTCPProtocol(listendialer antgo.ListenDialer) antgo.Protocol {
	return &TCPProtocol{listendialer}
}

func (p *TCPProtocol) ReadPacket(netConn net.Conn) antgo.Packet {
	listendialer := p.ListenDialer()
	command := listendialer.ReadPacket(netConn, endTag)
	parts := strings.Split(command, " ")
	if len(parts) > 1 {
		event := parts[0]
		msg := []byte(parts[1])
		return NewTCPPacket(event, msg)
	} else {
		if parts[0] == "quit" {
			return NewTCPPacket("quit", []byte(command))
		} else {
			return NewTCPPacket("unknow", []byte(command))
		}
	}
}

func (p *TCPProtocol) Deserialize(event string, msg []byte) antgo.Packet {
	return NewTCPPacket(event, msg)
}

func (p *TCPProtocol) ListenDialer() antgo.ListenDialer {
	return p.listendialer
}
