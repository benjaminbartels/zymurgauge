package bluetooth

import (
	"context"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"github.com/pkg/errors"
)

var _ Scanner = (*BLEScanner)(nil)

type BLEScanner struct{}

func NewBLEScanner() *BLEScanner {
	return &BLEScanner{}
}

func (b *BLEScanner) NewDevice() (*linux.Device, error) {
	device, err := linux.NewDevice()
	if err != nil {
		return nil, errors.Wrap(err, "could not create new device")
	}

	return device, nil
}

func (b *BLEScanner) SetDefaultDevice(device Device) {
	ble.SetDefaultDevice(device)
}

func (b *BLEScanner) WithSigHandler(ctx context.Context, cancel func()) context.Context {
	return ble.WithSigHandler(ctx, cancel)
}

func (b *BLEScanner) Scan(ctx context.Context, h func(a Advertisement), f func(a Advertisement) bool) error {
	handler := func(adv ble.Advertisement) {
		h(adv)
	}

	filter := func(adv ble.Advertisement) bool {
		return f(adv)
	}

	if err := ble.Scan(ctx, false, handler, filter); err != nil {
		return errors.Wrap(err, "could not scan")
	}

	return nil
}