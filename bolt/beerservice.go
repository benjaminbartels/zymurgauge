package bolt

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
	"github.com/benjaminbartels/zymurgauge"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// check implementation at compile time
var _ zymurgauge.BeerService = &BeerService{}

// BeerService represents a service for managing beers
type BeerService struct {
	db     *bolt.DB
	logger *logrus.Logger
}

// Get returns a Beer by its ID
func (s *BeerService) Get(id uint64) (*zymurgauge.Beer, error) {
	var b zymurgauge.Beer

	err := s.db.View(func(tx *bolt.Tx) error {
		if v := tx.Bucket([]byte("Beers")).Get(itob(id)); v == nil {
			return zymurgauge.ErrNotFound
		} else if err := json.Unmarshal(v, &b); err != nil {
			return errors.Wrapf(err, "Could not unmarshal Beer %d", id)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &b, nil
}

// Save creates or updates a Beer
func (s *BeerService) Save(b *zymurgauge.Beer) error {
	err := s.db.Update(func(tx *bolt.Tx) error {

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
