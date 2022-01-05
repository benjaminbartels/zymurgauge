//go:build !linux || !arm
// +build !linux !arm

package bluetooth

import (
	"context"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
)

// The program is only meant to run on linux on arm. This file only exists to prevent compilation issues on non
// linux/arm systems.

var _ Scanner = (*StubBLEScanner)(nil)

type Advertisement interface {
	ble.Advertisement
}

type Device interface {
	ble.Device
}

type StubBLEScanner struct{}

func NewBLEScanner() *StubBLEScanner {
	return &StubBLEScanner{}
}

func (b *StubBLEScanner) NewDevice() (*linux.Device, error) {
	device, _ := linux.NewDevice()

	return device, nil
}

func (b *StubBLEScanner) SetDefaultDevice(device Device) {}

func (b *StubBLEScanner) WithSigHandler(ctx context.Context, cancel func()) context.Context {
	return ctx
}

func (b *StubBLEScanner) Scan(ctx context.Context, h func(a Advertisement), f func(a Advertisement) bool) error {
	return nil
}
