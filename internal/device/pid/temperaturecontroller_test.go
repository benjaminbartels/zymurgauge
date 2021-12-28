package pid_test

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/device/pid"
	"github.com/benjaminbartels/zymurgauge/internal/test/mocks"
	"github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	chillerKp               float64 = -10
	chillerKi               float64 = 0
	chillerKd               float64 = 0
	heaterKp                float64 = 10
	heaterKi                float64 = 0
	heaterKd                float64 = 0
	thermometerReadErrorMsg         = "could not read thermometer"
	chillerOnErrorMsg               = "could not turn chiller actuator on"
	actuatorOffError                = "could not turn actuator off"
	actuatorQuitError               = "could not turn actuator off while quiting"
)

var (
	errDeadThermometer = errors.New("thermometer is dead")
	errDeadActuator    = errors.New("actuator is dead")
)

func TestRunActuatorsOn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		temperature float64
		setPoint    float64
		chillerOn   bool
		heaterOn    bool
	}{
		{name: "below", temperature: 10, setPoint: 15, chillerOn: false, heaterOn: true},
		{name: "same", temperature: 15, setPoint: 15, chillerOn: false, heaterOn: false},
		{name: "above", temperature: 20, setPoint: 15, chillerOn: true, heaterOn: false},
	}

	// nolint: paralleltest // False positives with ANY "Run" not in a loop
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l, _ := logtest.NewNullLogger()

			var (
				chillerOn bool
				heaterOn  bool
			)

			thermometerMock := &mocks.Thermometer{}
			thermometerMock.On("GetTemperature").Return(tc.temperature, nil)

			chillerCh := make(chan struct{}, 1)
			chillerMock := &mocks.Actuator{}
			chillerMock.Mock.On("On").Return(nil).Run(func(args mock.Arguments) { chillerCh <- struct{}{} })
			chillerMock.Mock.On("Off").Return(nil)

			heaterCh := make(chan struct{}, 1)
			heaterMock := &mocks.Actuator{}
			heaterMock.Mock.On("On").Return(nil).Run(func(args mock.Arguments) { heaterCh <- struct{}{} })
			heaterMock.Mock.On("Off").Return(nil)

			ctlr := pid.NewPIDTemperatureController(thermometerMock, chillerMock, heaterMock, chillerKp, chillerKi, chillerKd,
				heaterKp, heaterKi, heaterKd, l)

			ctx, stop := context.WithCancel(context.Background())

			if tc.chillerOn {
				// wait for chillerFake.On to be called and we receive the elapsed time on channel
				go func() {
					<-chillerCh
					chillerOn = true
					stop()
				}()
			}

			if tc.heaterOn {
				// wait for heaterFake.On to be called and we receive the elapsed time on channel
				go func() {
					<-heaterCh
					heaterOn = true
					stop()
				}()
			}

			// timeout
			go func() {
				<-time.After(200 * time.Millisecond)
				stop()
			}()

			// this will block until stop() is called
			err := ctlr.Run(ctx, tc.setPoint)
			assert.NoError(t, err)

			assert.Equal(t, tc.chillerOn, chillerOn, "expected chiller to be on")
			assert.Equal(t, tc.heaterOn, heaterOn, "expected heater to be on")
		})
	}
}

