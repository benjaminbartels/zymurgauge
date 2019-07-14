package boltdb

import (
	"github.com/boltdb/bolt"
)

func rollback(tx *bolt.Tx, err *error) {
	if cerr := tx.Rollback(); cerr != nil && cerr != bolt.ErrTxClosed && *err == nil {
		*err = cerr
	}
}
