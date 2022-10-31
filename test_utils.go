package forwarder

import (
	"io"
	"net"
)
import crand "crypto/rand"

func throw(e error) {
	if e != nil {
		panic(e)
	}
}

func freePort() (int, error) {
	sock, err := net.ListenUDP("udp", nil)
	if err != nil {
		return -1, err
	}
	defer sock.Close()
	return sock.LocalAddr().(*net.UDPAddr).Port, nil
}
func must[R any](r R, e error) R {
	throw(e)
	return r
}
func rand(n int) []byte {
	b := make([]byte, n)
	must(io.ReadFull(crand.Reader, b))
	return b
}
