package bluetooth

import (
	"context"

	"github.com/go-ble/ble/linux"
)

type Scanner interface {
	NewDevice() (*linux.Device, error)
	SetDefaultDevice(device Device)
	WithSigHandler(ctx context.Context, cancel func()) context.Context
	Scan(ctx context.Context, h func(a Advertisement), f func(a Advertisement) bool) error
}
