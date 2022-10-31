package forwarder

import (
	"github.com/itsabgr/forwarder-go/pkg/addr"
	"net/netip"
)

type Allocator interface {
	Init(ping func(relay netip.AddrPort, data []byte) error) error
	Release()
	Pong(bridge netip.AddrPort, data []byte) error
	Addr() *addr.Addr
}
