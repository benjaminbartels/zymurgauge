package tilt_test

import (
	"context"
	"sync"
	"testing"

	"github.com/benjaminbartels/zymurgauge/internal/device/tilt"
	"github.com/benjaminbartels/zymurgauge/internal/test/mocks"
	"github.com/go-ble/ble/linux"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRunAlreadyRunningError(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()
	ctx := context.Background()
	device := &linux.Device{}
	doneCh := make(chan struct{}, 1)

	scannerMock := &mocks.Scanner{}
	scannerMock.On("NewDevice").Return(device, nil)
	scannerMock.On("SetDefaultDevice", device).Return()
	scannerMock.On("WithSigHandler", mock.Anything, mock.Anything).Return(ctx)
	scannerMock.On("Scan", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		doneCh <- struct{}{}
	})

	monitor := tilt.NewMonitor(scannerMock, l)

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		// wait until first monitor.Run is called
		<-doneCh

		err := monitor.Run(ctx)
		assert.ErrorIs(t, err, tilt.ErrAlreadyRunning)
		wg.Done()
	}()

	go func() {
		// first monitor.Run is called
		wg.Done()

		_ = monitor.Run(ctx)
	}()

	wg.Wait()
}
