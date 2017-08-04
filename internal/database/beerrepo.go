package database

import (
	"encoding/json"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

// BeerRepo represents a boltdb repository for managing beers
type BeerRepo struct {
	db *bolt.DB
}

// NewBeerRepo returns a new Beer repository using the given bolt database. It also creates the Beers
// bucket if it is not yet created on disk.
func NewBeerRepo(db *bolt.DB) (*BeerRepo, error) {

	tx, err := db.Begin(true)
	if err != nil {
		return nil, errors.Wrap(err, "Could not begin transaction")
	}

	defer tx.Rollback()

	if _, err := tx.CreateBucketIfNotExists([]byte("Beers")); err != nil {
		return nil, errors.Wrap(err, "Could not create Beer bucket")
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "Could not commit transaction")
	}

	return &BeerRepo{
		db: db,
	}, nil
}

// Get returns a Beer by its ID
func (r *BeerRepo) Get(id uint64) (*internal.Beer, error) {
	var b *internal.Beer

	err := r.db.View(func(tx *bolt.Tx) error {
		if v := tx.Bucket([]byte("Beers")).Get(itob(id)); v != nil {
			if err := json.Unmarshal(v, &b); err != nil {
				return errors.Wrapf(err, "Could not unmarshal Beer %d", id)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return b, nil
}

// Save creates or updates a Beer
func (r *BeerRepo) Save(b *internal.Beer) error {
	err := r.db.Update(func(tx *bolt.Tx) error {

		bu := tx.Bucket([]byte("Beers"))

		if v := bu.Get(itob(b.ID)); v == nil {
			seq, _ := bu.NextSequence()
			b.ID = seq
		}

		b.ModTime = time.Now()

		if v, err := json.Marshal(b); err != nil {
			return errors.Wrapf(err, "Could not marshal Beer %d", b.ID)
		} else if err := bu.Put(itob(b.ID), v); err != nil {
			return errors.Wrapf(err, "Could not put Beer %d", b.ID)
		}

		return nil

	})

	return err
}
