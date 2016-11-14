package chatroom

import (
	"../../antgo"
	"../../antgo/protocol"
	"../../antgo/reactor"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type RegisterReactor struct {
	reactor.TCPReactor
	secret       string
	WorkerConns  map[string]*antgo.Conn
	GatewayConns map[string]interface{}
}

func (p *RegisterReactor) OnConnect(c *antgo.Conn) string {
	addr := c.RemoteAddr()
	fmt.Println("OnConnect:", addr)
	c.PutExtraData(addr)
	c.AsyncWritePacket(protocol.NewTCPPacket(0, "prompt", []byte("Welcome to p TCP Server")), 0)
	return addr
}

func (p *RegisterReactor) OnMessage(c *antgo.Conn, pt antgo.Packet) bool {
	// 删除定时器
	fmt.Println("OnMessage:", c.RemoteAddr())
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
		p.GatewayConns[c.RemoteAddr()] = 1
		p.BroadcastAddr(c.RemoteAddr())
		return true
	// 是 worker 连接
	case "worker_connect":
		p.WorkerConns[c.RemoteAddr()] = c
		p.UnicastAddrs(c)
		return true
	case "ping":
		return true
	default:
		c.AsyncWritePacket(protocol.NewTCPPacket(0, "prompt", "unknow msg"), 0)
		p.OnClose(c)
		c.Close()
		return true
	}
}

func (p *RegisterReactor) OnClose(c *antgo.Conn) {
	fmt.Println("OnClose:", c.GetExtraData())
}

func (p *RegisterReactor) BroadcastAddr(add string) {
	data := make(map[string]interface{})
	data["addresses"] = [1]string{add}
	for _, c := range p.WorkerConns {
		c.AsyncWritePacket(protocol.NewTCPPacket(0, "broadcast_addresses", data), 0)
	}
}

func (p *RegisterReactor) UnicastAddrs(c *antgo.Conn) {
	data := make(map[string]interface{})
	data["addresses"] = antgo.MapKeys(p.GatewayConns)
	c.AsyncWritePacket(protocol.NewTCPPacket(0, "broadcast_addresses", data), 0)
}

type Register struct {
	antgo.Ant
}

func NewRegister(transport string, ip string, port int, lType string, pType string) *Register {
	protocol := NewProtocol(pType, NewListenDialer(lType, transport, ip, port))
	reactor := &RegisterReactor{
		WorkerConns:  make(map[string]*antgo.Conn),
		GatewayConns: make(map[string]interface{}),
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
