package antgo

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"
	// "fmt"
)

// Error type
var (
	ErrConnClosing   = errors.New("use of closed network connection")
	ErrWriteBlocking = errors.New("write packet was blocking")
	ErrReadBlocking  = errors.New("read packet was blocking")
)

// Connection exposes a set of callbacks for the various events that occur on a connection
type Connection struct {
	srv         *Server
	conn        net.Conn      // the raw connection
	extraData   interface{}   // to save extra data
	closeOnce   sync.Once     // close the conn, once, per instance
	closeFlag   int32         // close flag
	closeChan   chan struct{} // close chanel
	sendChan    chan Packet   // packet send chanel
	receiveChan chan Packet   // packeet receive chanel
}

// newConn returns a wrapper of raw conn
func newConnection(conn net.Conn, srv *Server) *Connection {
	return &Connection{
		srv:         srv,
		conn:        conn,
		closeChan:   make(chan struct{}),
		sendChan:    make(chan Packet, srv.config.PacketSendChanLimit),
		receiveChan: make(chan Packet, srv.config.PacketReceiveChanLimit),
	}
}

// GetExtraData gets the extra data from the Connection
func (c *Connection) GetExtraData() interface{} {
	return c.extraData
}

// PutExtraData puts the extra data with the Connection
func (c *Connection) PutExtraData(data interface{}) {
	c.extraData = data
}

// GetRawConn returns the raw net.Conn from the Connection
func (c *Connection) GetRawConn() net.Conn {
	return c.conn
}

// Close closes the connection
func (c *Connection) Close() {
	// fmt.Println(&c.conn)
	c.closeOnce.Do(func() {
		atomic.StoreInt32(&c.closeFlag, 1)
		close(c.closeChan)
		close(c.sendChan)
		close(c.receiveChan)
		c.conn.Close()
		c.srv.reactor.OnClose(c)
	})
}

// IsClosed indicates whether or not the connection is closed
func (c *Connection) IsClosed() bool {
	return atomic.LoadInt32(&c.closeFlag) == 1
}

// AsyncWritePacket async writes a packet, this method will never block
func (c *Connection) AsyncWritePacket(p Packet, timeout time.Duration) (err error) {
	if c.IsClosed() {
		return ErrConnClosing
	}

	defer func() {
		if e := recover(); e != nil {
			err = ErrConnClosing
		}
	}()

	if timeout == 0 {
		select {
		case c.sendChan <- p:
			return nil

		default:
			return ErrWriteBlocking
		}

	} else {
		select {
		case c.sendChan <- p:
			return nil

		case <-c.closeChan:
			return ErrConnClosing

		case <-time.After(timeout):
			return ErrWriteBlocking
		}
	}
}

// Do it
func (c *Connection) Do() {
	if !c.srv.reactor.OnConnect(c) {
		return
	}

	asyncDo(c.handleLoop, c.srv.waitGroup)
	asyncDo(c.readLoop, c.srv.waitGroup)
	asyncDo(c.writeLoop, c.srv.waitGroup)
}

func (c *Connection) readLoop() {
	defer func() {
		recover()
		c.Close()
	}()

	for {
		select {
		case <-c.srv.exitChan:
			return

		case <-c.closeChan:
			return

		default:
		}

		p := c.srv.protocol.ReadPacket(c.conn)

		for q := range p {
			c.receiveChan <- q
		}
	}
}

func (c *Connection) writeLoop() {
	defer func() {
		recover()
		c.Close()
	}()

	for {
		select {
		case <-c.srv.exitChan:
			return

		case <-c.closeChan:
			return

		case p := <-c.sendChan:
			if c.IsClosed() {
				return
			}
			if _, err := c.conn.Write(p.Serialize()); err != nil {
				return
			}
		}
	}
}

func (c *Connection) handleLoop() {
	defer func() {
		recover()
		c.Close()
	}()

	for {
		select {
		case <-c.srv.exitChan:
			return

		case <-c.closeChan:
			return

		case p := <-c.receiveChan:
			if c.IsClosed() {
				return
			}
			if !c.srv.reactor.OnMessage(c, p) {
				return
			}
		}
	}
}

func asyncDo(fn func(), wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		fn()
		wg.Done()
	}()
}
