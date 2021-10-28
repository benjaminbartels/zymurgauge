package database

import (
	"encoding/json"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/chamber"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
)

const chamberBucket = "Chambers"

var _ internal.ChamberRepo = (*ChamberRepo)(nil)

// ChamberRepo represents a bbolt repository for managing Chambers.
type ChamberRepo struct {
	db *bbolt.DB
}

// NewChamberRepo returns a new Chamber repository using the given bbolt database. It also creates the Chambers
// bucket if it is not yet created on disk.
func NewChamberRepo(db *bbolt.DB) (*ChamberRepo, error) {
	tx, err := db.Begin(true)
	if err != nil {
		return nil, errors.Wrap(err, "could not begin transaction")
	}

	defer rollback(tx, &err)

	if _, err := tx.CreateBucketIfNotExists([]byte(chamberBucket)); err != nil {
		return nil, errors.Wrap(err, "could not create Chamber bucket")
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "could not commit transaction")
	}

	return &ChamberRepo{
		db: db,
	}, nil
}

// GetAll returns all Chambers.
func (r *ChamberRepo) GetAll() ([]chamber.Chamber, error) {
	chambers := []chamber.Chamber{}

	if err := r.db.View(func(tx *bbolt.Tx) error {
		err := tx.Bucket([]byte(chamberBucket)).ForEach(func(k, v []byte) error {
			var c chamber.Chamber
			if err := json.Unmarshal(v, &c); err != nil {
				return errors.Wrap(err, "could not unmarshal Chamber")
			}
			chambers = append(chambers, c)

			return nil
		})

		return errors.Wrap(err, "could not iterate over Chambers")
	}); err != nil {
		return nil, errors.Wrap(err, "could not execute view transaction")
	}

	return chambers, nil
}

// Get returns a Chamber by its ID.
func (r *ChamberRepo) Get(id string) (*chamber.Chamber, error) {
	var c *chamber.Chamber

	if err := r.db.View(func(tx *bbolt.Tx) error {
		if v := tx.Bucket([]byte(chamberBucket)).Get([]byte(id)); v != nil {
			if err := json.Unmarshal(v, &c); err != nil {
				return errors.Wrapf(err, "could not unmarshal Chamber %s", id)
			}
		}

		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "could not execute view transaction")
	}

	return c, nil
}

// Save creates or updates a Chamber.
func (r *ChamberRepo) Save(c *chamber.Chamber) error {

	if c.ID == "" {
		c.ID = uuid.NewString()
	}

	if err := r.db.Update(func(tx *bbolt.Tx) error {
		bu := tx.Bucket([]byte(chamberBucket))
		c.ModTime = time.Now()
		if v, err := json.Marshal(c); err != nil {
			return errors.Wrapf(err, "could not marshal Chamber %s", c.ID)
		} else if err := bu.Put([]byte(c.ID), v); err != nil {
			return errors.Wrapf(err, "could not put Chamber %s", c.ID)
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "could not execute update transaction")
	}

	return nil
}

// Delete permanently removes a Chamber.
func (r *ChamberRepo) Delete(id string) error {
	if err := r.db.Update(func(tx *bbolt.Tx) error {
		bu := tx.Bucket([]byte(chamberBucket))
		if err := bu.Delete([]byte(id)); err != nil {
			return errors.Wrapf(err, "could not delete Chamber %s", id)
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "could not execute update transaction")
	}

	return nil
}
