package chamber

import (
	"context"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var ErrNotFound = errors.New("chamber not found")

var _ Controller = (*Manager)(nil)

type Manager struct {
	repo         Repo
	chambers     sync.Map
	configurator Configurator
	logger       *logrus.Logger
}

func NewManager(repo Repo, configurator Configurator,
	logger *logrus.Logger) (*Manager, error) {
	m := &Manager{
		repo:         repo,
		configurator: configurator,
		logger:       logger,
	}

	chambers, err := m.repo.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "could not get all chambers from repository")
	}

	var errs error

	for i := range chambers {
		if err := chambers[i].Configure(configurator, logger); err != nil {
			errs = multierror.Append(errs,
				errors.Wrapf(err, "could not configure temperature controller for chamber %s", chambers[i].Name))
		}

		m.chambers.Store(chambers[i].ID, &chambers[i])
	}

	if err != nil {
		return m, errors.Wrap(errs, "could not configure all temperature controllers")
	}

	return m, nil
}

func (m *Manager) GetAll() ([]*Chamber, error) {
	chambers := []*Chamber{}

	m.chambers.Range(func(key, value interface{}) bool {
		spew.Dump(value)

		chambers = append(chambers, value.(*Chamber))

		return true
	})

	return chambers, nil
}

func (m *Manager) Get(id string) (*Chamber, error) {
	var chamber *Chamber

	value, ok := m.chambers.Load(id)
	if !ok {
		return chamber, nil
	}

	chamber, ok = value.(*Chamber)
	if !ok {
		return chamber, errors.Errorf("type assertion failed for chamber %s", id)
	}

	return chamber, nil
}

func (m *Manager) Save(chamber *Chamber) error {
	if err := m.repo.Save(chamber); err != nil {
		return errors.Wrap(err, "could not save chamber to repository")
	}

	m.chambers.Store(chamber.ID, chamber)

	if err := chamber.Configure(m.configurator, m.logger); err != nil {
		return errors.Wrap(err, "could not configure chamber")
	}

	return nil
}

func (m *Manager) Delete(id string) error {
	if err := m.repo.Delete(id); err != nil {
		return errors.Wrapf(err, "could not delete chamber %s from repository", id)
	}

	m.chambers.Delete(id)

	return nil
}

func (m *Manager) StartFermentation(ctx context.Context, chamberID string, step int) error {
	value, ok := m.chambers.Load(chamberID)
	if !ok {
		return ErrNotFound
	}

	chamber, ok := value.(*Chamber)
	if !ok {
		return errors.Errorf("type assertion failed for chamber %s", chamberID)
	}

	err := chamber.StartFermentation(ctx, step)
	if err != nil {
		return errors.Wrap(err, "could not start fermentation")
	}

	return nil
}

func (m *Manager) StopFermentation(chamberID string) error {
	value, ok := m.chambers.Load(chamberID)
	if !ok {
		return ErrNotFound
	}

	chamber, ok := value.(*Chamber)
	if !ok {
		return errors.Errorf("type assertion failed for chamber %s", chamberID)
	}

	err := chamber.StopFermentation()
	if err != nil {
		return errors.Wrap(err, "could not stop fermentation")
	}

	return nil
}
