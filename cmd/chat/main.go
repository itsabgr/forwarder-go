package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/go-log/log/log"
	"github.com/itsabgr/forwarder-go"
	"github.com/itsabgr/forwarder-go/internal/allocator"
	"github.com/itsabgr/forwarder-go/pkg/addr"
	"net/netip"
	"os"
	"runtime"
	"strings"
	"time"
)

var flagRelay = flag.String("relay", "", "listening addr")
var flagNetwork = flag.String("net", "udp", "network")

func init() {
	flag.Parse()
}

// TODO fix bugs
func main() {
	logger := log.New()
	alloc := allocator.New(netip.MustParseAddrPort(*flagRelay))
	cli, err := forwarder.Listen(*flagNetwork, netip.AddrPort{}, alloc, logger)
	if err != nil {
		logger.Log(err)
		os.Exit(1)
	}
	defer cli.Close()
	logger.Log("allocating...")
	for {
		if err := cli.SetDeadline(time.Now().Add(time.Millisecond * 100)); err != nil {
			logger.Log(err)
			os.Exit(1)
		}
		if err := alloc.Ping(); err != nil {
			logger.Log(err)
			os.Exit(1)
		}
		if _, _, err := cli.ReadFrom(nil); err != nil && !os.IsTimeout(err) {
			logger.Log(err)
			os.Exit(1)
		}
		if addr := cli.LocalAddr(); addr != nil {
			break
		}
	}
	if err := cli.SetDeadline(time.Time{}); err != nil {
		logger.Log(err)
		os.Exit(1)
	}
	fmt.Println(cli.LocalAddr())
	scanner := bufio.NewScanner(os.Stdin)
	for {
		runtime.Gosched()
		fmt.Print("peer >")
		if !scanner.Scan() {
			break
		}
		peer, err := addr.Resolve(strings.TrimSpace(scanner.Text()))
		if err != nil {
			fmt.Println("error: addr:", err)
			continue
		}
		fmt.Print("msg  >")
		if !scanner.Scan() {
			break
		}
		msg := scanner.Bytes()
		if len(msg) > 200 {
			fmt.Println("error: msg:", "too large msg")
			continue
		}
		_, err = cli.WriteTo(msg, peer)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
	}
}
