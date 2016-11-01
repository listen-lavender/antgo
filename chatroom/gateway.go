package chatroom

import (
    "fmt"
	"../../antgo"
	"../../antgo/reactor"
)

type GRegisterReactor struct {
	reactor.TCPReactor
}

func (p *GRegisterReactor) OnConnect(c *antgo.Conn) {
}

func (p *GRegisterReactor) OnMessage(c *antgo.Conn, pt antgo.Packet) {
}

func (p *GRegisterReactor) OnClose() {
}

type EndReactor struct {
	reactor.TCPReactor
}

func (p *EndReactor) OnConnect(c *antgo.Conn) {
}

func (p *EndReactor) OnMessage(c *antgo.Conn, pt antgo.Packet) {
	event := pt.Event()
	msg := antgo.JsonDecode(pt.Msg())
	secret := msg["secret"]
	fmt.Println(secret)
	fmt.Println(event)
}

func (p *EndReactor) OnClose(c *antgo.Conn) {
}

type GWorkerReactor struct {
	reactor.TCPReactor
	allGatewayAddr  map[string]string
	busyGatewayAddr map[string]string
	idleGatewayAddr map[string]string
}

func (p *GWorkerReactor) OnConnect(c *antgo.Conn) {
}

func (p *GWorkerReactor) OnMessage(c *antgo.Conn, pt antgo.Packet) {
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
		fmt.Println("Receive bad event:$event from Worker.\n")
	}
}

func (p *GWorkerReactor) OnClose(c *antgo.Conn) {
}

type Gateway struct {
	reactor.TCPReactor
	OuterAnt     *antgo.Ant
	InnerAnt     *antgo.Ant
	WorkerConns  []*antgo.Conn
	RegisterConn *antgo.Conn
	EndConns     []*antgo.Conn
}

func (*Gateway) pingRegister(){

}

func (*Gateway) pingWorker(){

}

func (*Gateway) run(){

}

