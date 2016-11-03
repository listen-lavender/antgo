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
	// listenspeaker *multinet.TCPListenSpeaker
	listenspeaker antgo.ListenSpeaker
}

func NewTCPProtocol(listenspeaker antgo.ListenSpeaker) antgo.Protocol {
	return &TCPProtocol{listenspeaker}
}

func (p TCPProtocol) ReadPacket(netConn net.Conn) <-chan antgo.Packet {
	queue := make(chan antgo.Packet)
	listenspeaker := p.ListenSpeaker()
	for command := range listenspeaker.ReadPacket(netConn, endTag) {
		parts := strings.Split(command, " ")
		if len(parts) > 1 {
			event := parts[0]
			msg := []byte(parts[1])
			queue <- NewTCPPacket(event, msg)
		} else {
			if parts[0] == "quit" {
				queue <- NewTCPPacket("quit", []byte(command))
			} else {
				queue <- NewTCPPacket("unknow", []byte(command))
			}
		}

	}
	return queue
}

func (p TCPProtocol) Deserialize(event string, msg []byte) antgo.Packet {
	return NewTCPPacket(event, msg)
}

func (p TCPProtocol) ListenSpeaker() antgo.ListenSpeaker {
	return p.listenspeaker
}
