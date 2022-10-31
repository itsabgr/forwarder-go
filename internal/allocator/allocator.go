package allocator

import (
	"errors"
	"fmt"
	"github.com/itsabgr/forwarder-go/pkg/addr"
	"net/netip"
	"sync/atomic"
)

type Allocator struct {
	relay netip.AddrPort
	addr  atomic.Pointer[addr.Addr]
	ping  func(relay netip.AddrPort, data []byte) error
}

func New(relay netip.AddrPort) *Allocator {
	return &Allocator{
		relay: relay,
	}
}
func (s *Allocator) Ping() error {
	return s.ping(s.relay, []byte("raw"))
}
func (s *Allocator) Init(ping func(relay netip.AddrPort, data []byte) error) error {
	s.ping = ping
	return nil
}

func (s *Allocator) Release() {
	s.ping = nil
}
func eqNetIPAddrPort(a, b netip.AddrPort) bool {
	if a.Port() != b.Port() {
		return false
	}
	if a.Addr().Is4In6() {
		a = netip.AddrPortFrom(netip.AddrFrom4(a.Addr().As4()), a.Port())
	}
	if b.Addr().Is4In6() {
		b = netip.AddrPortFrom(netip.AddrFrom4(b.Addr().As4()), b.Port())
	}
	return a.Addr() == b.Addr()

}
func (s *Allocator) Pong(relay netip.AddrPort, data []byte) error {
	if !eqNetIPAddrPort(s.relay, relay) {
		fmt.Println(s.relay, relay)
		return nil
	}
	if !addr.ID.Valid(data) {
		return errors.New("got invalid id from relay")
	}
	s.addr.Store(addr.New(relay, data))
	return nil
}

func (s *Allocator) Addr() *addr.Addr {
	addr := s.addr.Load()
	if addr == nil {
		return nil
	}
	if addr.ID() == nil {
		return nil
	}
	return addr
}
