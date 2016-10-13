package protocol

import (
	"../../antgo"
	"../../antgo/listener"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// Packet
type EchoPacket struct {
	buff []byte
}

func NewEchoPacket(buff []byte, hasLengthField bool) *EchoPacket {
	p := &EchoPacket{}

	if hasLengthField {
		p.buff = buff

	} else {
		p.buff = make([]byte, 4+len(buff))
		binary.BigEndian.PutUint32(p.buff[0:4], uint32(len(buff)))
		copy(p.buff[4:], buff)
	}
	return p
}

func (p EchoPacket) Serialize() []byte {
	return p.buff
}

func (p EchoPacket) GetLength() uint32 {
	return binary.BigEndian.Uint32(p.buff[0:4])
}

func (p EchoPacket) GetBody() []byte {
	return p.buff[4:]
}

type EchoProtocol struct {
	listener *listener.TCPListener
}

func NewEchoProtocol(listener *listener.TCPListener) *EchoProtocol {
	return &EchoProtocol{listener}
}

func (p EchoProtocol) ReadPacket(conn net.Conn) <-chan antgo.Packet {
	msg := make(chan antgo.Packet)

	var (
		lengthBytes []byte = make([]byte, 4)
		length      uint32
	)

	// read length
	if _, err := io.ReadFull(conn, lengthBytes); err != nil {
		fmt.Println(err)
		return nil
	}
	if length = binary.BigEndian.Uint32(lengthBytes); length > 1024 {
		fmt.Println("the size of packet is larger than the limit")
		return nil
	}

	buff := make([]byte, 4+length)
	copy(buff[0:4], lengthBytes)

	// read body ( buff = lengthBytes + body )
	if _, err := io.ReadFull(conn, buff[4:]); err != nil {
		fmt.Println(err)
		return nil
	}

	msg <- NewEchoPacket(buff, true)
	return msg
}

func (p EchoProtocol) GetListener() antgo.Listener {
	return p.listener
}

// func (p EchoProtocol)SetListener(listener antgo.Listener){
//     p.listener = listener
// }
