//go:build !linux || !arm
// +build !linux !arm

package bluetooth

import (
	"context"

	"github.com/go-ble/ble/linux"
)

// The program is only meant to run on linux on arm. This file only exists to prevent compilation issues on non
// linux/arm systems.

var _ Scanner = (*BLEScanner)(nil)

type BLEScanner struct{}

func NewBLEScanner() *BLEScanner {
	return &BLEScanner{}
}

func (b *BLEScanner) NewDevice() (*linux.Device, error) {
	device, _ := linux.NewDevice()

	return device, nil
}

func (b *BLEScanner) SetDefaultDevice(device Device) {}

func (b *BLEScanner) WithSigHandler(ctx context.Context, cancel func()) context.Context {
	return ctx
}

func (b *BLEScanner) Scan(ctx context.Context, h func(a Advertisement), f func(a Advertisement) bool) error {
	return nil
}
