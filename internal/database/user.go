package database

import (
	"encoding/json"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/auth"
	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
)

const (
	userBucket   = "Users"
	adminUserKey = "admin"
)

var _ auth.UserRepo = (*UserRepo)(nil)

// UserRepo represents a bbolt repository for managing User.
type UserRepo struct {
	db *bbolt.DB
}

// NewUserRepo returns a new User repository using the given bbolt database. It also creates the User
// bucket if it is not yet created on disk.
func NewUserRepo(db *bbolt.DB) (*UserRepo, error) {
	tx, err := db.Begin(true)
	if err != nil {
		return nil, errors.Wrap(err, "could not begin transaction")
	}

	defer rollback(tx, &err)

	if _, err := tx.CreateBucketIfNotExists([]byte(userBucket)); err != nil {
		return nil, errors.Wrap(err, "could not create User bucket")
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "could not commit transaction")
	}

	return &UserRepo{
		db: db,
	}, nil
}

// Get returns a User by its ID.
func (r *UserRepo) Get() (*auth.User, error) {
	var u *auth.User

	if err := r.db.View(func(tx *bbolt.Tx) error {
		if v := tx.Bucket([]byte(settingsBucket)).Get([]byte(adminUserKey)); v != nil {
			if err := json.Unmarshal(v, &u); err != nil {
				return errors.Wrap(err, "could not unmarshal Settings")
			}
		}

		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "could not execute view transaction")
	}

	return u, nil
}

// Save creates or updates Settings.
func (r *UserRepo) Save(c *auth.User) error {
	c.ModTime = time.Now().UTC()

	if err := r.db.Update(func(tx *bbolt.Tx) error {
		bu := tx.Bucket([]byte(settingsBucket))
		c.ModTime = time.Now()
		if v, err := json.Marshal(c); err != nil {
			return errors.Wrap(err, "could not marshal Settings")
		} else if err := bu.Put([]byte(adminUserKey), v); err != nil {
			return errors.Wrap(err, "could not put Settings")
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "could not execute update transaction")
	}

	return nil
}
