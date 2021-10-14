package localdb

import (
	"errors"

	"go.etcd.io/bbolt"
)

func rollback(tx *bbolt.Tx, err *error) {
	if cerr := tx.Rollback(); errors.Is(cerr, bbolt.ErrTxClosed) && *err == nil {
		*err = cerr
	}
}
