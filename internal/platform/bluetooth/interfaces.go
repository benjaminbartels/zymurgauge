package bluetooth

import (
	"context"

	"github.com/go-ble/ble"
)

type Scanner interface {
	WithSigHandler(ctx context.Context, cancel func()) context.Context
	Scan(ctx context.Context, h func(a Advertisement), f func(a Advertisement) bool) error
}

type Advertisement interface {
	ble.Advertisement
}

type Device interface {
	ble.Device
}
