package bluetooth

import (
	"context"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"github.com/pkg/errors"
)

var _ Scanner = (*BLEScanner)(nil)

type BLEScanner struct {
	device *linux.Device
}

func NewBLEScanner() (*BLEScanner, error) {
	device, err := linux.NewDevice()
	if err != nil {
		return nil, errors.Wrap(err, "could not create new ble device")
	}

	ble.SetDefaultDevice(device)

	return &BLEScanner{device: device}, nil
}

func (b *BLEScanner) WithSigHandler(ctx context.Context, cancel func()) context.Context {
	return ble.WithSigHandler(ctx, cancel)
}

func (b *BLEScanner) Scan(ctx context.Context, h func(a Advertisement), f func(a Advertisement) bool) error {
	if err := b.scan(ctx, h, f); err != nil {
		if err := b.restart(); err != nil {
			return errors.Wrapf(err, "could not restart device after error: %s")
		}

		// retry after restart
		if err := b.scan(ctx, h, f); err != nil {
			return errors.Wrap(err, "could not scan after restart")
		}
	}

	return nil
}

func (b *BLEScanner) scan(ctx context.Context, h func(a Advertisement), f func(a Advertisement) bool) error {
	handler := func(adv ble.Advertisement) {
		h(adv)
	}

	filter := func(adv ble.Advertisement) bool {
		return f(adv)
	}

	if err := ble.Scan(ctx, false, handler, filter); err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
		case errors.Is(err, context.Canceled):
			return err
		default:
			return errors.Wrap(err, "could not scan")
		}
	}

	return nil
}

func (b *BLEScanner) restart() error {
	if err := b.device.Stop(); err != nil {
		return errors.Wrap(err, "could not stop device")
	}

	device, err := linux.NewDevice()
	if err != nil {
		return errors.Wrap(err, "could not create new device")
	}
	ble.SetDefaultDevice(device)
	b.device = device

	return nil
}
