package forwarder

import (
	"context"
	"errors"
	"github.com/itsabgr/forwarder-go/internal/allocator"
	"github.com/itsabgr/forwarder-go/pkg/addr"
	"net"
	"net/netip"
	"runtime"
	"testing"
	"time"
)

func TestEcho(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	deadline, ok := ctx.Deadline()
	if !ok {
		panic(errors.New("not ok"))
	}
	rela := must(net.ListenUDP("udp", must(net.ResolveUDPAddr("udp", "127.0.0.1:0"))))
	defer rela.Close()
	allo := allocator.New(rela.LocalAddr().(*net.UDPAddr).AddrPort())
	ech := must(Listen("udp", must(net.ResolveUDPAddr("udp", "127.0.0.1:0")).AddrPort(), allo, t))
	defer ech.Close()
	cli := must(Listen("udp", must(net.ResolveUDPAddr("udp", "127.0.0.1:0")).AddrPort(), allo, t))
	defer cli.Close()
	go Relay(rela, ech.conn.LocalAddr(), t)
	go allocator.Serve(ech, t)
	throw(cli.SetDeadline(deadline))
	throw(allo.Ping())
	go cli.ReadFrom(nil)
	for {
		runtime.Gosched()
		if err := ctx.Err(); err != nil {
			t.Fatal(err)
		}
		if cli.LocalAddr() != nil {
			_ = cli.LocalAddr().String()
			break
		}
	}
	if netip.MustParseAddrPort(string(cli.LocalAddr().(*addr.Addr).ID())).String() != cli.conn.LocalAddr().String() {
		t.FailNow()
	}
}
