package main

import (
	"flag"
	"github.com/go-log/log/log"
	"github.com/itsabgr/forwarder-go"
	"net"
	"os"
)

var flagAddr = flag.String("addr", ":0", "listening addr")
var flagZero = flag.String("zero", "", "zero addr")
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
	zero, err := net.ResolveUDPAddr(*flagNetwork, *flagZero)
	if err != nil {
		logger.Log(err)
		os.Exit(1)
	}
	logger.Logf("zero: %s", zero)
	conn, err := net.ListenUDP(*flagNetwork, addr)
	if err != nil {
		logger.Log(err)
		os.Exit(1)
	}
	defer conn.Close()
	logger.Logf("addr: %s", conn.LocalAddr())
	err = forwarder.Relay(conn, zero, logger)
	if err != nil {
		logger.Log(err)
		os.Exit(1)
	}
}
