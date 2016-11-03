package chatroom

import (
	"../../antgo"
	"../../antgo/reactor"
	"fmt"
	"net"
    "os"
    "os/signal"
    "time"
    "syscall"
)

type WRegisterReactor struct {
	reactor.TCPReactor
	worker       *Worker
}

func (p *WRegisterReactor) OnMessage(c *antgo.Conn, pt antgo.Packet) bool {
	event := pt.Event()
	msg := antgo.JsonDecode(pt.Msg())
	secret := msg["secret"]
	fmt.Println(secret)
	switch event {
	case "broadcast_addresses":
		if msg["address"] == nil || msg["address"] == "" {
			fmt.Println("address not found\n")
			c.Close()
			p.OnClose(c)
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
	event := pt.Event()
	msg := antgo.JsonDecode(pt.Msg())
	secret := msg["secret"]
	fmt.Println(secret)
	if event == "ping" {
		return true
	}
	Handlers[event](c, msg["data"].(string), p.worker)
    return true
}

// type WorkerReactor struct {
//     reactor.Reactor

// }

type Worker struct {
	RegisterAnt *antgo.Ant
	GatewayAnt []*antgo.Ant
	
	AllGatewayAddr  map[string]net.Addr
	BusyGatewayAddr map[string]net.Addr
	IdleGatewayAddr map[string]net.Addr
}

func NewWorker(register_transport string, register_ip string, register_port int, register_lType string, register_pType string,
	sendLimit uint32, receiveLimit uint32) *Worker {
	worker := &Worker{
        RegisterAnt:     nil,
        GatewayAnt:     make([]*antgo.Ant, 0, 12),
        AllGatewayAddr:  make(map[string]net.Addr),
        BusyGatewayAddr: make(map[string]net.Addr),
        IdleGatewayAddr: make(map[string]net.Addr)}

    config := &antgo.Config{
        PacketSendChanLimit:    sendLimit,
        PacketReceiveChanLimit: receiveLimit}

    registerProtocol := NewProtocol(register_pType, NewListenSpeaker(register_lType, register_transport, register_ip, register_port))
    registerReactor := &WRegisterReactor{worker:worker}
    worker.RegisterAnt = antgo.NewAnt(register_transport, register_ip, register_port, config, registerProtocol, registerReactor)
    return worker
}

func (p *Worker) connectRegister() {
	go p.RegisterAnt.Speak(time.Second)
	p.RegisterAnt.Send("worker_connect", []byte("Welcome to p TCP Server"), 0)
}

func (p *Worker) connectGateway(gateway_transport string, gateway_ip string, gateway_port int, gateway_lType string, gateway_pType string) {
	config := &antgo.Config{
    PacketSendChanLimit:    20,
    PacketReceiveChanLimit: 20}
	gatewayProtocol := NewProtocol(gateway_pType, NewListenSpeaker(gateway_lType, gateway_transport, gateway_ip, gateway_port))
	gatewayReactor := &WGatewayReactor{worker:p}
	gatewayAnt := antgo.NewAnt(gateway_transport, gateway_ip, gateway_port, config, gatewayProtocol, gatewayReactor)
	go gatewayAnt.Speak(time.Second)
	gatewayAnt.Send("worker_connect", []byte("Welcome to p TCP Server"), 0)
	p.GatewayAnt = append(p.GatewayAnt, gatewayAnt)
}

func (p *Worker) PingRegister() {
    p.RegisterAnt.Send("ping", []byte(""), 0)
}

func (p *Worker) Run() {
    p.connectRegister()
    help := make(chan os.Signal)
    signal.Notify(help, syscall.SIGINT, syscall.SIGTERM)
    fmt.Println("Signal: ", <-help)
    p.RegisterAnt.Stop()
    for _, gatewayAnt := range(p.GatewayAnt){
        gatewayAnt.Stop()
    }
}

