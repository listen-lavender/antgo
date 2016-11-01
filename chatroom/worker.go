package chatroom

import (
	"fmt"
	"../../antgo"
	"../../antgo/protocol"
	"../../antgo/reactor"
)

type WRegisterReactor struct {
	reactor.TCPReactor
	allGatewayAddr  map[string]string
	busyGatewayAddr map[string]string
	idleGatewayAddr map[string]string
}

func (p *WRegisterReactor) OnMessage(c *antgo.Conn, pt antgo.Packet) {
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
			p.allGatewayAddr[addr] = addr
		}
	default:
		fmt.Println("Receive bad event:$event from Register.\n")
	}
}

type WGatewayReactor struct {
	reactor.TCPReactor
	WorkerConns  map[string]*antgo.Conn
	GatewayConns map[string]*antgo.Conn
	worker *Worker
}

func (p *WGatewayReactor) OnMessage(c *antgo.Conn, pt antgo.Packet) {
	event := pt.Event()
	msg := antgo.JsonDecode(pt.Msg())
	secret := msg["secret"]
	fmt.Println(secret)
	if event == "ping" {
		return
	}
	Handlers[event](c, msg["data"].(string), p.worker)
}

// type WorkerReactor struct {
//     reactor.Reactor

// }

type Worker struct {
	RegisterConn *antgo.Conn
	GatewayConns map[string]*antgo.Conn
}

func NewWorker(transport string, ip string, port int, lType string, pType string, rType string, sendLimit uint32, receiveLimit uint32) *Worker {
	// config := &antgo.Config{
	// 	PacketSendChanLimit:    sendLimit,
	// 	PacketReceiveChanLimit: receiveLimit}

    // protocol := NewProtocol(pType, NewListenSpeaker(lType, transport, ip, port))
    // reactor := &reactor.TCPReactor{}

	return &Worker{
		RegisterConn: nil,
		GatewayConns: make(map[string]*antgo.Conn)}
}

func (p *Worker) connectRegister() {
	// ant := &antgo.Ant{
	// 	Transport: transport,
	// 	IP:        ip,
	// 	Port:      port,

	// 	config:   config,
	// 	protocol: NewProtocol(pType, NewListenSpeaker(lType, transport, ip, port)),
	// 	reactor:  &WRegisterReactor{},

	// 	exitChan:  make(chan struct{}),
	// 	waitGroup: &sync.WaitGroup{},
	// }
	// conn := ant.Speak()
	// p.RegisterConn = conn
	// conn.AsyncWritePacket(protocol.NewTCPPacket("worker_connect", []byte("Welcome to p TCP Server")), 0)
}

func (p *Worker) connectGateway(addr string) {
	// ant := &antgo.Ant{
	// 	Transport: transport,
	// 	IP:        ip,
	// 	Port:      port,

	// 	config:   config,
	// 	protocol: NewProtocol(pType, NewListenSpeaker(lType, transport, ip, port)),
	// 	reactor:  &WGatewayReactor{},

	// 	exitChan:  make(chan struct{}),
	// 	waitGroup: &sync.WaitGroup{}}
	// conn := ant.Speak()
	// p.GatewayConns = append(p.GatewayConns, ant.Speak())
	// conn.AsyncWritePacket(protocol.NewTCPPacket("worker_connect", []byte("Welcome to p TCP Server")), 0)
}

func (p *Worker) PingRegister() {
	p.RegisterConn.AsyncWritePacket(protocol.NewTCPPacket("ping", []byte("")), 0)
}
