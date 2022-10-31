package forwarder

import (
	"errors"
	"github.com/itsabgr/forwarder-go/pkg/addr"
)

func decodePacket(b []byte) (addr.ID, []byte, bool) {
	if len(b) <= 1 {
		return nil, nil, false
	}
	l := b[0]
	if l == 0 {
		return nil, b[1:], true
	}
	if len(b)-1 < int(l) {
		return nil, nil, false
	}
	return b[1 : 1+l], b[1+l:], true
}
func encodePacket(id addr.ID, data []byte) ([]byte, bool) {
	if len(id) == 0 {
		return append([]byte{0}, data...), true
	}
	if !id.Valid() {
		panic(errors.New("invalid id"))
	}
	return append(id.Encode(), data...), true
}
