package addr

import (
	"encoding/hex"
)

const MaxIDLength = 250

func (id ID) String() string {
	return hex.EncodeToString(id)
}
func (id ID) Valid() bool {
	return len(id) < MaxIDLength
}
func ParseID(id string) (ID, error) {
	return hex.DecodeString(id)
}
