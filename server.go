package antgo

import (
	"sync"
	"time"
	"fmt"
)

type Config struct {
	PacketSendChanLimit    uint32 // the limit of packet send channel
	PacketReceiveChanLimit uint32 // the limit of packet receive channel
}

type Server struct {
	config    *Config            // server configuration
	reactor   Reactor    // message callbacks in connection
	protocol  Protocol  // customize packet protocol
	exitChan  chan struct{}      // notify all goroutines to shutdown
	waitGroup *sync.WaitGroup    // wait for all goroutines
}

// NewServer creates a server
func NewServer(config *Config, reactor Reactor, protocol Protocol) *Server {
	return &Server{
		config:    config,
		reactor:   reactor,
		protocol:  protocol,
		exitChan:  make(chan struct{}),
		waitGroup: &sync.WaitGroup{},
	}
}

// Start starts service
func (s *Server) Start(acceptTimeout time.Duration) {
	listener := s.protocol.GetListener()
	s.waitGroup.Add(1)

	defer func() {
		fmt.Println("abc")
		listener.Close()
		s.waitGroup.Done()
	}()

	for {
		select {
		case <-s.exitChan:
			return

		default:
		}

		listener.SetDeadline(time.Now().Add(acceptTimeout))
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		s.waitGroup.Add(1)
		go func() {
			newConnection(conn, s).Do()
			s.waitGroup.Done()
		}()
	}
}

// Stop stops service
func (s *Server) Stop() {
	close(s.exitChan)
	s.waitGroup.Wait()
}
