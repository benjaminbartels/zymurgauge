package tilt

import (
	"context"
	"sync"

	"github.com/benjaminbartels/zymurgauge/internal/platform/bluetooth/ibeacon"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Color string

type Monitor struct {
	ibeaconDiscoverer ibeacon.Discoverer
	logger            *logrus.Logger
	tilts             map[Color]*Tilt
	colors            map[string]Color
	isRunning         bool
	runMutex          sync.Mutex
	tiltMutex         sync.RWMutex
}

func NewMonitor(ibeaconDiscoverer ibeacon.Discoverer, logger *logrus.Logger) *Monitor {
	m := &Monitor{
		ibeaconDiscoverer: ibeaconDiscoverer,
		logger:            logger,
		tilts:             make(map[Color]*Tilt),
		colors: map[string]Color{
			"A495BB10C5B14B44B5121370F02D74DE": "red",
			"A495BB20C5B14B44B5121370F02D74DE": "green",
			"A495BB30C5B14B44B5121370F02D74DE": "black",
			"A495BB40C5B14B44B5121370F02D74DE": "purple",
			"A495BB50C5B14B44B5121370F02D74DE": "orange",
			"A495BB60C5B14B44B5121370F02D74DE": "blue",
			"A495BB70C5B14B44B5121370F02D74DE": "yellow",
			"A495BB80C5B14B44B5121370F02D74DE": "pink",
		},
	}

	for _, color := range m.colors {
		m.tilts[color] = &Tilt{color: color}
	}

	return m
}

func (m *Monitor) Run(ctx context.Context) error {
	m.runMutex.Lock()

	if m.isRunning {
		defer m.runMutex.Unlock()

		return ErrAlreadyRunning
	}

	m.isRunning = true

	m.runMutex.Unlock()

	return m.startCycle(ctx)
}

func (m *Monitor) startCycle(ctx context.Context) error {
	discovery, err := m.ibeaconDiscoverer.Discover(ctx)
	if err != nil {
		return errors.Wrap(err, "could not start discovery")
	}

	for {
		select {
		case event := <-discovery:
			if color, ok := m.colors[event.UUID]; ok {
				if event.Type == ibeacon.Offline {
					m.handleOffline(color)

					continue
				}

				m.handleOnline(color, event.IBeacon)
			} else {
				m.logger.Debugf("IBeacon %s is not a tilt.", event.IBeacon.ProximityUUID)
			}

		case <-ctx.Done():
			m.runMutex.Lock()
			defer m.runMutex.Unlock()
			m.isRunning = false

			return nil
		}
	}
}

func (m *Monitor) handleOnline(color Color, ibeacon ibeacon.IBeacon) {
	m.tiltMutex.Lock()
	defer m.tiltMutex.Unlock()

	if _, tiltFound := m.tilts[color]; tiltFound {
		m.logger.Debugf("Tilt online - Color: %s", color)
		m.tilts[color].ibeacon = &ibeacon
	}
}

func (m *Monitor) handleOffline(color Color) {
	m.tiltMutex.Lock()
	defer m.tiltMutex.Unlock()

	if _, tiltFound := m.tilts[color]; tiltFound {
		m.logger.Debugf("Tilt offline. - Color: %s", color)
		m.tilts[color].ibeacon = nil
	}
}

func (m *Monitor) GetTilt(color Color) (*Tilt, error) {
	m.tiltMutex.RLock()
	defer m.tiltMutex.RUnlock()

	tilt, ok := m.tilts[color]
	if !ok {
		return nil, ErrNotFound
	}

	return tilt, nil
}
