package forwarder

import (
	"errors"
	"github.com/go-log/log"
	"github.com/itsabgr/forwarder-go/pkg/addr"
	"net"
	"net/netip"
	"time"
)

type Client struct {
	conn       *net.UDPConn
	allocator  Allocator
	logger     log.Logger
	bufferSize int
}

func (c *Client) ReadFrom(p []byte) (int, net.Addr, error) {
	b, from, err := c.read(c.bufferSize)
	if err != nil {
		return 0, nil, err
	}
	if from.ID() == nil {
		c.logger.Log("forwarder: client: pong")
		if err := c.allocator.Pong(from.Relay(), b); err != nil {
			return 0, nil, err
		}
		return 0, nil, err
	}
	c.logger.Log("forwarder: client: msg")
	return copy(p, b), from, nil
}

func (c *Client) read(n int) ([]byte, *addr.Addr, error) {
again:
	b := make([]byte, n)
	n, from, err := c.conn.ReadFromUDPAddrPort(b)
	if err != nil {
		return nil, nil, err
	}
	c.logger.Log("read")
	id, b, ok := decodePacket(b[:n])
	if !ok {
		goto again
	}
	return b, addr.New(from, id), nil
}

func (c *Client) WriteTo(p []byte, netAddr net.Addr) (n int, err error) {
	a := netAddr.(*addr.Addr)
	if a.ID() == nil {
		panic("no id")
	}
	pack, ok := encodePacket(a.ID(), p)
	if !ok {
		panic(errors.New("not ok"))
	}
	_, err = c.conn.WriteToUDPAddrPort(pack, a.Relay())
	return len(p), err
}

func (c *Client) Close() error {
	err := c.conn.Close()
	if err == nil {
		c.logger.Log("forwarder: client: closed")
		if c.allocator != nil {
			c.allocator.Release()
		}
	}
	return err
}

func (c *Client) LocalAddr() net.Addr {
	if addr := c.allocator.Addr(); addr != nil {
		return addr
	}
	return nil
}

func (c *Client) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *Client) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)

}

func (c *Client) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
func Wrap(udp *net.UDPConn, allocator Allocator, logger log.Logger) (_ *Client, err error) {
	c := &Client{conn: udp, allocator: allocator, logger: logger, bufferSize: 1024}
	ping := func(relay netip.AddrPort, data []byte) error {
		c.logger.Log("forwarder: chat: ping")
		_, err := c.conn.WriteToUDPAddrPort(append([]byte{0}, data...), relay)
		return err
	}
	if allocator != nil {
		if err = allocator.Init(ping); err != nil {
			udp.Close()
			return nil, err
		}
		c.logger.Log("allocator: inited")
	}
	return c, err
}
func Listen(network string, local netip.AddrPort, allocator Allocator, logger log.Logger) (*Client, error) {
	udp, err := net.ListenUDP(network, net.UDPAddrFromAddrPort(local))
	if err != nil {
		return nil, err
	}
	return Wrap(udp, allocator, logger)
}
