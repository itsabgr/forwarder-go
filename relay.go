package forwarder

import (
	"errors"
	golog "github.com/go-log/log"
	"github.com/itsabgr/forwarder-go/pkg/addr"
	"net"
	"net/netip"
)

func Relay(conn *net.UDPConn, zero net.Addr, log golog.Logger) error {
	for {
		b := make([]byte, 1024)
		n, from, err := conn.ReadFrom(b)
		if err != nil {
			return err
		}
		id, b, ok := decodePacket(b[:n])
		if !ok {
			panic(errors.New("not ok"))
		}
		if id == nil {
			log.Log("relay: ping")
			pack, ok := encodePacket(addr.ID(from.String()), b)
			if !ok {
				panic(errors.New("not ok"))
			}
			must(conn.WriteTo(pack, zero))
		} else {
			if from.String() == zero.String() {
				log.Log("relay: pong")
				pack, ok := encodePacket(nil, b)
				if !ok {
					panic(errors.New("not ok"))
				}
				must(conn.WriteToUDPAddrPort(pack, netip.MustParseAddrPort(string(id))))
			} else {
				log.Log("relay: packet ")
				pack, ok := encodePacket(addr.ID(from.String()), b)
				if !ok {
					panic(errors.New("not ok"))
				}
				must(conn.WriteToUDPAddrPort(pack, netip.MustParseAddrPort(string(id))))
			}
		}
	}
}
