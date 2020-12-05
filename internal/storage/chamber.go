package storage

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
)

const chamberBucket = "Chambers"

// Chamber represents an insulated box (fridge) with internal heating/cooling elements that reacts to changes in
// monitored temperatures, by correcting small deviations from your desired fermentation temperature.
type Chamber struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	ThermometerID string    `json:"thermometerID"`
	ChillerPIN    string    `json:"chillerPIN"`
	HeaterPIN     string    `json:"heaterPIN"`
	ChillerKp     float64   `json:"chillerKp"`
	ChillerKi     float64   `json:"chillerKi"`
	ChillerKd     float64   `json:"chillerKd"`
	HeaterKp      float64   `json:"heaterKp"`
	HeaterKi      float64   `json:"heaterKi"`
	HeaterKd      float64   `json:"heaterKd"`
	ModTime       time.Time `json:"modTime"`
}

// ChamberRepo represents a bboltdb repository for managing Chambers.
type ChamberRepo struct {
	db *bbolt.DB
}

// NewChamberRepo returns a new Chamber repository using the given bbolt database. It also creates the Chambers
// bucket if it is not yet created on disk.
func NewChamberRepo(db *bbolt.DB) (*ChamberRepo, error) {
	tx, err := db.Begin(true)
	if err != nil {
		return nil, errors.Wrap(err, "Could not begin transaction")
	}

	defer rollback(tx, &err)

	if _, err := tx.CreateBucketIfNotExists([]byte(chamberBucket)); err != nil {
		return nil, errors.Wrap(err, "Could not create Chamber bucket")
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "Could not commit transaction")
	}

	return &ChamberRepo{
		db: db,
	}, nil
}

// GetAll returns all Chambers.
func (r *ChamberRepo) GetAll() ([]Chamber, error) {
	var chambers []Chamber

	err := r.db.View(func(tx *bbolt.Tx) error {
		err := tx.Bucket([]byte(chamberBucket)).ForEach(func(k, v []byte) error {
			var c Chamber
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

// Get returns a Chamber by its ID.
func (r *ChamberRepo) Get(id string) (*Chamber, error) {
	var c *Chamber

	err := r.db.View(func(tx *bbolt.Tx) error {
		if v := tx.Bucket([]byte(chamberBucket)).Get([]byte(id)); v != nil {
			if err := json.Unmarshal(v, &c); err != nil {
				return errors.Wrapf(err, "Could not unmarshal Chamber %s", id)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Save creates or updates a Chamber.
func (r *ChamberRepo) Save(c *Chamber) error {
	err := r.db.Update(func(tx *bbolt.Tx) error {
		bu := tx.Bucket([]byte(chamberBucket))
		c.ModTime = time.Now()
		if v, err := json.Marshal(c); err != nil {
			return errors.Wrapf(err, "Could not marshal Chamber %s", c.ID)
		} else if err := bu.Put([]byte(c.ID), v); err != nil {
			return errors.Wrapf(err, "Could not put Chamber %s", c.ID)
		}

		return nil
	})

	return err
}

// Delete permanently removes a Chamber.
func (r *ChamberRepo) Delete(id string) error {
	err := r.db.Update(func(tx *bbolt.Tx) error {
		bu := tx.Bucket([]byte(chamberBucket))
		if err := bu.Delete([]byte(id)); err != nil {
			return errors.Wrapf(err, "Could not delete Chamber %s", id)
		}

		return nil
	})

	return err
}
