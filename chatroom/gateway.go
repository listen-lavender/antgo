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
	addr := c.RemoteAddr()
	fmt.Println("OnConnect:", addr)
	p.gateway.RegisterAnt.Send(0, "gateway_connect", map[string]interface{}{"secret": "", "addresses":p.gateway.Addresses}, c, 0)
	return addr
}

func (p *GRegisterReactor) OnMessage(c *antgo.Conn, pt antgo.Packet) bool {
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

type GEndReactor struct {
	reactor.TCPReactor
	gateway *Gateway
}

func (p *GEndReactor) OnConnect(c *antgo.Conn) string {
	addr := c.RemoteAddr()
	fmt.Println("OnConnect:", addr)
	return addr
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

type GWorkerReactor struct {
	reactor.TCPReactor
	gateway *Gateway
	ant *antgo.Ant
}

func (p *GWorkerReactor) OnConnect(c *antgo.Conn) string {
	addr := c.RemoteAddr()
	fmt.Println("OnConnect:", addr)
	p.ant.Send(0, "prompt", "Welcome to Gateway Server, worker.", c, 0)
	return addr
}

func (p *GWorkerReactor) OnMessage(c *antgo.Conn, pt antgo.Packet) bool {
	fmt.Println("OnMessage:", c.RemoteAddr())
	code := pt.Code()
	event := pt.Event()
	msg := pt.Msg()
	fmt.Println(code)
	switch event {
	case "prompt":
		fmt.Println(msg)
	case "worker_connect":
		fmt.Println("=====worker")
		fmt.Println(c.RemoteAddr())
	case "ping":
		fmt.Println("ping")
	default:
		p.ant.Send(0, "prompt", "unknow msg", c, 0)
		// p.OnClose(c)
		// c.Close()
	}
	return true
}

type Gateway struct {
	EndAnt      *antgo.Ant
	WorkerAnt   []*antgo.Ant
	RegisterAnt *antgo.Ant

	Addresses []string

	// EndConns     map[string]*antgo.Conn
	// WorkerConns  map[string]*antgo.Conn
	// RegisterConn *antgo.Conn
}

func NewGateway(endTransport string, endIP string, endPort int, endType string,
	workerTransport string, workerIP string, workerPort []int, workerType string,
	registerTransport string, registerIP string, registerPort int, registerType string) *Gateway {
	gateway := &Gateway{
		EndAnt:      nil,
		RegisterAnt: nil,
		WorkerAnt:   make([]*antgo.Ant, 0, 12),
		Addresses:   make([]string, 0, 12)}

	endProtocol := NewProtocol(endType, NewListenDialer(endType, endTransport, endIP, endPort))
	endReactor := &GEndReactor{gateway: gateway}
	gateway.EndAnt = antgo.NewAnt(endTransport, endIP, endPort, antgo.DefaultConfig, endProtocol, endReactor)

	registerProtocol := NewProtocol(registerType, NewListenDialer(registerType, registerTransport, registerIP, registerPort))
	registerReactor := &GRegisterReactor{gateway: gateway}
	gateway.RegisterAnt = antgo.NewAnt(registerTransport, registerIP, registerPort, antgo.DefaultConfig, registerProtocol, registerReactor)

	for _, port := range workerPort {
		workerProtocol := NewProtocol(workerType, NewListenDialer(workerType, workerTransport, workerIP, port))
		workerReactor := &GWorkerReactor{gateway: gateway}
		workerAnt := antgo.NewAnt(workerTransport, workerIP, port, antgo.DefaultConfig, workerProtocol, workerReactor)
		workerReactor.ant = workerAnt
		gateway.WorkerAnt = append(gateway.WorkerAnt, workerAnt)
		gateway.Addresses = append(gateway.Addresses, workerAnt.AddressJoin())
	}
	return gateway
}

func (p *Gateway) connectRegister() {
	go p.RegisterAnt.Dial(Timeout)
}

func (p *Gateway) pingEnd() {
	p.EndAnt.Send(0, "ping", "", nil, 0)
}

func (p *Gateway) pingRegister() {
	p.RegisterAnt.Send(0, "ping", "", nil, 0)
}

func (p *Gateway) pingWorker() {
	for _, WorkerAnt := range p.WorkerAnt {
		WorkerAnt.Send(0, "ping", "", nil, 0)
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
	antgo.Stop()
}
