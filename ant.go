package antgo

import (
	"fmt"
	"strconv"
	"time"
)

type Config struct {
	PacketSendChanLimit    uint32 // the limit of packet send channel
	PacketReceiveChanLimit uint32 // the limit of packet receive channel
}

type Ant struct {
	Transport string // ant transport tcp or udp
	IP        string // ip
	Port      int    // port
	Address   string // listen address
	Conns     []*Conn

	Config   *Config  // ant configuration
	protocol Protocol // ant protocol
	Reactor           // ant reactor
}

// NewAnt creates a ant
func NewAnt(transport string, ip string, port int, config *Config, protocol Protocol, reactor Reactor) *Ant {
	address := Fastjoin("", transport, "://", ip, ":", strconv.Itoa(port))
	return &Ant{
		Transport: transport,
		IP:        ip,
		Port:      port,
		Address:   address,
		Conns:     make([]*Conn, 0, 1000),

		Config:   config,
		Reactor:  reactor,
		protocol: protocol,
	}
}

func (p *Ant) Listen(acceptTimeout time.Duration) {
	listendialer := p.protocol.ListenDialer()
	if err := listendialer.Listen(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Listen ", p.Address, "...")
	WaitGroup.Add(1)

	defer func() {
		listendialer.Close()
		WaitGroup.Done()
	}()

	for {
		select {
		case <-ExitChan:
			return

		default:
		}

		listendialer.SetDeadline(time.Now().Add(acceptTimeout))
		netConn, err := listendialer.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		WaitGroup.Add(1)
		go func() {
			conn := newConn(netConn, p)
			conn.Do(p)
			p.Conns = append(p.Conns, conn)
			WaitGroup.Done()
		}()
	}
}

func (p *Ant) Dial(acceptTimeout time.Duration) {
	listendialer := p.protocol.ListenDialer()
	WaitGroup.Add(1)

	defer func() {
		listendialer.Close()
		WaitGroup.Done()
	}()

	listendialer.SetDeadline(time.Now().Add(acceptTimeout))
	netConn, err := listendialer.Dial()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Dial ", p.Address, "...")

	WaitGroup.Add(1)
	go func() {
		conn := newConn(netConn, p)
		conn.Do(p)
		p.Conns = append(p.Conns, conn)
		WaitGroup.Done()
	}()
}

func (p *Ant) Send(code int, event string, msg interface{}, conn *Conn, timeout time.Duration) {
	if conn == nil {
		for _, conn := range p.Conns {
			conn.AsyncWritePacket(p.protocol.Deserialize(code, event, msg), timeout)
		}
	} else {
		conn.AsyncWritePacket(p.protocol.Deserialize(code, event, msg), timeout)
	}
}

func (p *Ant) AddressJoin() string {
	return Fastjoin(":", p.Transport, p.IP, strconv.Itoa(p.Port), p.protocol.Type())
}

func (p *Ant) Run() {

}

func (p *Ant) Reload() {

}

func (p *Ant) BeforeStart() {

}

func (p *Ant) AfterStart() {

}

func (p *Ant) BeforeReload() {

}

func (p *Ant) AfterReload() {

}

func (p *Ant) BeforeStop() {

}

func (p *Ant) AfterStop() {

}

func Stop() {
	close(ExitChan)
	WaitGroup.Wait()
}
