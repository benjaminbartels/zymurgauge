package bolt

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
	"github.com/orangesword/zymurgauge"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// check implementation at compile time
var _ zymurgauge.FermentationService = &FermentationService{}

// FermentationService represents a service for managing fermentations
type FermentationService struct {
	db     *bolt.DB
	logger *logrus.Logger
}

// Get returns a Fermentation by its ID
func (s *FermentationService) Get(id uint64) (*zymurgauge.Fermentation, error) {
	var f zymurgauge.Fermentation

	err := s.db.View(func(tx *bolt.Tx) error {
		if v := tx.Bucket([]byte("Fermentations")).Get(itob(id)); v == nil {
			return zymurgauge.ErrNotFound
		} else if err := json.Unmarshal(v, &f); err != nil {
			return errors.Wrapf(err, "Could not unmarshal Fermentation %d", id)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &f, nil
}

// Save creates or updates a Fermentation
func (s *FermentationService) Save(f *zymurgauge.Fermentation) error {
	err := s.db.Update(func(tx *bolt.Tx) error {

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

// ToDo: Move to its own service
// // LogEvent logs the given event for the given Fermentation by it FermentationID
// func (s *FermentationService) LogEvent(fermentationID uint64, event zymurgauge.FermentationEvent) error {
// 	tx, err := s.db.Begin(true)
// 	if err != nil {
// 		return err
// 	}

// 	bu := tx.Bucket([]byte("Fermentations"))

// 	defer func() { _ = tx.Rollback() }()

// 	v := bu.Get(itob(fermentationID))
// 	if v == nil {
// 		return zymurgauge.ErrNotFound
// 	}

// 	var f zymurgauge.Fermentation
// 	if err := json.Unmarshal(v, &f); err != nil {
// 		return err
// 	}

// 	f.Events = append(f.Events, event)

// 	if v, err := json.Marshal(f); err != nil {
// 		return err
// 	} else if err := bu.Put(itob(f.ID), v); err != nil {
// 		return err
// 	}

// 	return tx.Commit()

// }
