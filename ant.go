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

func (ant *Ant) Listen(acceptTimeout time.Duration) {
	listendialer := ant.protocol.ListenDialer()
	listendialer.Listen()
	fmt.Println("Listen ", ant.Address, "...")
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
			conn := newConn(netConn, ant)
			conn.Do(ant)
			ant.Conns = append(ant.Conns, conn)
			WaitGroup.Done()
		}()
	}
}

func (ant *Ant) Dial(acceptTimeout time.Duration) {
	listendialer := ant.protocol.ListenDialer()
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
	fmt.Println("Dial ", ant.Address, "...")

	WaitGroup.Add(1)
	go func() {
		conn := newConn(netConn, ant)
		conn.Do(ant)
		ant.Conns = append(ant.Conns, conn)
		WaitGroup.Done()
	}()
}

func (ant *Ant) Send(code int, event string, msg interface{}, conn *Conn, timeout time.Duration) {
	if conn == nil {
		for _, conn := range ant.Conns {
			conn.AsyncWritePacket(ant.protocol.Deserialize(code, event, msg), timeout)
		}
	} else {
		conn.AsyncWritePacket(ant.protocol.Deserialize(code, event, msg), timeout)
	}
}

func (ant *Ant) AddressJoin() string{
    return Fastjoin(":", ant.Transport, ant.IP, strconv.Itoa(ant.Port), ant.protocol.Type())
}

func (ant *Ant) Run() {

}

func (ant *Ant) Reload() {

}

func (ant *Ant) BeforeStart() {

}

func (ant *Ant) AfterStart() {

}

func (ant *Ant) BeforeReload() {

}

func (ant *Ant) AfterReload() {

}

func (ant *Ant) BeforeStop() {

}

func (ant *Ant) AfterStop() {

}

func Stop() {
	close(ExitChan)
	WaitGroup.Wait()
}
