package database

import (
	"encoding/binary"
)

// itob encodes unsigned 64-bit integer to byte slices to improve performance when used as BoltDB keys
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
