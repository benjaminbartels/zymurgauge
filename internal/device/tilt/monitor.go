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
	defaultTimeout  time.Duration = 5 * time.Second
	defaultInterval time.Duration = 1 * time.Second
	Red             Color         = "Red"
	Green           Color         = "Green"
	Black           Color         = "Black"
	Purple          Color         = "Purple"
	Orange          Color         = "Orange"
	Blue            Color         = "Blue"
	Yellow          Color         = "Yellow"
	Pink            Color         = "Pink"
)

var (
	ErrAlreadyRunning = errors.New("monitor is already running")
	ErrNotFound       = errors.New("tilt not found")
)

//nolint:gochecknoglobals
var colors = map[string]Color{
	"a495bb10c5b14b44b5121370f02d74de": "Red",
	"a495bb20c5b14b44b5121370f02d74de": "Green",
	"a495bb30c5b14b44b5121370f02d74de": "Black",
	"a495bb40c5b14b44b5121370f02d74de": "Purple",
	"a495bb50c5b14b44b5121370f02d74de": "Orange",
	"a495bb60c5b14b44b5121370f02d74de": "Blue",
	"a495bb70c5b14b44b5121370f02d74de": "Yellow",
	"a495bb80c5b14b44b5121370f02d74de": "Pink",
}

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
}

func NewMonitor(scanner bluetooth.Scanner, logger *logrus.Logger, options ...OptionsFunc) *Monitor {
	m := &Monitor{
		scanner:  scanner,
		timeout:  defaultTimeout,
		interval: defaultInterval,
		tilts:    make(map[Color]*Tilt),
		logger:   logger,
	}

	for _, option := range options {
		option(m)
	}

	for _, color := range colors {
		m.tilts[color] = &Tilt{}
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
	if m.tilts[color].ibeacon == nil {
		return nil, ErrNotFound
	}

	return m.tilts[color], nil
}

func (m *Monitor) Run(ctx context.Context) error {
	m.runMutex.Lock()
	if m.isRunning {
		return ErrAlreadyRunning
	}

	m.isRunning = true

	m.runMutex.Unlock()

	device, err := m.scanner.NewDevice()
	if err != nil {
		return errors.Wrap(err, "could not create new device")
	}

	m.scanner.SetDefaultDevice(device)

	return m.startCycle(ctx)
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
				return errors.Wrap(err, "could not scan")
			}
		}

		for _, color := range colors {
			if !containsColor(m.availableColors, color) {
				m.tilts[color].ibeacon = nil
			}
		}

		timer := time.NewTimer(m.interval)
		defer timer.Stop()

		select {
		case <-timer.C:
		case <-ctx.Done():
			m.runMutex.Lock()
			defer m.runMutex.Unlock()
			m.isRunning = false

			return nil // TODO: should this return ctx.Err()
		}
	}
}

func (m *Monitor) handler(adv bluetooth.Advertisement) {
	ibeacon, err := NewIBeacon(adv.ManufacturerData())
	if err != nil {
		m.logger.WithError(err).Error("could not create new IBeacon")

		return
	}

	color := colors[ibeacon.UUID]
	m.availableColors = append(m.availableColors, color)
	m.tilts[color].ibeacon = ibeacon
	m.logger.Debugf("Discovered Tilt - Color: %s, UUID: %s, Major: %d, Minor: %d\n", colors[ibeacon.UUID],
		ibeacon.UUID, ibeacon.Major, ibeacon.Minor)
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
