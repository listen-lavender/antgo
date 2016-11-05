package antgo

import (
	"fmt"
	"strconv"
	"sync"
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
	Ptype     string // protocol type, name of New method
	Rtype     string // reactor type, name of struct
	Conns     []*Conn

	config   *Config  // ant configuration
	protocol Protocol // ant protocol
	Reactor           // ant reactor

	exitChan  chan struct{}   // notify all goroutines to shutdown
	waitGroup *sync.WaitGroup // wait for all goroutines
}

// NewAnt creates a ant
func NewAnt(transport string, ip string, port int, config *Config, protocol Protocol, reactor Reactor) *Ant {
	address := Fastjoin(transport, "://", ip, ":", strconv.Itoa(port))
	return &Ant{
		Transport: transport,
		IP:        ip,
		Port:      port,
		Address:   address,
		Conns:     make([]*Conn, 0, 1000),

		config:   config,
		Reactor:  reactor,
		protocol: protocol,

		exitChan:  make(chan struct{}),
		waitGroup: &sync.WaitGroup{},
	}
}

func (ant *Ant) Listen(acceptTimeout time.Duration) {
	listendialer := ant.protocol.ListenDialer()
	listendialer.Listen()
	fmt.Println("Listen ", ant.Address, "...")
	ant.waitGroup.Add(1)

	defer func() {
		listendialer.Close()
		ant.waitGroup.Done()
	}()

	for {
		select {
		case <-ant.exitChan:
			return

		default:
		}

		listendialer.SetDeadline(time.Now().Add(acceptTimeout))
		netConn, err := listendialer.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		ant.waitGroup.Add(1)
		go func() {
			conn := newConn(netConn, ant)
			conn.Do(ant)
			ant.Conns = append(ant.Conns, conn)
			ant.waitGroup.Done()
		}()
	}
}

func (ant *Ant) Dial(acceptTimeout time.Duration) {
	listendialer := ant.protocol.ListenDialer()
	ant.waitGroup.Add(1)
	defer func() {
		listendialer.Close()
		ant.waitGroup.Done()
	}()
	listendialer.SetDeadline(time.Now().Add(acceptTimeout))
	netConn, err := listendialer.Dial()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Dial ", ant.Address, "...")

	ant.waitGroup.Add(1)
	go func() {
		conn := newConn(netConn, ant)
		conn.Do(ant)
		ant.Conns = append(ant.Conns, conn)
		ant.waitGroup.Done()
	}()

	for {
		select {
		case <-ant.exitChan:
			return

		default:
		}
	}
}

func (ant *Ant) Send(event string, msg []byte, timeout time.Duration) {
	for _, conn := range ant.Conns {
		conn.AsyncWritePacket(ant.protocol.Deserialize(event, msg), timeout)
	}
}

func (ant *Ant) Run() {

}

func (ant *Ant) Reload() {

}

func (ant *Ant) Stop() {
	close(ant.exitChan)
	ant.waitGroup.Wait()
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
