package tilt

import (
	"context"
	"strings"
	"testing"
	"time"

	mocks "github.com/benjaminbartels/zymurgauge/internal/test/mocks/bluetooth"
	"github.com/go-ble/ble/linux"
	"github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	temperature     = 20.0
	specificGravity = 0.999
	orange          = "orange"
)

func getOrangeTiltManufacurerData() []byte {
	return []byte{
		76, 0, 2, 21, 164, 149, 187, 80, 197, 177, 75, 68, 181, 18, 19, 112, 240, 45, 116, 222, 0, 68, 3, 231, 197,
	}
}

func getInvalidTiltManufacurerData() []byte {
	return []byte{
		76, 0, 2, 21, 164, 149, 187,
	}
}

func TestGetTilt(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()
	device := &linux.Device{}
	ctx, stop := context.WithCancel(context.Background())

	scannerMock := &mocks.Scanner{}
	scannerMock.On("NewDevice").Return(device, nil)
	scannerMock.On("SetDefaultDevice", device).Return()
	scannerMock.On("WithSigHandler", mock.Anything, mock.Anything).Return(ctx)

	monitor := NewMonitor(scannerMock, l)

	scannerMock.On("Scan", mock.Anything, mock.Anything, mock.Anything).Return(context.DeadlineExceeded).Run(
		func(args mock.Arguments) {
			adv := &mocks.Advertisement{}
			adv.Mock.On("ManufacturerData").Return(getOrangeTiltManufacurerData)

			isTilt := monitor.filter(adv)
			assert.True(t, isTilt, "expected manufacurerData to represent a tilt")

			monitor.handler(adv)

			ti, err := monitor.GetTilt(orange)
			assert.NoError(t, err)

			temp, err := ti.GetTemperature()
			assert.NoError(t, err)
			assert.Equal(t, temp, temperature)

			sg, err := ti.GetGravity()
			assert.NoError(t, err)
			assert.Equal(t, sg, specificGravity)

			stop()
		})

	err := monitor.Run(ctx)
	assert.NoError(t, err)
}

func TestGetTiltIBeaconIsNilError(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()
	device := &linux.Device{}
	ctx, stop := context.WithCancel(context.Background())

	scannerMock := &mocks.Scanner{}
	scannerMock.On("NewDevice").Return(device, nil)
	scannerMock.On("SetDefaultDevice", device).Return()
	scannerMock.On("WithSigHandler", mock.Anything, mock.Anything).Return(ctx)

	timeout := SetTimeout(1 * time.Millisecond)
	interval := SetInterval(1 * time.Millisecond)
	monitor := NewMonitor(scannerMock, l, timeout, interval)

	ctr := 0

	var orangeTilt *Tilt

	scannerMock.On("Scan", mock.Anything, mock.Anything, mock.Anything).Return(context.DeadlineExceeded).Run(
		func(args mock.Arguments) {
			ctr++

			switch ctr {
			case 1:
				// on first call to Scan add orange tilt to monitor.tilts and monitor.availableColors
				adv := &mocks.Advertisement{}
				adv.Mock.On("ManufacturerData").Return(getOrangeTiltManufacurerData)

				isTilt := monitor.filter(adv)
				assert.True(t, isTilt, "expected manufacurerData to represent a tilt")

				monitor.handler(adv)
			case 2:
				// on second first call to Scan check to see if orange tilt exists
				var err error
				orangeTilt, err = monitor.GetTilt(orange)
				assert.NoError(t, err)

				temp, err := orangeTilt.GetTemperature()
				assert.NoError(t, err)
				assert.Equal(t, temp, temperature)

				sg, err := orangeTilt.GetGravity()
				assert.NoError(t, err)
				assert.Equal(t, sg, specificGravity)

				// after second call to Scan orange tilt should be removed from monitor.tilts and monitor.availableColors
			case 3:
				// check the tilt pointer points to a tilt with a nil beacon
				temp, err := orangeTilt.GetTemperature()
				assert.Error(t, err, ErrIBeaconIsNil)
				assert.Equal(t, temp, 0.0)

				sg, err := orangeTilt.GetGravity()
				assert.Error(t, err, ErrIBeaconIsNil)
				assert.Equal(t, sg, 0.0)

				// check the GetTilt return not found now
				orangeTilt, err = monitor.GetTilt(orange)
				assert.NoError(t, err)

				// after third call to Scan stop monitor
				stop()
			}
		})

	err := monitor.Run(ctx)
	assert.NoError(t, err)
}

func TestHandlerInvalidManufacturerDataLengthError(t *testing.T) {
	t.Parallel()

	l, hook := logtest.NewNullLogger()
	device := &linux.Device{}
	ctx, stop := context.WithCancel(context.Background())

	scannerMock := &mocks.Scanner{}
	scannerMock.On("NewDevice").Return(device, nil)
	scannerMock.On("SetDefaultDevice", device).Return()
	scannerMock.On("WithSigHandler", mock.Anything, mock.Anything).Return(ctx)

	monitor := NewMonitor(scannerMock, l)

	scannerMock.On("Scan", mock.Anything, mock.Anything, mock.Anything).Return(context.DeadlineExceeded).Run(
		func(args mock.Arguments) {
			adv := &mocks.Advertisement{}
			adv.Mock.On("ManufacturerData").Return(getInvalidTiltManufacurerData)

			monitor.handler(adv)

			stop()
		})

	err := monitor.Run(ctx)
	assert.NoError(t, err)

	assert.True(t, logContains(hook.AllEntries(), logrus.ErrorLevel, "could not create new IBeacon"))
}

func TestOptions(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()
	expected := time.Nanosecond

	monitor := NewMonitor(nil, l,
		SetTimeout(expected),
		SetInterval(expected))

	if monitor.timeout != expected {
		assert.Equal(t, expected, monitor.timeout)
	}

	if monitor.interval != expected {
		assert.Equal(t, expected, monitor.timeout)
	}
}

func logContains(logs []*logrus.Entry, level logrus.Level, substr string) bool {
	found := false

	for _, v := range logs {
		if strings.Contains(v.Message, substr) && v.Level == level {
			found = true
		}
	}

	return found
}
