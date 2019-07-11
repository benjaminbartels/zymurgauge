package boltdb

import (
	"encoding/json"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

// TemperatureChangeRepo represents a boltdb repository for managing fermentationTemperatureChanges
type TemperatureChangeRepo struct {
	db *bolt.DB
}

// NewTemperatureChangeRepo returns a new TemperatureChange repository using the given bolt boltdb. It also creates
// the TemperatureChanges bucket if it is not yet created on disk.
func NewTemperatureChangeRepo(db *bolt.DB) (*TemperatureChangeRepo, error) {
	tx, err := db.Begin(true)
	if err != nil {
		return nil, errors.Wrap(err, "Could not begin transaction")
	}
	defer rollback(tx, &err)
	if _, err := tx.CreateBucketIfNotExists([]byte("TemperatureChanges")); err != nil {
		return nil, errors.Wrap(err, "Could not create TemperatureChange bucket")
	}
	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "Could not commit transaction")
	}
	return &TemperatureChangeRepo{
		db: db,
	}, nil
}

// GetRangeByFermentationID returns a all temperature changes for the given fermentation id for the given range
func (r *TemperatureChangeRepo) GetRangeByFermentationID(fermentationID string, start,
	end time.Time) ([]internal.TemperatureChange, error) {
	temperatureChanges := []internal.TemperatureChange{}

	err := r.db.View(func(tx *bolt.Tx) error {
		bu := tx.Bucket([]byte("TemperatureChanges"))
		return bu.ForEach(func(k, v []byte) error {
			var f internal.TemperatureChange
			if err := json.Unmarshal(v, &f); err != nil {
				return err
			}
			if f.FermentationID == fermentationID && f.InsertTime.After(start) && f.InsertTime.Before(end) {
				temperatureChanges = append(temperatureChanges, f)
			}
			return nil
		})

	})

	return temperatureChanges, err
}

// Save creates or updates a TemperatureChange
func (r *TemperatureChangeRepo) Save(b *internal.TemperatureChange) error {
	err := r.db.Update(func(tx *bolt.Tx) error {
		bu := tx.Bucket([]byte("TemperatureChanges"))
		if v, err := json.Marshal(b); err != nil {
			return errors.Wrapf(err, "Could not marshal TemperatureChange %s", b.ID)
		} else if err := bu.Put([]byte(b.ID), v); err != nil {
			return errors.Wrapf(err, "Could not put TemperatureChange %s", b.ID)
		}
		return nil
	})
	return err
}

// ToDo: Delete all by FermentationID
