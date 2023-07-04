package ibeacon

import (
	"context"
	"sync"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/platform/bluetooth"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/api/beacon"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	ErrAlreadyDiscovering = bluetooth.Error("ibeacon discoverer is already running")
	WatchTimeout          = 180 * time.Second
)

var _ Discoverer = (*BluezDiscoverer)(nil)

type BluezDiscoverer struct {
	logger    *logrus.Logger
	ch        chan Event
	adapter   *adapter.Adapter1
	ibeacons  map[string]IBeacon
	isRunning bool
	runMutex  sync.Mutex
}

func NewDiscoverer(logger *logrus.Logger) (*BluezDiscoverer, error) {
	a, err := adapter.GetAdapter(adapter.GetDefaultAdapterID())
	if err != nil {
		return nil, errors.Wrap(err, "could not get adapter")
	}

	if err := a.FlushDevices(); err != nil {
		return nil, errors.Wrap(err, "could not flush devices")
	}

	discoverer := &BluezDiscoverer{
		logger:   logger,
		ch:       make(chan Event),
		adapter:  a,
		ibeacons: make(map[string]IBeacon),
	}

	return discoverer, nil
}

//nolint:funlen // TODO: Shorten
func (d *BluezDiscoverer) Discover(ctx context.Context) (chan Event, error) {
	if d.isRunning {
		defer d.runMutex.Unlock()

		return nil, ErrAlreadyDiscovering
	}

	d.isRunning = true

	d.runMutex.Unlock()

	filter := &adapter.DiscoveryFilter{}
	filter.Transport = "le"

	discoveryCh, cancel, err := api.Discover(d.adapter, filter)
	if err != nil {
		return nil, errors.Wrap(err, "could not start discovery")
	}

	d.logger.Debug("IBeacon discovery starting.")

	go func() {
		for {
			select {
			case event := <-discoveryCh:
				if event.Type == adapter.DeviceRemoved {
					if ibeacon, ok := d.ibeacons[string(event.Path)]; ok {
						d.logger.Debugf("Removing IBeacon %s.", ibeacon.ProximityUUID)
						delete(d.ibeacons, string(event.Path))
						d.ch <- Event{
							UUID: ibeacon.ProximityUUID,
							Type: Offline,
						}
					}

					continue
				}

				device, err := device.NewDevice1(event.Path)
				if err != nil {
					d.logger.WithError(err).Errorf("Problem creating new device for %s.", event.Path)

					continue
				}

				if d == nil {
					d.logger.Errorf("Device is nil for %s.", event.Path)

					continue
				}

				go d.processDevice(ctx, device)

			case <-ctx.Done():
				d.runMutex.Lock()
				defer d.runMutex.Unlock()
				d.isRunning = false

				cancel()
				close(d.ch)

				break
			}
		}
	}()

	return d.ch, nil
}

//nolint:funlen // TODO: Shorten
func (d *BluezDiscoverer) processDevice(ctx context.Context, device *device.Device1) {
	d.logger.Debugf("Processing device %s.", device.Properties.Address)

	b, err := beacon.NewBeacon(device)
	if err != nil {
		d.logger.WithError(err).Errorf("Problem creating new beacon from device %s.", device.Properties.Address)
	}

	defer func() {
		if err := b.UnwatchDeviceChanges(); err != nil {
			d.logger.WithError(err).Errorf("Problem un-watching for device changes on device %s.", device.Properties.Address)
		}
	}()

	propChanged, err := b.WatchDeviceChanges(ctx)
	if err != nil {
		d.logger.WithError(err).Errorf("Problem watching for device changes on device %s.", device.Properties.Address)
	}

	for {
		timer := time.NewTimer(WatchTimeout)
		defer timer.Stop()

		select {
		case isChanged := <-propChanged:
			if !isChanged {
				d.logger.Debugf("Device %s is not a beacon.", device.Properties.Address)

				return
			}

			if !b.IsIBeacon() {
				d.logger.Debugf("Device %s is not a ibeacon.", device.Properties.Address)

				return
			}

			ibeaconBeacon := b.GetIBeacon()
			d.logger.Debugf("Found IBeacon Address: %s, UUID: %s Major: %d, Minor: %d, Power: %d", device.Properties.Address,
				ibeaconBeacon.ProximityUUID, ibeaconBeacon.Major, ibeaconBeacon.Minor, ibeaconBeacon.MeasuredPower)

			ibeacon := IBeacon{
				ProximityUUID: ibeaconBeacon.ProximityUUID,
				Major:         ibeaconBeacon.Major,
				Minor:         ibeaconBeacon.Minor,
				MeasuredPower: ibeaconBeacon.MeasuredPower,
			}

			d.ibeacons[string(device.Path())] = ibeacon

			event := Event{
				UUID:    ibeacon.ProximityUUID,
				Type:    Online,
				IBeacon: ibeacon,
			}

			d.ch <- event

		case <-timer.C:
			d.logger.Debugf("Timed out waiting for property to change on %s.", device.Properties.Address)

			return

		case <-ctx.Done():
			return
		}
	}
}
