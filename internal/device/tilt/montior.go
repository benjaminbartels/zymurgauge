package tilt

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"tinygo.org/x/bluetooth"
)

const (
	iBeaconCompanyID = 76 // Apple's IBeacon company Id
	tiltTLL          = 60 * time.Second
)

type Color string

type Monitor struct {
	logger    *logrus.Logger
	tilts     map[Color]*Tilt
	colors    map[string]Color
	isRunning bool
	runMutex  sync.Mutex
	tiltMutex sync.RWMutex
}

func NewMonitor(logger *logrus.Logger) *Monitor {
	m := &Monitor{
		logger: logger,
		tilts:  make(map[Color]*Tilt),
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

func (m *Monitor) removeExpiredTilts() {
	m.tiltMutex.Lock()
	defer m.tiltMutex.Unlock()

	for _, tilt := range m.tilts {
		if tilt.lastSeen.Before(time.Now().Add(-tiltTLL)) {
			m.logger.Debugf("Removing expired Tilt: %s", tilt.color)
			delete(m.tilts, tilt.color)
		}
	}
}

func (m *Monitor) startCycle(ctx context.Context) error {
	adapter := bluetooth.DefaultAdapter

	if err := adapter.Enable(); err != nil {
		return errors.Wrap(err, "could not enable bluetooth adapter")
	}

	if err := adapter.Scan(m.scan); err != nil {
		return errors.Wrap(err, "could not scan bluetooth adapter")
	}

	for {
		timer := time.NewTimer(tiltTLL)
		defer timer.Stop()

		select {
		case <-timer.C:
			m.removeExpiredTilts()

		case <-ctx.Done():
			m.runMutex.Lock()

			if err := adapter.StopScan(); err != nil {
				return errors.Wrap(err, "could not stop scaning bluetooth adapter")
			}

			defer m.runMutex.Unlock()
			m.isRunning = false

			return nil
		}
	}
}

func (m *Monitor) scan(_ *bluetooth.Adapter, device bluetooth.ScanResult) {
	manufacturerData := device.ManufacturerData()

	if len(manufacturerData) > 0 {
		for _, element := range manufacturerData {
			if element.CompanyID == iBeaconCompanyID && len(element.Data) == 23 {
				uuid := hex.EncodeToString(element.Data[2:18])

				if color, ok := m.colors[uuid]; ok {
					major := binary.BigEndian.Uint16(element.Data[18:20])
					minor := binary.BigEndian.Uint16(element.Data[20:22])
					measuredPower := int8(element.Data[22])

					m.logger.Debugf("Tilt Online: Color: %s UUID: %s Major: %d, Minor: %d, Power: %d dBm",
						color, uuid, major, minor, measuredPower)

					m.tiltMutex.Lock()
					m.tilts[color] = &Tilt{
						color:       color,
						temperature: float64(major),
						gravity:     float64(minor),
						lastSeen:    time.Now(),
					}
					m.tiltMutex.Unlock()
				}
			}
		}
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
