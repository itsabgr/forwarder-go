package main

import (
	"flag"
	"github.com/go-log/log/log"
	"github.com/itsabgr/forwarder-go"
	"github.com/itsabgr/forwarder-go/internal/allocator"
	"net"
	"os"
)

var flagAddr = flag.String("addr", "0.0.0.0:0", "listening addr")
var flagNetwork = flag.String("net", "udp", "network")

func init() {
	flag.Parse()
}
func main() {
	logger := log.New()
	addr, err := net.ResolveUDPAddr(*flagNetwork, *flagAddr)
	if err != nil {
		logger.Log(err)
		os.Exit(1)
	}
	conn, err := net.ListenUDP(*flagNetwork, addr)
	if err != nil {
		logger.Log(err)
		os.Exit(1)
	}
	defer conn.Close()
	logger.Logf("addr: %s", conn.LocalAddr())
	cli, err := forwarder.Wrap(conn, nil, logger)
	if err != nil {
		logger.Log(err)
		os.Exit(1)
	}
	defer cli.Close()
	err = allocator.Serve(cli, logger)
	if err != nil {
		logger.Log(err)
		os.Exit(1)
	}
}
