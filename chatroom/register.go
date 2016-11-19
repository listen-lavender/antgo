package chatroom

import (
	"../../antgo"
	// "../../antgo/protocol"
	"../../antgo/reactor"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type RegisterReactor struct {
	register *Register
	reactor.TCPReactor
	secret       string
	WorkerConns  map[string]*antgo.Conn
	GatewayConns map[string]interface{}
}

func (p *RegisterReactor) OnConnect(c *antgo.Conn) string {
	addr := c.RemoteAddr()
	fmt.Println("OnConnect:", addr)
	c.PutExtraData(addr)
	p.register.Ant.Send(0, "prompt", "Welcome to Register Server", c, 0)
	return addr
}

func (p *RegisterReactor) OnMessage(c *antgo.Conn, pt antgo.Packet) bool {
	fmt.Println("OnMessage:", c.RemoteAddr())
	code := pt.Code()
	event := pt.Event()
	msg := pt.Msg()
	fmt.Println(code)
	switch event {
	case "prompt":
		fmt.Println(msg)
	case "gateway_connect":
		data := msg.(map[string]interface{})
		addresses := data["addresses"].([]interface{})
		for _, addr := range addresses {
			address := addr.(string)
			println(address)
			p.GatewayConns[address] = 1
			p.BroadcastAddr(address)
		}
	// 是 worker 连接
	case "worker_connect":
		data := msg.(map[string]interface{})
		secret := data["secret"]
		fmt.Println(secret)
		p.WorkerConns[c.RemoteAddr()] = c
		p.UnicastAddrs(c)
	case "ping":
		fmt.Println("ping")
	default:
		p.register.Ant.Send(0, "prompt", "unknow msg", c, 0)
		// p.OnClose(c)
		// c.Close()
	}
	return true
}

func (p *RegisterReactor) BroadcastAddr(addr string) {
	data := make(map[string]interface{})
	data["addresses"] = [1]string{addr}
	for _, c := range p.WorkerConns {
		p.register.Ant.Send(0, "broadcast_addresses", data, c, 0)
	}
}

func (p *RegisterReactor) UnicastAddrs(c *antgo.Conn) {
	data := make(map[string]interface{})
	data["addresses"] = antgo.MapKeys(p.GatewayConns)
	p.register.Ant.Send(0, "broadcast_addresses", data, c, 0)
}

type Register struct {
	*antgo.Ant
}

func NewRegister(localTransport string, localIP string, localPort int, localType string) *Register {
	// &Register{*antgo.NewAnt(localTransport, localIP, localPort, antgo.DefaultConfig, protocol, reactor)}
	register := &Register{Ant: nil}

	protocol := NewProtocol(localType, NewListenDialer(localType, localTransport, localIP, localPort))
	reactor := &RegisterReactor{
		register:     register,
		WorkerConns:  make(map[string]*antgo.Conn),
		GatewayConns: make(map[string]interface{}),
	}
	register.Ant = antgo.NewAnt(localTransport, localIP, localPort, antgo.DefaultConfig, protocol, reactor)
	return register
}

func (p *Register) Run() {
	go p.Listen(Timeout)
	help := make(chan os.Signal)
	signal.Notify(help, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-help)
	antgo.Stop()
}
