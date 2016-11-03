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

type GRegisterReactor struct {
	reactor.TCPReactor
	gateway *Gateway
}

func (p *GRegisterReactor) OnConnect(c *antgo.Conn)net.Addr {
	return nil
}

func (p *GRegisterReactor) OnMessage(c *antgo.Conn, pt antgo.Packet) bool{
	return true
}

func (p *GRegisterReactor) OnClose(c *antgo.Conn) {
}

type GEndReactor struct {
	reactor.TCPReactor
	gateway *Gateway
}

func (p *GEndReactor) OnConnect(c *antgo.Conn)net.Addr{
	return nil
}

func (p *GEndReactor) OnMessage(c *antgo.Conn, pt antgo.Packet) bool{
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

func (p *GWorkerReactor) OnConnect(c *antgo.Conn)net.Addr {
	return nil
}

func (p *GWorkerReactor) OnMessage(c *antgo.Conn, pt antgo.Packet) bool {
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
			p.gateway.AllGatewayAddr[addr] = nil
		}
	default:
		fmt.Println("Receive bad event:$event from Worker.\n")
	}
	return true
}

func (p *GWorkerReactor) OnClose(c *antgo.Conn) {
}

type Gateway struct {
	OuterAnt *antgo.Ant
	InnerAnt []*antgo.Ant

	EndConns     map[string]*antgo.Conn
	WorkerConns  map[string]*antgo.Conn
	RegisterConn *antgo.Conn

	AllGatewayAddr  map[string]net.Addr
	BusyGatewayAddr map[string]net.Addr
	IdleGatewayAddr map[string]net.Addr
}

func NewGateway(end_transport string, end_ip string, end_port int, end_lType string, end_pType string,
	worker_transport string, worker_ip string, worker_port []int, worker_lType string, worker_pType string,
	register_transport string, register_ip string, register_port int, register_lType string, register_pType string,
	sendLimit uint32, receiveLimit uint32) *Gateway {
	gateway := &Gateway{
		OuterAnt:     nil,
		InnerAnt:     make([]*antgo.Ant, 0, 12),
		WorkerConns:  make(map[string]*antgo.Conn),
		EndConns: make(map[string]*antgo.Conn),
		RegisterConn: nil,
		AllGatewayAddr:  make(map[string]net.Addr),
		BusyGatewayAddr: make(map[string]net.Addr),
		IdleGatewayAddr: make(map[string]net.Addr)}

	config := &antgo.Config{
		PacketSendChanLimit:    sendLimit,
		PacketReceiveChanLimit: receiveLimit}

	outProtocol := NewProtocol(end_pType, NewListenSpeaker(end_lType, end_transport, end_ip, end_port))
	outReactor := &GEndReactor{gateway:gateway}
	gateway.OuterAnt = antgo.NewAnt(end_transport, end_ip, end_port, config, outProtocol, outReactor)

	for _, port := range worker_port {
		innerProtocol := NewProtocol(worker_pType, NewListenSpeaker(worker_lType, worker_transport, worker_ip, port))
		innerReactor := &GWorkerReactor{gateway:gateway}
		gateway.InnerAnt = append(gateway.InnerAnt, antgo.NewAnt(worker_transport, worker_ip, port, config, innerProtocol, innerReactor))
	}
	return gateway
}

func (*Gateway) pingEnd() {

}

func (*Gateway) pingRegister() {

}

func (*Gateway) pingWorker() {

}

func (p *Gateway) Run() {
	go p.OuterAnt.Listen(time.Second)
	for _, innerAnt := range p.InnerAnt {
		go innerAnt.Listen(time.Second)
	}

	help := make(chan os.Signal)
	signal.Notify(help, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-help)
	p.OuterAnt.Stop()
	for _, innerAnt := range p.InnerAnt {
		innerAnt.Stop()
	}
}
