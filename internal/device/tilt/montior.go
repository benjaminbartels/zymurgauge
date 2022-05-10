package tilt

import (
	"context"
	"math"
	"sync"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/api/beacon"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Color string

type Monitor struct {
	logger    *logrus.Logger
	tilts     map[Color]*Tilt
	colors    map[string]Color
	isRunning bool
	runMutex  sync.Mutex
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

	btAdapter, err := adapter.GetAdapter(adapter.GetDefaultAdapterID())
	if err != nil {
		return errors.Wrap(err, "could not get adapter")
	}

	if err := btAdapter.FlushDevices(); err != nil {
		return errors.Wrap(err, "could not flush devices")
	}

	filter := &adapter.DiscoveryFilter{}
	filter.Transport = "le"

	discovery, cancel, err := api.Discover(btAdapter, filter)
	if err != nil {
		return errors.Wrap(err, "could not start discovery")
	}

	for {
		select {
		case event := <-discovery:

			if event.Type == adapter.DeviceRemoved {
				for _, v := range m.tilts {
					if v.path == string(event.Path) {
						m.tilts[v.color].ibeacon = nil
						m.tilts[v.color].path = ""
					}
				}

				continue
			}

			d, err := device.NewDevice1(event.Path)
			if err != nil {
				m.logger.WithError(err).Errorf("Problem creating new device %s.", event.Path)
				continue
			}

			if d == nil {
				m.logger.Errorf("Device %s was nil.", event.Path)
				continue
			}

			go m.processDevice(ctx, d)

		case <-ctx.Done():
			m.runMutex.Lock()
			defer m.runMutex.Unlock()
			m.isRunning = false
			cancel()
			return nil
		}
	}
}

func (m *Monitor) processDevice(ctx context.Context, d *device.Device1) {
	b, err := beacon.NewBeacon(d)
	if err != nil {
		m.logger.WithError(err).Errorf("Problem creating new beacon from device %s.", d.Path)
	}

	beaconUpdated, err := b.WatchDeviceChanges(ctx)
	if err != nil {
		m.logger.WithError(err).Errorf("Problem watching for device changes on device %s.", d.Path)
	}

	for {
		isBeacon := <-beaconUpdated
		if !isBeacon {
			m.logger.Debugf("Device device %s is not a beacon.", d.Path)
			return
		}

		name := b.Device.Properties.Alias
		if name == "" {
			name = b.Device.Properties.Name
		}

		if b.IsIBeacon() {
			m.logger.Debugf("Found IBeacon Type: %s, Name: %s", b.Type, name)
			ibeacon := b.GetIBeacon()
			m.logger.Debugf("IBeacon %s (%ddbi) (major=%d minor=%d)\n",
				ibeacon.ProximityUUID, ibeacon.MeasuredPower, ibeacon.Major, ibeacon.Minor)

			if color, ok := m.colors[ibeacon.ProximityUUID]; ok {
				if tilt, ok := m.tilts[color]; ok {
					if tilt.ibeacon == nil {
						m.logger.Debugf("Tilt online - Color: %s, UUID: %s, Major: %d, Minor: %d, Power: %d",
							m.colors[ibeacon.ProximityUUID], ibeacon.ProximityUUID, ibeacon.Major, ibeacon.Minor,
							ibeacon.MeasuredPower)
					}

					m.tilts[color].ibeacon = &ibeacon
					m.tilts[color].path = string(d.Path())
				}
			}
		}
	}
}

func (m *Monitor) GetTilt(color Color) (*Tilt, error) {
	tilt, ok := m.tilts[color]
	if !ok {
		return nil, ErrNotFound
	}

	return tilt, nil
}

// ==========================================================================================

var ErrIBeaconIsNil = errors.New("underlying IBeacon is nil")

type Tilt struct {
	ibeacon *beacon.BeaconIBeacon
	path    string
	color   Color
}

func (t *Tilt) GetID() string {
	return string(t.color)
}

func (t *Tilt) GetTemperature() (float64, error) {
	if t.ibeacon == nil {
		return 0, ErrIBeaconIsNil
	}

	return math.Round(float64(t.ibeacon.Major-32)/1.8*100) / 100, nil
}

func (t *Tilt) GetGravity() (float64, error) {
	if t.ibeacon == nil {
		return 0, ErrIBeaconIsNil
	}

	return float64(t.ibeacon.Minor) / 1000, nil
}

// TODO: protect against 2 tilts with same color
