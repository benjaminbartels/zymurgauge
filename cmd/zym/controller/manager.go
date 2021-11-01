package controller

import (
	"context"

	"github.com/benjaminbartels/zymurgauge/internal/chamber"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	ErrNotFound       = errors.New("chamber not found")
	ErrNoCurrentBatch = errors.New("chamber does not have a current batch")
)

type ChamberController interface {
	chamber.Repo
	StartFermentation(chamberID string, step int) error
	StopFermentation(chamberID string) error
}

type ChamberManager struct {
	repo     chamber.Repo
	chambers map[string]*chamber.Chamber
	mainCtx  context.Context
	logger   *logrus.Logger
}

func NewChamberManager(ctx context.Context, repo chamber.Repo, logger *logrus.Logger) (*ChamberManager, error) {
	c := &ChamberManager{
		repo:     repo,
		chambers: make(map[string]*chamber.Chamber),
		mainCtx:  ctx,
		logger:   logger,
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

		c.chambers[chambers[i].ID] = &chambers[i]
	}

	if err != nil {
		return c, errors.Wrap(errs, "could not configure all temperature controllers")
	}

	return c, nil
}

func (c *ChamberManager) GetAllChambers() ([]chamber.Chamber, error) {
	chambers := make([]chamber.Chamber, 0, len(c.chambers))

	for _, chamber := range c.chambers {
		chambers = append(chambers, *chamber)
	}

	return chambers, nil
}

func (c *ChamberManager) GetChamber(id string) (*chamber.Chamber, error) {
	return c.chambers[id], nil
}

func (c *ChamberManager) SaveChamber(chamber *chamber.Chamber) error {
	if err := c.repo.SaveChamber(chamber); err != nil {
		return errors.Wrap(err, "could not save chamber to repository")
	}

	c.chambers[chamber.ID] = chamber

	if err := c.chambers[chamber.ID].Configure(c.mainCtx, c.logger); err != nil {
		return errors.Wrap(err, "could not save chamber to repository")
	}

	return nil
}

func (c *ChamberManager) DeleteChamber(id string) error {
	if err := c.repo.DeleteChamber(id); err != nil {
		return errors.Wrapf(err, "could not delete chamber %s from repository", id)
	}

	delete(c.chambers, id)

	return nil
}

func (c *ChamberManager) StartFermentation(chamberID string, step int) error {
	chamber, ok := c.chambers[chamberID]
	if !ok {
		return ErrNotFound
	}

	if chamber.CurrentBatch == nil {
		return ErrNoCurrentBatch
	}

	if err := chamber.StartFermentation(step); err != nil {
		return errors.Wrap(err, "could not start fermentation")
	}

	chamber.CurrentFermentationStep = step

	if err := c.repo.SaveChamber(chamber); err != nil {
		return errors.Wrap(err, "could not save chamber to repository")
	}

	return nil
}

func (c *ChamberManager) StopFermentation(chamberID string) error {
	if err := c.chambers[chamberID].StopFermentation(); err != nil {
		return errors.Wrap(err, "could not stop fermentation")
	}

	return nil
}
