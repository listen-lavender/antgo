package chatroom

import (
	"../../antgo"
	"../../antgo/protocol"
	"../../antgo/reactor"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
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
	c.AsyncWritePacket(protocol.NewTCPPacket(0, "prompt", []byte("Welcome to p TCP Server")), 0)
	return addr
}

func (p *RegisterReactor) OnMessage(c *antgo.Conn, pt antgo.Packet) bool {
	// 删除定时器
	code := pt.Code()
	event := pt.Event()
	msg := pt.Msg().(map[string]interface{})
	secret := msg["secret"]
	fmt.Println(secret)
	fmt.Println(code)
	// 开始验证
	switch event {
	// 是 gateway 连接
	case "gateway_connect":
		if msg["address"] == nil || msg["address"] == "" {
			fmt.Println("address not found\n")
			p.OnClose(c)
			c.Close()
			return true
		}
		if msg["secret"] != p.secret {
			fmt.Println("Register: Key does not match secret_key !== {this->secretKey}\n")
			p.OnClose(c)
			c.Close()
			return true
		}
		p.GatewayConns[c.Id] = c.RemoteAddr() // msg["address"]
		p.BroadcastAddrs()
		return true
	// 是 worker 连接
	case "worker_connect":
		if msg["secret"] != p.secret {
			fmt.Println("Register: Key does not match secret_key !== {this->secretKey}\n")
			p.OnClose(c)
			c.Close()
			return true
		}
		p.WorkerConns[c.Id] = c
		p.UnicastAddrs(c)
		return true
	case "ping":
		return true
	default:
		c.AsyncWritePacket(protocol.NewTCPPacket(0, "prompt", []byte("unknow msg")), 0)
		p.OnClose(c)
		c.Close()
		return true
	}
}

func (p *RegisterReactor) OnClose(c *antgo.Conn) {
	fmt.Println("OnClose:", c.GetExtraData())
}

func (p *RegisterReactor) BroadcastAddrs() {
	buffer := antgo.JsonEncode(p.GatewayConns)
	for _, c := range p.WorkerConns {
		c.AsyncWritePacket(protocol.NewTCPPacket(0, "broadcast_addresses", buffer), 0)
	}
}

func (p *RegisterReactor) UnicastAddrs(c *antgo.Conn) {
	buffer := antgo.JsonEncode(p.GatewayConns)
	c.AsyncWritePacket(protocol.NewTCPPacket(0, "broadcast_addresses", buffer), 0)
}

type Register struct {
	antgo.Ant
}

func NewRegister(transport string, ip string, port int, lType string, pType string) *Register {
	protocol := NewProtocol(pType, NewListenDialer(lType, transport, ip, port))
	reactor := &RegisterReactor{
		WorkerConns:  make(map[string]*antgo.Conn),
		GatewayConns: make(map[string]net.Addr),
	}
	return &Register{
		*antgo.NewAnt(transport, ip, port, antgo.DefaultConfig, protocol, reactor)}
}

func (p *Register) Run() {
	go p.Listen(Timeout)
	help := make(chan os.Signal)
	signal.Notify(help, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-help)
	p.Stop()
}
