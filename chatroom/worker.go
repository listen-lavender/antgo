package chatroom

import (
	"../../antgo"
	"../../antgo/reactor"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type WRegisterReactor struct {
	reactor.TCPReactor
	worker *Worker
	ant    *antgo.Ant
}

func (p *WRegisterReactor) OnConnect(c *antgo.Conn) string {
	addr := c.RemoteAddr()
	fmt.Println("OnConnect:", addr)
	p.ant.Send(0, "worker_connect", map[string]string{"secret": ""}, c, 0)
	return addr
}

func (p *WRegisterReactor) OnMessage(c *antgo.Conn, pt antgo.Packet) bool {
	fmt.Println("OnMessage:", c.RemoteAddr())
	code := pt.Code()
	event := pt.Event()
	msg := pt.Msg()
	fmt.Println(code)
	switch event {
	case "prompt":
		fmt.Println("prompt")
		fmt.Println(msg)
	case "broadcast_addresses":
		data := msg.(map[string]interface{})
		secret := data["secret"]
		fmt.Println("broadcast_addresses")
		fmt.Println(secret)
		addresses := data["addresses"].([]interface{})
		for _, addr := range addresses {
			address := addr.(string)
			gatewayTransport, gatewayIP, gatewayPort, gatewayType := antgo.AddressSplit(address)
			p.worker.connectGateway(gatewayTransport, gatewayIP, gatewayPort, gatewayType)
			p.worker.AllGatewayAddr[address] = nil
			fmt.Println(address)
		}
	case "ping":
		fmt.Println("ping")
	default:
		fmt.Println("Receive bad event:$event from Register.\n")
	}
	return true
}

type WGatewayReactor struct {
	reactor.TCPReactor
	// WorkerConns  map[string]*antgo.Conn
	// GatewayConns map[string]*antgo.Conn
	worker *Worker
	ant    *antgo.Ant
}

func (p *WGatewayReactor) OnConnect(c *antgo.Conn) string {
	addr := c.RemoteAddr()
	fmt.Println("OnConnect:", addr)
	p.ant.Send(0, "worker_connect", map[string]string{"secret": "abc"}, c, 0)
	return addr
}

func (p *WGatewayReactor) OnMessage(c *antgo.Conn, pt antgo.Packet) bool {
	fmt.Println("OnMessage:", c.RemoteAddr())
	code := pt.Code()
	event := pt.Event()
	msg := pt.Msg()
	fmt.Println(code)
	switch event {
	case "prompt":
		fmt.Println("prompt")
		fmt.Println(msg)
	case "ping":
		fmt.Println("ping")
	default:
		fmt.Println("Receive bad event:$event from Register.\n")
	}
	return true
}

type Worker struct {
	RegisterAnt *antgo.Ant
	GatewayAnt  []*antgo.Ant

	AllGatewayAddr  map[string]net.Addr
	BusyGatewayAddr map[string]net.Addr
	IdleGatewayAddr map[string]net.Addr
}

func NewWorker(registerTransport string, registerIP string, registerPort int, registerType string) *Worker {
	worker := &Worker{
		RegisterAnt:     nil,
		GatewayAnt:      make([]*antgo.Ant, 0, 12),
		AllGatewayAddr:  make(map[string]net.Addr),
		BusyGatewayAddr: make(map[string]net.Addr),
		IdleGatewayAddr: make(map[string]net.Addr)}

	registerProtocol := NewProtocol(registerType, NewListenDialer(registerType, registerTransport, registerIP, registerPort))
	registerReactor := &WRegisterReactor{worker: worker}
	registerAnt := antgo.NewAnt(registerTransport, registerIP, registerPort, antgo.DefaultConfig, registerProtocol, registerReactor)
	registerReactor.ant = registerAnt
	worker.RegisterAnt = registerAnt
	return worker
}

func (p *Worker) connectRegister() {
	go p.RegisterAnt.Dial(Timeout)
}

func (p *Worker) connectGateway(gatewayTransport string, gatewayIP string, gatewayPort int, gatewayType string) {
	gatewayProtocol := NewProtocol(gatewayType, NewListenDialer(gatewayType, gatewayTransport, gatewayIP, gatewayPort))
	gatewayReactor := &WGatewayReactor{worker: p}
	gatewayAnt := antgo.NewAnt(gatewayTransport, gatewayIP, gatewayPort, antgo.DefaultConfig, gatewayProtocol, gatewayReactor)
	gatewayReactor.ant = gatewayAnt
	go gatewayAnt.Dial(Timeout)
	p.GatewayAnt = append(p.GatewayAnt, gatewayAnt)
}

func (p *Worker) PingRegister() {
	p.RegisterAnt.Send(0, "ping", "", nil, 0)
}

func (p *Worker) Run() {
	p.connectRegister()
	help := make(chan os.Signal)
	signal.Notify(help, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-help)
	antgo.Stop()
}
