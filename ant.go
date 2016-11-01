package antgo

import (
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
	return &Ant{
		Transport: transport,
		IP:        ip,
		Port:      port,

		config:   config,
		Reactor:  reactor,
		protocol: protocol,

		exitChan:  make(chan struct{}),
		waitGroup: &sync.WaitGroup{},
	}
}

func (ant *Ant) Listen(acceptTimeout time.Duration) {
	listenspeaker := ant.protocol.ListenSpeaker()
	ant.waitGroup.Add(1)

	defer func() {
		listenspeaker.Close()
		ant.waitGroup.Done()
	}()

	for {
		select {
		case <-ant.exitChan:
			return

		default:
		}

		listenspeaker.SetDeadline(time.Now().Add(acceptTimeout))
		netConn, err := listenspeaker.Accept()
		if err != nil {
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

func (ant *Ant) Speak(acceptTimeout time.Duration) *Conn {
	listenspeaker := ant.protocol.ListenSpeaker()
	listenspeaker.SetDeadline(time.Now().Add(acceptTimeout))
	netConn, err := listenspeaker.Dial()
	if err != nil {
		return nil
	}
	conn := newConn(netConn, ant)
	conn.Do(ant)
	return conn
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
