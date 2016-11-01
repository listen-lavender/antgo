package chatroom

import (
	"../../antgo"
	// "../../antgo/multinet"
	"../../antgo/protocol"
	"../../antgo/reactor"
	"fmt"
	"net"
)

type RegisterReactor struct {
	reactor.TCPReactor
	secret       string
	WorkerConns  map[string]*antgo.Conn
	GatewayConns map[string]net.Addr
}

func (p *RegisterReactor) OnConnect(c *antgo.Conn) net.Addr {
	addr := c.RemoteAddr()
	fmt.Println("OnConnect:", addr)
	c.PutExtraData(addr)
	c.AsyncWritePacket(protocol.NewTCPPacket("prompt", []byte("Welcome to p TCP Server")), 0)
	return addr
}

func (p *RegisterReactor) OnMessage(c *antgo.Conn, pt antgo.Packet) bool {
	// 删除定时器
	event := pt.Event()
	msg := antgo.JsonDecode(pt.Msg())
	secret := msg["secret"]
	fmt.Println(secret)
	// 开始验证
	switch event {
	// 是 gateway 连接
	case "gateway_connect":
		if msg["address"] == nil || msg["address"] == "" {
			fmt.Println("address not found\n")
			c.Close()
			p.OnClose(c)
			return true
		}
		if msg["secret"] != p.secret {
			fmt.Println("Register: Key does not match secret_key !== {this->secretKey}\n")
			c.Close()
			p.OnClose(c)
			return true
		}
		p.GatewayConns[c.Id] = c.RemoteAddr() // msg["address"]
		p.BroadcastAddrs()
		return true
	// 是 worker 连接
	case "worker_connect":
		if msg["secret"] != p.secret {
			fmt.Println("Register: Key does not match secret_key !== {this->secretKey}\n")
			c.Close()
			p.OnClose(c)
			return true
		}
		p.WorkerConns[c.Id] = c
		p.UnicastAddrs(c)
		return true
	case "ping":
		return true
	default:
		c.AsyncWritePacket(protocol.NewTCPPacket("prompt", []byte("unknow msg")), 0)
		c.Close()
		p.OnClose(c)
		return true
	}
}

func (p *RegisterReactor) OnClose(c *antgo.Conn) {
	fmt.Println("OnClose:", c.GetExtraData())
}

func (p *RegisterReactor) BroadcastAddrs() {
	buffer := antgo.JsonEncode(p.GatewayConns)
	for _, c := range p.WorkerConns {
		c.AsyncWritePacket(protocol.NewTCPPacket("broadcast_addresses", buffer), 0)
	}
}

func (p *RegisterReactor) UnicastAddrs(c *antgo.Conn) {
	buffer := antgo.JsonEncode(p.GatewayConns)
	c.AsyncWritePacket(protocol.NewTCPPacket("broadcast_addresses", buffer), 0)
}

type Register struct {
	antgo.Ant
}

func NewRegister(transport string, ip string, port int, lType string, pType string, sendLimit uint32, receiveLimit uint32) *Register {
	config := &antgo.Config{
		PacketSendChanLimit:    sendLimit,
		PacketReceiveChanLimit: receiveLimit}

	protocol := NewProtocol(pType, NewListenSpeaker(lType, transport, ip, port))
	reactor := &RegisterReactor{
		WorkerConns:  make(map[string]*antgo.Conn),
		GatewayConns: make(map[string]net.Addr),
	}
	return &Register{
		*antgo.NewAnt(transport, ip, port, config, protocol, reactor)}
}