func TestRunDutyCycle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		temperature         float64
		expectedElapsedTime time.Duration
		waitTime            time.Duration
	}{
		{name: "0% duty", temperature: 20, expectedElapsedTime: 0 * time.Millisecond},
		{name: "minimum duty", temperature: 20.5, expectedElapsedTime: 10 * time.Millisecond},
		{name: "50% duty", temperature: 25, expectedElapsedTime: 50 * time.Millisecond},
		{name: "100% duty", temperature: 30, expectedElapsedTime: 100 * time.Millisecond},
	}

	// nolint: paralleltest // False positives with ANY "Run" not in a loop
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l, _ := logtest.NewNullLogger()

			thermometerMock := &mocks.Thermometer{}
			thermometerMock.On("GetTemperature").Return(tc.temperature, nil)

			var startTime time.Time
			doneCh := make(chan time.Duration, 1)
			on := func(args mock.Arguments) {
				if startTime.IsZero() {
					startTime = time.Now()
				}
			}

			off := func(args mock.Arguments) {
				if startTime.IsZero() {
					doneCh <- 0
				} else {
					doneCh <- time.Since(startTime)
				}
			}

			chillerMock := &mocks.Actuator{}
			chillerMock.Mock.On("On").Return(nil).Run(on)
			chillerMock.Mock.On("Off").Return(nil).Run(off)

			heaterMock := &mocks.Actuator{}
			heaterMock.Mock.On("Off").Return(nil)

			therm := pid.NewPIDTemperatureController(thermometerMock, chillerMock, heaterMock,
				chillerKp, chillerKi, chillerKd,
				heaterKp, heaterKi, heaterKd,
				l, pid.SetChillingCyclePeriod(100*time.Millisecond),
				pid.SetChillingMinimum(10*time.Millisecond))

			var elapsedTime time.Duration
			var mu sync.RWMutex

			ctx, stop := context.WithCancel(context.Background())

			// wait for chillerFake.Off to be called and we receive the elapsed time on channel
			go func() {
				mu.Lock()
				elapsedTime = <-doneCh
				stop()
				mu.Unlock()
			}()

			// timeout
			go func() {
				<-time.After(200 * time.Millisecond)
				stop()
			}()

			// this will block until stop() is called
			err := therm.Run(ctx, 20)
			assert.NoError(t, err)
			mu.RLock()
			if tc.name == "100% duty" {
				assert.GreaterOrEqual(t, elapsedTime, tc.expectedElapsedTime)
			} else {
				assert.Equal(t, tc.expectedElapsedTime, elapsedTime.Round(10*time.Millisecond))
			}
			mu.RUnlock()
		})
	}
}

func TestLogging(t *testing.T) {
	t.Parallel()

	l, hook := logtest.NewNullLogger()
	l.Level = logrus.DebugLevel

	thermometerMock := &mocks.Thermometer{}
	thermometerMock.On("GetTemperature").Return(25.0, nil)

	chillerMock := &mocks.Actuator{}
	chillerMock.Mock.On("Off").Return(nil)

	doneCh := make(chan struct{}, 1)
	heaterMock := &mocks.Actuator{}
	heaterMock.Mock.On("On").Return(nil)
	heaterMock.Mock.On("Off").Return(nil).Run(func(args mock.Arguments) { doneCh <- struct{}{} })

	therm := pid.NewPIDTemperatureController(thermometerMock, chillerMock, heaterMock, chillerKp, chillerKi, chillerKd,
		heaterKp, heaterKi, heaterKd, l, pid.SetHeatingCyclePeriod(100*time.Millisecond),
		pid.SetHeatingMinimum(50*time.Millisecond))

	ctx, stop := context.WithCancel(context.Background())

	go func() {
		// wait for 2 cycles
		<-doneCh
		<-doneCh
		stop()
	}()

	err := therm.Run(ctx, 30)
	assert.NoError(t, err)

	<-time.After(1 * time.Second)

	assert.True(t, logContains(hook.AllEntries(), logrus.DebugLevel, "Actuator heater current temperature is 25.0000°C"))
	assert.True(t, logContains(hook.AllEntries(), logrus.DebugLevel, "Actuator heater dutyCycle is 50.00%"))
	assert.True(t, logContains(hook.AllEntries(), logrus.DebugLevel, "Actuator heater waitTime is 50ms"))
	assert.True(t, logContains(hook.AllEntries(), logrus.DebugLevel, "Actuator heater acting for 50ms"))
	assert.True(t, logContains(hook.AllEntries(), logrus.DebugLevel, "Actuator heater acted for 50ms"))
	assert.True(t, logContains(hook.AllEntries(), logrus.DebugLevel, "Actuator heater waiting for 50ms"))
	assert.True(t, logContains(hook.AllEntries(), logrus.DebugLevel, "Actuator heater waited for 50ms"))
	assert.True(t, logContains(hook.AllEntries(), logrus.DebugLevel, "Actuator chiller current temperature is 25.0000°C"))
	assert.True(t, logContains(hook.AllEntries(), logrus.DebugLevel, "Actuator chiller dutyCycle is 0.00%"))
	assert.True(t, logContains(hook.AllEntries(), logrus.DebugLevel, "Actuator chiller waitTime is 30m0s"))
	assert.True(t, logContains(hook.AllEntries(), logrus.DebugLevel, "Actuator chiller waiting for 30m0s"))
}

