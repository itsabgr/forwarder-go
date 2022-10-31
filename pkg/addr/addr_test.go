package addr

import (
	crand "crypto/rand"
	"encoding/hex"
	"io"
	"testing"
)

func throw(e error) {
	if e != nil {
		panic(e)
	}
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
func TestAddr(t *testing.T) {
	str := "1.1.1.1:8080-" + hex.EncodeToString(rand(128))
	addr := must(Resolve(str))
	if addr.String() != str {
		t.FailNow()
	}
}
