package addr

import (
	"errors"
	"net/netip"
	"strings"
)

type ID []byte

type Addr struct {
	relay netip.AddrPort
	id    ID
}

func New(relay netip.AddrPort, id ID) *Addr {
	return &Addr{id: id, relay: relay}
}
func (id ID) Encode() []byte {
	return append([]byte{byte(len(id))}, id...)
}

var ErrInvalidAddr = errors.New("forwarder: invalid addr")

func Resolve(addr string) (*Addr, error) {
	if len(addr) < 2 {
		return nil, ErrInvalidAddr
	}
	last := strings.LastIndex(addr, "-")
	if last < 0 || last == len(addr)-1 {
		return nil, ErrInvalidAddr
	}
	addrPort, err := netip.ParseAddrPort(addr[:last])
	if err != nil {
		return nil, ErrInvalidAddr
	}
	id, err := ParseID(addr[last+1:])
	if err != nil {
		return nil, ErrInvalidAddr
	}
	if !id.Valid() {
		return nil, ErrInvalidAddr
	}
	return &Addr{addrPort, id}, nil
}

func (a *Addr) ID() ID {
	return a.id
}
func (a *Addr) Relay() netip.AddrPort {
	return a.relay
}

func (a *Addr) String() string {
	return a.relay.String() + "-" + a.ID().String()
}
func (a *Addr) Network() string {
	if a.relay.Addr().Is4() {
		return "forwarder+udp4"
	}
	if a.relay.Addr().Is4In6() {
		return "forwarder+udp"
	}
	return "forwarder+udp6"
}
