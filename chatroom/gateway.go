package chatroom

import (
	"../../antgo"
	"../../antgo/reactor"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type GRegisterReactor struct {
	reactor.TCPReactor
	gateway *Gateway
}

func (p *GRegisterReactor) OnConnect(c *antgo.Conn) string {
	return ""
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

func (p *GEndReactor) OnConnect(c *antgo.Conn) string {
	return ""
}

func (p *GEndReactor) OnMessage(c *antgo.Conn, pt antgo.Packet) bool {
	code := pt.Code()
	event := pt.Event()
	msg := pt.Msg().(map[string]interface{})
	secret := msg["secret"]
	fmt.Println(code)
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

func (p *GWorkerReactor) OnConnect(c *antgo.Conn) string {
	return ""
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
	register_transport string, register_ip string, register_port int, register_lType string, register_pType string) *Gateway {
	gateway := &Gateway{
		EndAnt:      nil,
		RegisterAnt: nil,
		WorkerAnt:   make([]*antgo.Ant, 0, 12)}

	endProtocol := NewProtocol(end_pType, NewListenDialer(end_lType, end_transport, end_ip, end_port))
	endReactor := &GEndReactor{gateway: gateway}
	gateway.EndAnt = antgo.NewAnt(end_transport, end_ip, end_port, antgo.DefaultConfig, endProtocol, endReactor)

	registerProtocol := NewProtocol(register_pType, NewListenDialer(register_lType, register_transport, register_ip, register_port))
	registerReactor := &GRegisterReactor{gateway: gateway}
	gateway.RegisterAnt = antgo.NewAnt(register_transport, register_ip, register_port, antgo.DefaultConfig, registerProtocol, registerReactor)

	for _, port := range worker_port {
		workerProtocol := NewProtocol(worker_pType, NewListenDialer(worker_lType, worker_transport, worker_ip, port))
		workerReactor := &GWorkerReactor{gateway: gateway}
		gateway.WorkerAnt = append(gateway.WorkerAnt, antgo.NewAnt(worker_transport, worker_ip, port, antgo.DefaultConfig, workerProtocol, workerReactor))
	}
	return gateway
}

func (p *Gateway) connectRegister() {
	go p.RegisterAnt.Dial(Timeout)
	p.RegisterAnt.Send(0, "gateway_connect", []byte("Welcome to p TCP Server"), nil, 0)
}

func (p *Gateway) pingEnd() {
	p.EndAnt.Send(0, "ping", []byte("Welcome to p TCP Server"), nil, 0)
}

func (p *Gateway) pingRegister() {
	p.RegisterAnt.Send(0, "ping", []byte("Welcome to p TCP Server"), nil, 0)
}

func (p *Gateway) pingWorker() {
	for _, WorkerAnt := range p.WorkerAnt {
		WorkerAnt.Send(0, "ping", []byte("Welcome to p TCP Server"), nil, 0)
	}
}

func (p *Gateway) Run() {
	go p.EndAnt.Listen(Timeout)
	for _, WorkerAnt := range p.WorkerAnt {
		go WorkerAnt.Listen(Timeout)
	}
	p.connectRegister()

	help := make(chan os.Signal)
	signal.Notify(help, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-help)
	p.EndAnt.Stop()
	for _, WorkerAnt := range p.WorkerAnt {
		WorkerAnt.Stop()
	}
	p.RegisterAnt.Stop()
}
