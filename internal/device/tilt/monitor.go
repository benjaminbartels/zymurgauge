package tilt

import (
	"context"
	"encoding/hex"
	"sync"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	tiltID       = "4c000215a495" // TODO: remove preamble
	Red    Color = "Red"
	Green  Color = "Green"
	Black  Color = "Black"
	Purple Color = "Purple"
	Orange Color = "Orange"
	Blue   Color = "Blue"
	Yellow Color = "Yellow"
	Pink   Color = "Pink"
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
	timeout         time.Duration
	interval        time.Duration
	logger          *logrus.Logger
	availibleColors []Color
	tilts           map[Color]*Tilt
	isRunning       bool
	runMutex        sync.Mutex
}

func NewMonitor(timeout, interval time.Duration, logger *logrus.Logger) *Monitor {
	m := &Monitor{
		timeout:  timeout,
		interval: interval,
		tilts:    make(map[Color]*Tilt),
		logger:   logger,
	}

	for _, color := range colors {
		m.tilts[color] = &Tilt{}
	}

	return m
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

	device, err := linux.NewDevice()
	if err != nil {
		return errors.Wrap(err, "could not create new device")
	}

	ble.SetDefaultDevice(device)

	return m.startCycle(ctx)
}

func (m *Monitor) startCycle(ctx context.Context) error {
	for {
		m.availibleColors = []Color{}

		scanCtx := ble.WithSigHandler(context.WithTimeout(ctx, m.timeout))

		if err := ble.Scan(scanCtx, false, m.advHandler, m.advFilter); err != nil {
			switch {
			case errors.Is(err, context.DeadlineExceeded):
			case errors.Is(err, context.Canceled):
				return nil // TODO: should this return ctx.Err()
			default:
				return errors.Wrap(err, "could not scan")
			}
		}

		for _, color := range colors {
			if !containsColor(m.availibleColors, color) {
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

func (m *Monitor) advFilter(adv ble.Advertisement) bool {
	return len(adv.ManufacturerData()) >= 25 && hex.EncodeToString(adv.ManufacturerData())[0:12] == tiltID
}

func (m *Monitor) advHandler(adv ble.Advertisement) {
	ibeacon := NewIBeacon(adv.ManufacturerData())
	color := colors[ibeacon.UUID]
	m.availibleColors = append(m.availibleColors, color)
	m.tilts[color].ibeacon = ibeacon
	m.logger.Debugf("Discovered Tilt - Color: %s, UUID: %s, Major: %d, Minor: %d\n", colors[ibeacon.UUID],
		ibeacon.UUID, ibeacon.Major, ibeacon.Minor)
}

func containsColor(elems []Color, v Color) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}

	return false
}
