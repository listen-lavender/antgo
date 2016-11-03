package chatroom

import (
	"../../antgo"
	"../../antgo/reactor"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type GRegisterReactor struct {
	reactor.TCPReactor
	gateway *Gateway
}

func (p *GRegisterReactor) OnConnect(c *antgo.Conn) net.Addr {
	return nil
}

func (p *GRegisterReactor) OnMessage(c *antgo.Conn, pt antgo.Packet) bool {
	return true
}

func (p *GRegisterReactor) OnClose(c *antgo.Conn) {
}

type GEndReactor struct {
	reactor.TCPReactor
	gateway *Gateway
}

func (p *GEndReactor) OnConnect(c *antgo.Conn) net.Addr {
	return nil
}

func (p *GEndReactor) OnMessage(c *antgo.Conn, pt antgo.Packet) bool {
	event := pt.Event()
	msg := antgo.JsonDecode(pt.Msg())
	secret := msg["secret"]
	fmt.Println(secret)
	fmt.Println(event)
	return true
}

func (p *GEndReactor) OnClose(c *antgo.Conn) {
}

type GWorkerReactor struct {
	reactor.TCPReactor
	gateway *Gateway
}

func (p *GWorkerReactor) OnConnect(c *antgo.Conn) net.Addr {
	return nil
}

func (p *GWorkerReactor) OnMessage(c *antgo.Conn, pt antgo.Packet) bool {
	return true
}

func (p *GWorkerReactor) OnClose(c *antgo.Conn) {
}

type Gateway struct {
	EndAnt      *antgo.Ant
	WorkerAnt   []*antgo.Ant
	RegisterAnt *antgo.Ant

	// EndConns     map[string]*antgo.Conn
	// WorkerConns  map[string]*antgo.Conn
	// RegisterConn *antgo.Conn
}

func NewGateway(end_transport string, end_ip string, end_port int, end_lType string, end_pType string,
	worker_transport string, worker_ip string, worker_port []int, worker_lType string, worker_pType string,
	register_transport string, register_ip string, register_port int, register_lType string, register_pType string,
	sendLimit uint32, receiveLimit uint32) *Gateway {
	gateway := &Gateway{
		EndAnt:      nil,
		RegisterAnt: nil,
		WorkerAnt:   make([]*antgo.Ant, 0, 12)}

	config := &antgo.Config{
		PacketSendChanLimit:    sendLimit,
		PacketReceiveChanLimit: receiveLimit}

	endProtocol := NewProtocol(end_pType, NewListenSpeaker(end_lType, end_transport, end_ip, end_port))
	endReactor := &GEndReactor{gateway: gateway}
	gateway.EndAnt = antgo.NewAnt(end_transport, end_ip, end_port, config, endProtocol, endReactor)

	registerProtocol := NewProtocol(register_pType, NewListenSpeaker(register_lType, register_transport, register_ip, register_port))
	registerReactor := &GRegisterReactor{gateway: gateway}
	gateway.RegisterAnt = antgo.NewAnt(register_transport, register_ip, register_port, config, registerProtocol, registerReactor)

	for _, port := range worker_port {
		workerProtocol := NewProtocol(worker_pType, NewListenSpeaker(worker_lType, worker_transport, worker_ip, port))
		workerReactor := &GWorkerReactor{gateway: gateway}
		gateway.WorkerAnt = append(gateway.WorkerAnt, antgo.NewAnt(worker_transport, worker_ip, port, config, workerProtocol, workerReactor))
	}
	return gateway
}

func (p *Gateway) connectRegister() {
	go p.RegisterAnt.Speak(time.Second)
	p.RegisterAnt.Send("gateway_connect", []byte("Welcome to p TCP Server"), 0)
}

func (p *Gateway) pingEnd() {
	p.EndAnt.Send("ping", []byte("Welcome to p TCP Server"), 0)
}

func (p *Gateway) pingRegister() {
	p.RegisterAnt.Send("ping", []byte("Welcome to p TCP Server"), 0)
}

func (p *Gateway) pingWorker() {
	for _, WorkerAnt := range p.WorkerAnt {
		WorkerAnt.Send("ping", []byte("Welcome to p TCP Server"), 0)
	}
}

func (p *Gateway) Run() {
	go p.EndAnt.Listen(time.Second)
	for _, WorkerAnt := range p.WorkerAnt {
		go WorkerAnt.Listen(time.Second)
	}
	// p.RegisterAnt.connectRegister()

	help := make(chan os.Signal)
	signal.Notify(help, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-help)
	p.EndAnt.Stop()
	for _, WorkerAnt := range p.WorkerAnt {
		WorkerAnt.Stop()
	}
	// p.RegisterAnt.Stop()
}
