package allocator

import (
	"github.com/go-log/log"
	"github.com/itsabgr/forwarder-go/pkg/addr"
	"net"
)

func Serve(client net.PacketConn, logger log.Logger) error {
	b := make([]byte, 1)
	for {
		_, from, err := client.ReadFrom(b)
		if err != nil {
			return err
		}
		logger.Log("echo")
		client.WriteTo(from.(*addr.Addr).ID(), from)
	}
}
