package boltdb

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

// NewFermentationRepo returns a new Fermentation repository using the given bolt boltdb. It also creates the
// Fermentations bucket if it is not yet created on disk.
func NewFermentationRepo(db *bolt.DB) (*FermentationRepo, error) {
	tx, err := db.Begin(true)
	if err != nil {
		return nil, errors.Wrap(err, "Could not begin transaction")
	}
	defer rollback(tx, &err)
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
func (r *FermentationRepo) Get(id string) (*internal.Fermentation, error) {
	var f *internal.Fermentation
	err := r.db.View(func(tx *bolt.Tx) error {
		if v := tx.Bucket([]byte("Fermentations")).Get([]byte(id)); v != nil {
			if err := json.Unmarshal(v, &f); err != nil {
				return errors.Wrapf(err, "Could not unmarshal Fermentation %s", id)
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
	fermentations := []internal.Fermentation{} // ToDo: init all array this way in repos
	err := r.db.View(func(tx *bolt.Tx) error {
		bu := tx.Bucket([]byte("Fermentations"))
		return bu.ForEach(func(k, v []byte) error {
			var f internal.Fermentation
			if err := json.Unmarshal(v, &f); err != nil {
				return err
			}
			fermentations = append(fermentations, f)
			return nil
		})

	})

	return fermentations, err
}

// Save creates or updates a Fermentation
func (r *FermentationRepo) Save(f *internal.Fermentation) error {
	err := r.db.Update(func(tx *bolt.Tx) error {
		bu := tx.Bucket([]byte("Fermentations"))
		f.ModTime = time.Now()
		if v, err := json.Marshal(f); err != nil {
			return errors.Wrapf(err, "Could not marshal Fermentation %s", f.ID)
		} else if err := bu.Put([]byte(f.ID), v); err != nil {
			return errors.Wrapf(err, "Could not put Fermentation %s", f.ID)
		}
		return nil
	})
	return err
}

// Delete permanently removes a Fermentation
func (r *FermentationRepo) Delete(id string) error {
	err := r.db.Update(func(tx *bolt.Tx) error {
		bu := tx.Bucket([]byte("Fermentations"))
		if err := bu.Delete([]byte(id)); err != nil {
			return errors.Wrapf(err, "Could not delete Fermentation %s", id)
		}
		return nil
	})
	return err
}
