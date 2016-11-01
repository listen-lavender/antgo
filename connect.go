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
		sendChan:    make(chan Packet, ant.config.PacketSendChanLimit),
		receiveChan: make(chan Packet, ant.config.PacketReceiveChanLimit),
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
	if ant.OnConnect(c) == nil {
		return
	}

	asyncDo(c.handleLoop, ant)
	asyncDo(c.readLoop, ant)
	asyncDo(c.writeLoop, ant)
}

func asyncDo(fn func(ant *Ant), ant *Ant) {
	wg := ant.waitGroup
	wg.Add(1)
	go func() {
		fn(ant)
		wg.Done()
	}()
}

func (c *Conn) readLoop(ant *Ant) {
	defer func() {
		recover()
		c.Close()
		ant.OnClose(c)
	}()

	for {
		select {
		case <-ant.exitChan:
			return

		case <-c.closeChan:
			return

		default:
		}

		p := ant.protocol.ReadPacket(c.netConn)

		for q := range p {
			c.receiveChan <- q
		}
	}
}

func (c *Conn) handleLoop(ant *Ant) {
	defer func() {
		recover()
		c.Close()
		ant.OnClose(c)
	}()

	for {
		select {
		case <-ant.exitChan:
			return

		case <-c.closeChan:
			return

		case p := <-c.receiveChan:
			if c.IsClosed() {
				return
			}
			if !ant.OnMessage(c, p) {
				return
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
		c.Close()
		ant.OnClose(c)
	}()

	for {
		select {
		case <-ant.exitChan:
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

func (c *Conn) LocalAddr() net.Addr {
	return c.netConn.LocalAddr()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.netConn.RemoteAddr()
}
