package tilt_test

import (
	"context"
	"testing"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/device/tilt"
	mocks "github.com/benjaminbartels/zymurgauge/internal/test/mocks/bluetooth"
	"github.com/go-ble/ble/linux"
	"github.com/pkg/errors"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	errDevice = errors.New("device error")
	errScan   = errors.New("unrecoverable scan error")
)

func TestRunAlreadyRunningError(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()
	ctx, stop := context.WithCancel(context.Background())
	device := &linux.Device{}

	scannerMock := &mocks.Scanner{}
	scannerMock.On("NewDevice").Return(device, nil)
	scannerMock.On("SetDefaultDevice", device).Return()
	scannerMock.On("WithSigHandler", mock.Anything, mock.Anything).Return(ctx)

	timeout := tilt.SetTimeout(1 * time.Second)
	interval := tilt.SetInterval(1 * time.Second)
	monitor := tilt.NewMonitor(scannerMock, l, timeout, interval)

	scannerMock.On("Scan", mock.Anything, mock.Anything, mock.Anything).Return(context.DeadlineExceeded).Run(
		func(args mock.Arguments) {
			err := monitor.Run(ctx)
			assert.ErrorIs(t, err, tilt.ErrAlreadyRunning)

			stop()
		})

	// first monitor.Run is called
	err := monitor.Run(ctx)
	assert.NoError(t, err)
}

func TestRunNewDeviceError(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()

	scannerMock := &mocks.Scanner{}
	scannerMock.On("NewDevice").Return(nil, errDevice)

	monitor := tilt.NewMonitor(scannerMock, l)

	err := monitor.Run(context.Background())
	assert.Contains(t, err.Error(), "could not create new device: device error")
}

func TestScanDeadlineExceeded(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()
	device := &linux.Device{}
	ctx, stop := context.WithCancel(context.Background())
	ctr := 0

	scannerMock := &mocks.Scanner{}
	scannerMock.On("NewDevice").Return(device, nil)
	scannerMock.On("SetDefaultDevice", device).Return()
	scannerMock.On("WithSigHandler", mock.Anything, mock.Anything).Return(ctx)
	scannerMock.On("Scan", mock.Anything, mock.Anything, mock.Anything).Return(context.DeadlineExceeded).Run(
		func(args mock.Arguments) {
			ctr++
			if ctr == 2 {
				stop() // stop after second call to Scan
			}
		})

	timeout := tilt.SetTimeout(1 * time.Millisecond)
	interval := tilt.SetInterval(1 * time.Millisecond)
	monitor := tilt.NewMonitor(scannerMock, l, timeout, interval)

	err := monitor.Run(ctx)
	assert.NoError(t, err)
}

func TestScanCancelled(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()
	device := &linux.Device{}
	ctx := context.Background()

	scannerMock := &mocks.Scanner{}
	scannerMock.On("NewDevice").Return(device, nil)
	scannerMock.On("SetDefaultDevice", device).Return()
	scannerMock.On("WithSigHandler", mock.Anything, mock.Anything).Return(ctx)
	scannerMock.On("Scan", mock.Anything, mock.Anything, mock.Anything).Return(context.Canceled)

	timeout := tilt.SetTimeout(1 * time.Millisecond)
	interval := tilt.SetInterval(1 * time.Millisecond)
	monitor := tilt.NewMonitor(scannerMock, l, timeout, interval)

	err := monitor.Run(ctx)

	assert.NoError(t, err)
}

func TestScanOtherError(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()
	device := &linux.Device{}
	ctx := context.Background()

	scannerMock := &mocks.Scanner{}
	scannerMock.On("NewDevice").Return(device, nil)
	scannerMock.On("SetDefaultDevice", device).Return()
	scannerMock.On("WithSigHandler", mock.Anything, mock.Anything).Return(ctx)
	scannerMock.On("Scan", mock.Anything, mock.Anything, mock.Anything).Return(errScan)

	timeout := tilt.SetTimeout(1 * time.Millisecond)
	interval := tilt.SetInterval(1 * time.Millisecond)
	monitor := tilt.NewMonitor(scannerMock, l, timeout, interval)

	err := monitor.Run(ctx)
	assert.Contains(t, err.Error(), "could not scan: unrecoverable scan error")
}

func TestGetTiltNotFound(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()
	device := &linux.Device{}
	ctx, stop := context.WithCancel(context.Background())

	scannerMock := &mocks.Scanner{}
	scannerMock.On("NewDevice").Return(device, nil)
	scannerMock.On("SetDefaultDevice", device).Return()
	scannerMock.On("WithSigHandler", mock.Anything, mock.Anything).Return(ctx)

	timeout := tilt.SetTimeout(1 * time.Millisecond)
	interval := tilt.SetInterval(1 * time.Millisecond)
	monitor := tilt.NewMonitor(scannerMock, l, timeout, interval)

	scannerMock.On("Scan", mock.Anything, mock.Anything, mock.Anything).Return(context.DeadlineExceeded).Run(
		func(args mock.Arguments) {
			ti, err := monitor.GetTilt("badColor")
			assert.Nil(t, ti)
			assert.ErrorIs(t, err, tilt.ErrNotFound)

			stop()
		})

	err := monitor.Run(ctx)
	assert.NoError(t, err)
}
