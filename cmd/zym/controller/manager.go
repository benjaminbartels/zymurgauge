package controller

import (
	"context"
	"sync"

	"github.com/benjaminbartels/zymurgauge/internal/chamber"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// var ErrNotFound = errors.New("chamber not found")

type ChamberManager struct {
	repo     chamber.Repo
	chambers sync.Map
	// chambers map[string]*chamber.Chamber
	mainCtx context.Context
	logger  *logrus.Logger
}

func NewChamberManager(ctx context.Context, repo chamber.Repo, logger *logrus.Logger) (*ChamberManager, error) {
	c := &ChamberManager{
		repo: repo,
		// chambers: make(map[string]*chamber.Chamber),
		mainCtx: ctx,
		logger:  logger,
	}

	chambers, err := c.repo.GetAllChambers()
	if err != nil {
		return nil, errors.Wrap(err, "could not get all chambers from repository")
	}

	var errs error

	for i := range chambers {
		if err := chambers[i].Configure(c.mainCtx, logger); err != nil {
			errs = multierror.Append(errs,
				errors.Wrapf(err, "could not configure temperature controller for chamber %s", chambers[i].Name))
		}

		c.chambers.Store(chambers[i].ID, &chambers[i])
	}

	if err != nil {
		return c, errors.Wrap(errs, "could not configure all temperature controllers")
	}

	return c, nil
}

func (c *ChamberManager) GetAllChambers() []*chamber.Chamber {
	chambers := []*chamber.Chamber{}

	c.chambers.Range(func(key, value interface{}) bool {
		chambers = append(chambers, value.(*chamber.Chamber))

		return true
	})

	return chambers
}

func (c *ChamberManager) GetChamber(id string) *chamber.Chamber {
	value, ok := c.chambers.Load(id)
	if ok {
		return value.(*chamber.Chamber)
	}

	return nil
}

func (c *ChamberManager) SaveChamber(chamber *chamber.Chamber) error {
	if err := c.repo.SaveChamber(chamber); err != nil {
		return errors.Wrap(err, "could not save chamber to repository")
	}

	c.chambers.Store(chamber.ID, chamber)

	if err := chamber.Configure(c.mainCtx, c.logger); err != nil {
		return errors.Wrap(err, "could not configure chamber")
	}

	return nil
}

func (c *ChamberManager) DeleteChamber(id string) error {
	if err := c.repo.DeleteChamber(id); err != nil {
		return errors.Wrapf(err, "could not delete chamber %s from repository", id)
	}

	c.chambers.Delete(id)

	return nil
}