func TestRunAlreadyRunningError(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()

	thermometerMock := &mocks.Thermometer{}
	thermometerMock.On("GetTemperature").Return(20.0, nil)

	doneCh := make(chan struct{}, 1)
	chillerMock := &mocks.Actuator{}
	chillerMock.Mock.On("Off").Return(nil).Run(func(args mock.Arguments) { doneCh <- struct{}{} })

	heaterMock := &mocks.Actuator{}
	heaterMock.Mock.On("Off").Return(nil)

	therm := pid.NewPIDTemperatureController(thermometerMock, chillerMock, heaterMock, chillerKp, chillerKi, chillerKd,
		heaterKp, heaterKi, heaterKd, l)

	ctx := context.Background()

	go func() {
		// first therm.Run is called
		err := therm.Run(ctx, 20)
		assert.NoError(t, err)
	}()

	<-doneCh

	err := therm.Run(ctx, 66)
	assert.ErrorIs(t, err, pid.ErrAlreadyRunning)
}

func TestThermometerError(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()

	thermometerMock := &mocks.Thermometer{}
	thermometerMock.On("GetTemperature").Return(0.0, errDeadThermometer)

	chillerMock := &mocks.Actuator{}
	heaterMock := &mocks.Actuator{}

	therm := pid.NewPIDTemperatureController(thermometerMock, chillerMock, heaterMock, chillerKp, chillerKi, chillerKd,
		heaterKp, heaterKi, heaterKd, l)

	ctx := context.Background()
	err := therm.Run(ctx, 15)
	assert.Contains(t, err.Error(), thermometerReadErrorMsg)
}

func TestActuatorOnError(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()

	thermometerMock := &mocks.Thermometer{}
	thermometerMock.On("GetTemperature").Return(20.0, nil)

	chillerMock := &mocks.Actuator{}
	chillerMock.Mock.On("On").Return(errDeadActuator)

	heaterMock := &mocks.Actuator{}
	heaterMock.Mock.On("Off").Return(nil)

	ctlr := pid.NewPIDTemperatureController(thermometerMock, chillerMock, heaterMock, chillerKp, chillerKi, chillerKd,
		heaterKp, heaterKi, heaterKd, l)

	ctx := context.Background()
	err := ctlr.Run(ctx, 15)
	assert.Contains(t, err.Error(), chillerOnErrorMsg)
}

func TestActuatorOffErrorAfterDuty(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()

	thermometerMock := &mocks.Thermometer{}
	thermometerMock.On("GetTemperature").Return(20.0, nil)

	chillerMock := &mocks.Actuator{}
	chillerMock.Mock.On("On").Return(nil)
	chillerMock.Mock.On("Off").Return(errDeadActuator)

	heaterMock := &mocks.Actuator{}
	heaterMock.Mock.On("Off").Return(nil)

	ctlr := pid.NewPIDTemperatureController(thermometerMock, chillerMock, heaterMock, chillerKp, chillerKi, chillerKd,
		heaterKp, heaterKi, heaterKd, l, pid.SetChillingCyclePeriod(100*time.Millisecond),
		pid.SetChillingMinimum(10*time.Millisecond))

	ctx := context.Background()
	err := ctlr.Run(ctx, 15)
	assert.Contains(t, err.Error(), actuatorOffError)
}

func TestActuatorOffErrorOnQuit(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()

	thermometerMock := &mocks.Thermometer{}
	thermometerMock.On("GetTemperature").Return(20.0, nil)

	doneCh := make(chan struct{}, 1)
	chillerMock := &mocks.Actuator{}
	chillerMock.Mock.On("On").Return(nil).Run(func(args mock.Arguments) { doneCh <- struct{}{} })
	chillerMock.Mock.On("Off").Return(errDeadActuator)

	heaterMock := &mocks.Actuator{}
	heaterMock.Mock.On("Off").Return(nil)

	ctlr := pid.NewPIDTemperatureController(thermometerMock, chillerMock, heaterMock, chillerKp, chillerKi, chillerKd,
		heaterKp, heaterKi, heaterKd, l)

	ctx, stop := context.WithCancel(context.Background())

	go func() {
		// wait until first therm.Run is called
		<-doneCh
		stop()
	}()

	err := ctlr.Run(ctx, 15)
	assert.Contains(t, err.Error(), actuatorQuitError)
}

//nolint: unparam
func logContains(logs []*logrus.Entry, level logrus.Level, substr string) bool {
	found := false

	for _, v := range logs {
		if strings.Contains(v.Message, substr) && v.Level == level {
			found = true
		}
	}

	return found
}
