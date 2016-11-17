package protocol

import (
	"../../antgo"
	"net"
	// "strings"
)

var (
	endTag = []byte("\r\n") // TCP command's end tag
)

// Packet
type TCPPacket struct {
	code  int
	event string
	msg   interface{}
}

func NewTCPPacket(code int, event string, msg interface{}) *TCPPacket {
	return &TCPPacket{
		code:  code,
		event: event,
		msg:   msg,
	}
}

func (p TCPPacket) Serialize() []byte {
	data := make(map[string]interface{})
	data["code"] = p.code
	data["event"] = p.event
	data["msg"] = p.msg
	buf := antgo.JsonEncode(data)
	buf = append(buf, endTag...)
	return buf
}

func (p TCPPacket) Code() int {
	return p.code
}

func (p TCPPacket) Event() string {
	return p.event
}

func (p TCPPacket) Msg() interface{} {
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
	buf := listendialer.ReadPacket(netConn, endTag)
	if buf == nil{
		return nil
	}
	data := antgo.JsonDecode(buf)
	// parts := strings.Split(command, " ")
	code, _ := data["code"].(int)
	event, _ := data["event"].(string)
	msg := data["msg"]
	return NewTCPPacket(code, event, msg)
	// if len(parts) > 1 {
	// 	event := parts[0]
	// 	msg := []byte(parts[1])
	// 	return NewTCPPacket(event, msg)
	// } else {
	// 	if parts[0] == "quit" {
	// 		return NewTCPPacket("quit", []byte(command))
	// 	} else {
	// 		return NewTCPPacket("unknow", []byte(command))
	// 	}
	// }
}

func (p *TCPProtocol) Deserialize(code int, event string, msg interface{}) antgo.Packet {
	return NewTCPPacket(code, event, msg)
}

func (p *TCPProtocol) ListenDialer() antgo.ListenDialer {
	return p.listendialer
}
