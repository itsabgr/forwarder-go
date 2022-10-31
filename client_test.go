package forwarder

import (
	"context"
	"errors"
	"github.com/itsabgr/forwarder-go/internal/allocator"
	"net"
	"runtime"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	deadline, ok := ctx.Deadline()
	if !ok {
		panic(errors.New("not ok"))
	}
	rela := must(net.ListenUDP("udp", must(net.ResolveUDPAddr("udp", "127.0.0.1:0"))))
	defer rela.Close()
	allo1 := allocator.New(rela.LocalAddr().(*net.UDPAddr).AddrPort())
	allo2 := allocator.New(rela.LocalAddr().(*net.UDPAddr).AddrPort())
	ech := must(Listen("udp", must(net.ResolveUDPAddr("udp", "127.0.0.1:0")).AddrPort(), nil, t))
	defer ech.Close()
	cli1 := must(Listen("udp", must(net.ResolveUDPAddr("udp", "127.0.0.1:0")).AddrPort(), allo1, t))
	defer cli1.Close()
	cli2 := must(Listen("udp", must(net.ResolveUDPAddr("udp", "127.0.0.1:0")).AddrPort(), allo2, t))
	defer cli2.Close()
	go Relay(rela, ech.conn.LocalAddr(), t)
	go allocator.Serve(ech, t)
	throw(cli1.SetDeadline(deadline))
	throw(allo1.Ping())
	go cli1.ReadFrom(nil)
	for {
		runtime.Gosched()
		if err := ctx.Err(); err != nil {
			t.Fatal(err)
		}
		if cli1.LocalAddr() != nil {
			_ = cli1.LocalAddr().String()
			break
		}
	}
	throw(cli2.SetDeadline(deadline))
	throw(allo2.Ping())
	go cli2.ReadFrom(nil)
	for {
		runtime.Gosched()
		if err := ctx.Err(); err != nil {
			t.Fatal(err)
		}
		if cli2.LocalAddr() != nil {
			_ = cli2.LocalAddr().String()
			break
		}
	}
	for range make([]struct{}, 100) {
		data := rand(50)
		must(cli1.WriteTo(data, cli2.LocalAddr()))
		n, from, err := cli2.ReadFrom(data)
		throw(err)
		if from.String() != cli1.LocalAddr().String() {
			t.FailNow()
		}
		if n != len(data) {
			t.FailNow()
		}
		data = rand(len(data) + 10)
		must(cli2.WriteTo(data, cli1.LocalAddr()))
		n, from, err = cli1.ReadFrom(data)
		throw(err)
		if from.String() != cli2.LocalAddr().String() {
			t.FailNow()
		}
		if n != len(data) {
			t.FailNow()
		}
	}
}
