package antgo

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrConnClosing   = errors.New("use of closed network connection")
	ErrWriteBlocking = errors.New("write packet was blocking")
	ErrReadBlocking  = errors.New("read packet was blocking")
)

type Conn struct {
	Id          string
	netConn     net.Conn      // the raw connection
	extraData   interface{}   // to save extra data
	closeOnce   sync.Once     // close the netConn, once, per instance
	closeFlag   int32         // close flag
	closeChan   chan struct{} // close chanel
	sendChan    chan Packet   // packet send chanel
	receiveChan chan Packet   // packeet receive chanel
}

func newConn(netConn net.Conn, ant *Ant) *Conn {
	return &Conn{
		netConn:     netConn,
		closeChan:   make(chan struct{}),
		sendChan:    make(chan Packet, ant.Config.PacketSendChanLimit),
		receiveChan: make(chan Packet, ant.Config.PacketReceiveChanLimit),
	}
}

func (c *Conn) GetExtraData() interface{} {
	return c.extraData
}

func (c *Conn) PutExtraData(data interface{}) {
	c.extraData = data
}

func (c *Conn) GetRawConn() net.Conn {
	return c.netConn
}

func (c *Conn) Close() {
	c.closeOnce.Do(func() {
		atomic.StoreInt32(&c.closeFlag, 1)
		close(c.closeChan)
		close(c.sendChan)
		close(c.receiveChan)
		c.netConn.Close()
	})
}

func (c *Conn) IsClosed() bool {
	return atomic.LoadInt32(&c.closeFlag) == 1
}

func (c *Conn) Do(ant *Ant) {
	if ant.OnConnect(c) == "" {
		return
	}

	asyncDo(c.handleLoop, ant)
	asyncDo(c.readLoop, ant)
	asyncDo(c.writeLoop, ant)
}

func asyncDo(fn func(ant *Ant), ant *Ant) {
	WaitGroup.Add(1)
	go func() {
		fn(ant)
		WaitGroup.Done()
	}()
}

func (c *Conn) readLoop(ant *Ant) {
	defer func() {
		recover()
		ant.OnClose(c)
		c.Close()
	}()
	for {
		select {
		case <-ExitChan:
			return

		case <-c.closeChan:
			return

		default:
		}
		p := ant.protocol.ReadPacket(c.netConn)
		if p == nil {
			return
		} else {
			c.receiveChan <- p
		}
	}
}

func (c *Conn) handleLoop(ant *Ant) {
	defer func() {
		recover()
		ant.OnClose(c)
		c.Close()
	}()

	for {
		select {
		case <-ExitChan:
			return

		case <-c.closeChan:
			return

		case p := <-c.receiveChan:
			if c.IsClosed() {
				return
			} else {
				ant.OnMessage(c, p)
			}
		}
	}
}

func (c *Conn) AsyncWritePacket(p Packet, timeout time.Duration) (err error) {
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

func (c *Conn) writeLoop(ant *Ant) {
	defer func() {
		recover()
		ant.OnClose(c)
		c.Close()
	}()

	for {
		select {
		case <-ExitChan:
			return

		case <-c.closeChan:
			return

		case p := <-c.sendChan:
			if c.IsClosed() {
				return
			}
			if _, err := c.netConn.Write(p.Serialize()); err != nil {
				return
			}
		}
	}
}

func (c *Conn) LocalAddr() string {
	return c.netConn.LocalAddr().String()
}

func (c *Conn) RemoteAddr() string {
	return c.netConn.RemoteAddr().String()
}
