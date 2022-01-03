package chamber

import (
	"context"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	ErrNotFound   = errors.New("chamber not found")
	ErrFermenting = errors.New("fermentation has started")
)

var _ Controller = (*Manager)(nil)

type Manager struct {
	repo         Repo
	chambers     map[string]*Chamber
	configurator Configurator
	logger       *logrus.Logger
	mutex        sync.RWMutex
}

func NewManager(repo Repo, configurator Configurator,
	logger *logrus.Logger) (*Manager, error) {
	m := &Manager{
		repo:         repo,
		chambers:     make(map[string]*Chamber),
		configurator: configurator,
		logger:       logger,
	}

	chambers, err := m.repo.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "could not get all chambers from repository")
	}

	var errs error

	m.mutex.Lock()
	defer m.mutex.Unlock()

	for i := range chambers {
		// TODO: Configure implementation should vary based on arch
		if err := chambers[i].Configure(configurator, logger); err != nil {
			errs = multierror.Append(errs,
				errors.Wrapf(err, "could not configure temperature controller for chamber %s", chambers[i].Name))
		}

		m.chambers[chambers[i].ID] = chambers[i]
	}

	if errs != nil {
		return m, errors.Wrap(errs, "could not configure all temperature controllers")
	}

	return m, nil
}

func (m *Manager) GetAll() ([]*Chamber, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	chambers := make([]*Chamber, 0, len(m.chambers))
	for _, chamber := range m.chambers {
		chambers = append(chambers, chamber)
	}

	// It is not possible for GetAll() to return an error
	return chambers, nil
}

func (m *Manager) Get(id string) (*Chamber, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	chamber, ok := m.chambers[id]
	if !ok {
		return chamber, nil
	}

	return chamber, nil
}

func (m *Manager) Save(chamber *Chamber) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if c, ok := m.chambers[chamber.ID]; ok && c.IsFermenting() {
		return ErrFermenting
	}

	m.chambers[chamber.ID] = chamber

	if err := m.repo.Save(chamber); err != nil {
		return errors.Wrap(err, "could not save chamber to repository")
	}

	if err := chamber.Configure(m.configurator, m.logger); err != nil {
		return errors.Wrap(err, "could not configure chamber")
	}

	return nil
}

func (m *Manager) Delete(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if c, ok := m.chambers[id]; ok && c.IsFermenting() {
		return ErrFermenting
	}

	if err := m.repo.Delete(id); err != nil {
		return errors.Wrapf(err, "could not delete chamber %s from repository", id)
	}

	delete(m.chambers, id)

	return nil
}

func (m *Manager) StartFermentation(ctx context.Context, chamberID string, step int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	chamber, ok := m.chambers[chamberID]
	if !ok {
		return ErrNotFound
	}

	err := chamber.StartFermentation(ctx, step)
	if err != nil {
		return errors.Wrap(err, "could not start fermentation")
	}

	return nil
}

func (m *Manager) StopFermentation(chamberID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	chamber, ok := m.chambers[chamberID]
	if !ok {
		return ErrNotFound
	}

	err := chamber.StopFermentation()
	if err != nil {
		return errors.Wrap(err, "could not stop fermentation")
	}

	return nil
}