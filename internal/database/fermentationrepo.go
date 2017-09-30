package database

import (
	"encoding/json"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

// FermentationRepo represents a boltdb repository for managing Fermentations
type FermentationRepo struct {
	db *bolt.DB
}

// NewFermentationRepo returns a new Fermentation repository using the given bolt database. It also creates the
// Fermentations bucket if it is not yet created on disk.
func NewFermentationRepo(db *bolt.DB) (*FermentationRepo, error) {

	tx, err := db.Begin(true)
	if err != nil {
		return nil, errors.Wrap(err, "Could not begin transaction")
	}

	defer tx.Rollback()

	if _, err := tx.CreateBucketIfNotExists([]byte("Fermentations")); err != nil {
		return nil, errors.Wrap(err, "Could not create Fermentation bucket")
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "Could not commit transaction")
	}

	return &FermentationRepo{
		db: db,
	}, nil
}

// Get returns a Fermentation by its ID
func (r *FermentationRepo) Get(id uint64) (*internal.Fermentation, error) {
	var f *internal.Fermentation

	err := r.db.View(func(tx *bolt.Tx) error {
		if v := tx.Bucket([]byte("Fermentations")).Get(itob(id)); v != nil {
			if err := json.Unmarshal(v, &f); err != nil {
				return errors.Wrapf(err, "Could not unmarshal Fermentation %d", id)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return f, nil
}

// GetAll returns all Fermentations
func (r *FermentationRepo) GetAll() ([]internal.Fermentation, error) {
	var fermentations []internal.Fermentation

	err := r.db.View(func(tx *bolt.Tx) error {

		v := tx.Bucket([]byte("Fermentations"))
		return v.ForEach(func(k, v []byte) error {
			var f internal.Fermentation
			if err := json.Unmarshal(v, &f); err != nil {
				return err
			}
			fermentations = append(fermentations, f)
			return nil
		})

	})

	if err != nil {
		return []internal.Fermentation{}, err
	}
	return fermentations, nil
}

// Save creates or updates a Fermentation
func (r *FermentationRepo) Save(f *internal.Fermentation) error {
	err := r.db.Update(func(tx *bolt.Tx) error {

		bu := tx.Bucket([]byte("Fermentations"))

		if v := bu.Get(itob(f.ID)); v == nil {
			seq, _ := bu.NextSequence()
			f.ID = seq
		}

		f.ModTime = time.Now()

		if v, err := json.Marshal(f); err != nil {
			return errors.Wrapf(err, "Could not marshal Fermentation %d", f.ID)
		} else if err := bu.Put(itob(f.ID), v); err != nil {
			return errors.Wrapf(err, "Could not put Fermentation %d", f.ID)
		}
		return nil
	})
	return err
}
