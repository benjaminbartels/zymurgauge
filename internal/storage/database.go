package storage

import (
	"encoding/binary"
	"errors"

	"github.com/boltdb/bolt"
)

// itob encodes unsigned 64-bit integer to byte slices to improve performance when used as BoltDB keys.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)

	return b
}

func rollback(tx *bolt.Tx, err *error) {
	if cerr := tx.Rollback(); errors.Is(cerr, bolt.ErrTxClosed) && *err == nil {
		*err = cerr
	}
}
