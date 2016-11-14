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
}

func (p *WRegisterReactor) OnConnect(c *antgo.Conn) string {
	addr := c.RemoteAddr()
	fmt.Println("OnConnect:", addr)
	p.worker.RegisterAnt.Send(0, "worker_connect", map[string]string{"secret": ""}, c, 0)
	return addr
}

func (p *WRegisterReactor) OnMessage(c *antgo.Conn, pt antgo.Packet) bool {
	code := pt.Code()
	event := pt.Event()
	msg := pt.Msg().(map[string]interface{})
	secret := msg["secret"]
	fmt.Println(code)
	fmt.Println(secret)
	switch event {
	case "broadcast_addresses":
		if msg["address"] == nil || msg["address"] == "" {
			fmt.Println("address not found\n")
			p.OnClose(c)
			c.Close()
		}
		addresses := msg["addresses"].([]string)
		for _, addr := range addresses {
			p.worker.AllGatewayAddr[addr] = nil
		}
	default:
		fmt.Println("Receive bad event:$event from Register.\n")
	}
	return true
}

type WGatewayReactor struct {
	reactor.TCPReactor
	WorkerConns  map[string]*antgo.Conn
	GatewayConns map[string]*antgo.Conn
	worker       *Worker
}

func (p *WGatewayReactor) OnMessage(c *antgo.Conn, pt antgo.Packet) bool {
	code := pt.Code()
	event := pt.Event()
	msg := pt.Msg().(map[string]interface{})
	secret := msg["secret"]
	fmt.Println(code)
	fmt.Println(secret)
	if event == "ping" {
		return true
	}
	Handlers[event](c, msg["data"].(string), p.worker)
	return true
}

type Worker struct {
	RegisterAnt *antgo.Ant
	GatewayAnt  []*antgo.Ant

	AllGatewayAddr  map[string]net.Addr
	BusyGatewayAddr map[string]net.Addr
	IdleGatewayAddr map[string]net.Addr
}

func NewWorker(register_transport string, register_ip string, register_port int, register_lType string, register_pType string) *Worker {
	worker := &Worker{
		RegisterAnt:     nil,
		GatewayAnt:      make([]*antgo.Ant, 0, 12),
		AllGatewayAddr:  make(map[string]net.Addr),
		BusyGatewayAddr: make(map[string]net.Addr),
		IdleGatewayAddr: make(map[string]net.Addr)}

	registerProtocol := NewProtocol(register_pType, NewListenDialer(register_lType, register_transport, register_ip, register_port))
	registerReactor := &WRegisterReactor{worker: worker}
	worker.RegisterAnt = antgo.NewAnt(register_transport, register_ip, register_port, antgo.DefaultConfig, registerProtocol, registerReactor)
	return worker
}

func (p *Worker) connectRegister() {
	p.RegisterAnt.Dial(Timeout)
}

func (p *Worker) connectGateway(gateway_transport string, gateway_ip string, gateway_port int, gateway_lType string, gateway_pType string) {
	gatewayProtocol := NewProtocol(gateway_pType, NewListenDialer(gateway_lType, gateway_transport, gateway_ip, gateway_port))
	gatewayReactor := &WGatewayReactor{worker: p}
	gatewayAnt := antgo.NewAnt(gateway_transport, gateway_ip, gateway_port, antgo.DefaultConfig, gatewayProtocol, gatewayReactor)
	go gatewayAnt.Dial(Timeout)
	gatewayAnt.Send(0, "worker_connect", []byte("Welcome to p TCP Server"), nil, 0)
	p.GatewayAnt = append(p.GatewayAnt, gatewayAnt)
}

func (p *Worker) PingRegister() {
	p.RegisterAnt.Send(0, "ping", []byte(""), nil, 0)
}

func (p *Worker) Run() {
	p.connectRegister()
	help := make(chan os.Signal)
	signal.Notify(help, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-help)
	p.RegisterAnt.Stop()
	// for _, gatewayAnt := range p.GatewayAnt {
	// 	gatewayAnt.Stop()
	// }
}
