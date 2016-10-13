package protocol

import (
    "net"
    "strings"
    "../../antgo/listener"
    "../../antgo"
)

var (
    endTag = []byte("\r\n") // TCP command's end tag
)

// Packet
type TCPPacket struct {
    event string
    message []byte
}

func NewTCPPacket(event string, message []byte) *TCPPacket {
    return &TCPPacket{
        event: event,
        message: message,
    }
}

func (p TCPPacket) Serialize() []byte {
    buf := p.message
    buf = append(buf, endTag...)
    return buf
}

func (p TCPPacket) GetType() string {
    return p.event
}

func (p TCPPacket) GetData() []byte {
    return p.message
}

type TCPProtocol struct {
    listener *listener.TCPListener
}

func NewTCPProtocol(listener *listener.TCPListener) *TCPProtocol {
    return &TCPProtocol{listener}
}

func (p TCPProtocol) ReadPacket(conn net.Conn) (<-chan antgo.Packet) {
    msg := make(chan antgo.Packet)
    listener := p.GetListener()
    for command := range(listener.ReadPacket(conn, endTag)){
        parts := strings.Split(command, " ")
        if len(parts) > 1 {
            event := parts[0]
            message := []byte(parts[1])
            msg <- NewTCPPacket(event, message)
        } else {
            if parts[0] == "quit" {
                msg <- NewTCPPacket("quit", []byte(command))
            } else {
                msg <- NewTCPPacket("unknow", []byte(command))
            }
        }

    }
    return msg
}

func (p TCPProtocol)GetListener() antgo.Listener{
    return p.listener
}

// func (p TCPProtocol)SetListener(listener antgo.Listener){
//     p.listener = listener
// }

