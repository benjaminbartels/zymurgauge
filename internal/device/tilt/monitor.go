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
	tiltID = "4c000215a495" // TODO: remove preamble
)

var (
	ErrAlreadyRunning = errors.New("monitor is already running")
	ErrNotFound       = errors.New("tilt not found is nil")
)

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
	timeout      time.Duration
	interval     time.Duration
	tilts        map[Color]*Tilt
	isRunning    bool
	runningMutex sync.Mutex
	logger       *logrus.Logger
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

func (m *Monitor) GetTilt(color Color) (*Tilt, error) {
	if m.tilts[color].ibeacon == nil {
		return nil, ErrNotFound
	}

	return m.tilts[color], nil
}

func (m *Monitor) Run(ctx context.Context) error {
	var discoveredIBeacons map[string]*IBeacon

	m.runningMutex.Lock()
	if m.isRunning {
		return ErrAlreadyRunning
	}

	device, err := linux.NewDevice()
	if err != nil {
		return errors.Wrap(err, "could not create new device")
	}

	ble.SetDefaultDevice(device)

	filter := func(adv ble.Advertisement) bool {
		return len(adv.ManufacturerData()) >= 25 && hex.EncodeToString(adv.ManufacturerData())[0:12] == tiltID
	}

	handler := func(adv ble.Advertisement) {
		ibeacon := NewIBeacon(adv.ManufacturerData())
		discoveredIBeacons[ibeacon.UUID] = ibeacon
	}

	ctx = ble.WithSigHandler(context.WithTimeout(ctx, m.timeout))

	for {
		discoveredIBeacons = make(map[string]*IBeacon)

		if err := ble.Scan(ctx, false, handler, filter); err != nil {
			return errors.Wrap(err, "could not scan")
		}

		for uuid, color := range colors {
			if ibeacon, ok := discoveredIBeacons[uuid]; ok {
				m.tilts[color].ibeacon = ibeacon
			} else {
				m.tilts[color].ibeacon = nil
			}
		}

		timer := time.NewTimer(m.interval)
		defer timer.Stop()

		select {
		case <-timer.C:

		case <-ctx.Done():
			m.runningMutex.Lock()
			defer m.runningMutex.Unlock()
			m.isRunning = false

			return nil
		}
	}
}
