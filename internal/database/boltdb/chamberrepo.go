package boltdb

import (
	"encoding/json"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

// ChamberRepo represents a boltdb repository for managing Chambers
type ChamberRepo struct {
	db *bolt.DB
}

// NewChamberRepo returns a new Chamber repository using the given bolt boltdb. It also creates the Chambers
// bucket if it is not yet created on disk.
func NewChamberRepo(db *bolt.DB) (*ChamberRepo, error) {
	tx, err := db.Begin(true)
	if err != nil {
		return nil, errors.Wrap(err, "Could not begin transaction")
	}
	defer rollback(tx, &err)
	if _, err := tx.CreateBucketIfNotExists([]byte("Chambers")); err != nil {
		return nil, errors.Wrap(err, "Could not create Chamber bucket")
	}
	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "Could not commit transaction")
	}
	return &ChamberRepo{
		db: db,
	}, nil
}

// GetAll returns all Chambers
func (r *ChamberRepo) GetAll() ([]internal.Chamber, error) {
	var chambers []internal.Chamber
	err := r.db.View(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte("Chambers")).ForEach(func(k, v []byte) error {
			var c internal.Chamber
			if err := json.Unmarshal(v, &c); err != nil {
				return errors.Wrap(err, "Could not unmarshal Chambers")
			}
			chambers = append(chambers, c)
			return nil
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return chambers, nil
}

// Get returns a Chamber by its MAC address
func (r *ChamberRepo) Get(mac string) (*internal.Chamber, error) {
	var c *internal.Chamber
	err := r.db.View(func(tx *bolt.Tx) error {
		if v := tx.Bucket([]byte("Chambers")).Get([]byte(mac)); v != nil {
			if err := json.Unmarshal(v, &c); err != nil {
				return errors.Wrapf(err, "Could not unmarshal Chamber %s", mac)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Save creates or updates a Chamber
func (r *ChamberRepo) Save(c *internal.Chamber) error {
	err := r.db.Update(func(tx *bolt.Tx) error {
		bu := tx.Bucket([]byte("Chambers"))
		c.ModTime = time.Now()
		if v, err := json.Marshal(c); err != nil {
			return errors.Wrapf(err, "Could not marshal Chamber %s", c.MacAddress)
		} else if err := bu.Put([]byte(c.ID), v); err != nil {
			return errors.Wrapf(err, "Could not put Chamber %s", c.MacAddress)
		}
		return nil
	})
	return err
}

// Delete permanently removes a Chamber
func (r *ChamberRepo) Delete(mac string) error {
	err := r.db.Update(func(tx *bolt.Tx) error {
		bu := tx.Bucket([]byte("Chambers"))
		if err := bu.Delete([]byte(mac)); err != nil {
			return errors.Wrapf(err, "Could not delete Chamber %s", mac)
		}
		return nil
	})
	return err
}
