package database

import (
	"encoding/json"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/settings"
	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
)

const (
	settingsBucket = "Settings"
	settingsKey    = "settings"
)

var _ settings.Repo = (*SettingsRepo)(nil)

// SettingsRepo represents a bbolt repository for managing Settings.
type SettingsRepo struct {
	db *bbolt.DB
}

// NewSettingsRepo returns a new Settings repository using the given bbolt database. It also creates the Settings
// bucket if it is not yet created on disk.
func NewSettingsRepo(db *bbolt.DB) (*SettingsRepo, error) {
	tx, err := db.Begin(true)
	if err != nil {
		return nil, errors.Wrap(err, "could not begin transaction")
	}

	defer rollback(tx, &err)

	if _, err := tx.CreateBucketIfNotExists([]byte(settingsBucket)); err != nil {
		return nil, errors.Wrap(err, "could not create Settings bucket")
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "could not commit transaction")
	}

	return &SettingsRepo{
		db: db,
	}, nil
}

// Get returns Settings.
func (r *SettingsRepo) Get() (*settings.Settings, error) {
	var c *settings.Settings

	if err := r.db.View(func(tx *bbolt.Tx) error {
		if v := tx.Bucket([]byte(settingsBucket)).Get([]byte(settingsKey)); v != nil {
			if err := json.Unmarshal(v, &c); err != nil {
				return errors.Wrap(err, "could not unmarshal Settings")
			}
		}

		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "could not execute view transaction")
	}

	return c, nil
}

// Save creates or updates Settings.
func (r *SettingsRepo) Save(c *settings.Settings) error {
	c.ModTime = time.Now().UTC()

	if err := r.db.Update(func(tx *bbolt.Tx) error {
		bu := tx.Bucket([]byte(settingsBucket))
		c.ModTime = time.Now()
		if v, err := json.Marshal(c); err != nil {
			return errors.Wrap(err, "could not marshal Settings")
		} else if err := bu.Put([]byte(settingsKey), v); err != nil {
			return errors.Wrap(err, "could not put Settings")
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "could not execute update transaction")
	}

	return nil
}
