package tilt

import (
	"context"
	"encoding/hex"
	"sync"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/platform/bluetooth"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	tiltID                        = "4c000215a495" // TODO: remove preamble
	defaultTimeout  time.Duration = 10 * time.Second
	defaultInterval time.Duration = 1 * time.Second
)

type Color string

type Monitor struct {
	scanner         bluetooth.Scanner
	logger          *logrus.Logger
	timeout         time.Duration
	interval        time.Duration
	availableColors []Color
	tilts           map[Color]*Tilt
	isRunning       bool
	runMutex        sync.Mutex
	colors          map[string]Color
}

func NewMonitor(scanner bluetooth.Scanner, logger *logrus.Logger, options ...OptionsFunc) *Monitor {
	m := &Monitor{
		scanner:  scanner,
		timeout:  defaultTimeout,
		interval: defaultInterval,
		tilts:    make(map[Color]*Tilt),
		logger:   logger,
		colors: map[string]Color{
			"a495bb10c5b14b44b5121370f02d74de": "red",
			"a495bb20c5b14b44b5121370f02d74de": "green",
			"a495bb30c5b14b44b5121370f02d74de": "black",
			"a495bb40c5b14b44b5121370f02d74de": "purple",
			"a495bb50c5b14b44b5121370f02d74de": "orange",
			"a495bb60c5b14b44b5121370f02d74de": "blue",
			"a495bb70c5b14b44b5121370f02d74de": "yellow",
			"a495bb80c5b14b44b5121370f02d74de": "pink",
		},
	}

	for _, option := range options {
		option(m)
	}

	for _, color := range m.colors {
		m.tilts[color] = &Tilt{color: color}
	}

	return m
}

type OptionsFunc func(*Monitor)

func SetTimeout(timeout time.Duration) OptionsFunc {
	return func(t *Monitor) {
		t.timeout = timeout
	}
}

func SetInterval(interval time.Duration) OptionsFunc {
	return func(t *Monitor) {
		t.interval = interval
	}
}

func (m *Monitor) GetTilt(color Color) (*Tilt, error) { // TODO: check for race?
	tilt, ok := m.tilts[color]
	if !ok {
		return nil, ErrNotFound
	}

	return tilt, nil
}

func (m *Monitor) Run(ctx context.Context) error {
	m.runMutex.Lock()

	if m.isRunning {
		defer m.runMutex.Unlock()

		return ErrAlreadyRunning
	}

	m.isRunning = true

	m.runMutex.Unlock()

	if err := m.setupHCIDevice(); err != nil {
		return errors.Wrap(err, "could not setup hci device")
	}

	return m.startCycle(ctx)
}

func (m *Monitor) setupHCIDevice() error {
	device, err := m.scanner.NewDevice()
	if err != nil {
		return errors.Wrap(err, "could not create new device")
	}

	m.scanner.SetDefaultDevice(device)

	return nil
}

func (m *Monitor) startCycle(ctx context.Context) error {
	for {
		m.availableColors = []Color{}

		scanCtx := m.scanner.WithSigHandler(context.WithTimeout(ctx, m.timeout))

		if err := m.scanner.Scan(scanCtx, m.handler, m.filter); err != nil {
			switch {
			case errors.Is(err, context.DeadlineExceeded):
			case errors.Is(err, context.Canceled):
				return nil // TODO: should this return ctx.Err()
			default:
				m.logger.WithError(err).Warn("Error occurred while scanning. Resetting hci device.")

				if err := m.setupHCIDevice(); err != nil {
					return errors.Wrap(err, "could not setup hci device")
				}
			}
		}

		m.handleOfflineTilts()

		if m.wait(ctx) {
			return nil
		}
	}
}

func (m *Monitor) handleOfflineTilts() {
	for _, color := range m.colors {
		if !containsColor(m.availableColors, color) && m.tilts[color].ibeacon != nil {
			m.logger.Debugf("Tilt offline - Color: %s", color)
			m.tilts[color].ibeacon = nil
		}
	}
}

func (m *Monitor) wait(ctx context.Context) bool {
	timer := time.NewTimer(m.interval)
	defer timer.Stop()

	select {
	case <-timer.C:
	case <-ctx.Done():
		m.runMutex.Lock()
		defer m.runMutex.Unlock()
		m.isRunning = false

		return true
	}

	return false
}

func (m *Monitor) handler(adv bluetooth.Advertisement) {
	ibeacon, err := NewIBeacon(adv.ManufacturerData())
	if err != nil {
		m.logger.WithError(err).Error("could not create new IBeacon")

		return
	}

	color := m.colors[ibeacon.UUID]
	m.availableColors = append(m.availableColors, color)

	if m.tilts[color].ibeacon == nil {
		m.logger.Debugf("Tilt online - Color: %s, UUID: %s, Major: %d, Minor: %d", m.colors[ibeacon.UUID],
			ibeacon.UUID, ibeacon.Major, ibeacon.Minor)
	}

	m.tilts[color].ibeacon = ibeacon
}

func (m *Monitor) filter(adv bluetooth.Advertisement) bool {
	return len(adv.ManufacturerData()) >= 25 && hex.EncodeToString(adv.ManufacturerData())[0:12] == tiltID
}

func containsColor(colors []Color, v Color) bool {
	for _, s := range colors {
		if v == s {
			return true
		}
	}

	return false
}
