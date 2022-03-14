package pid_test

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/temperaturecontrol/pid"
	"github.com/benjaminbartels/zymurgauge/internal/test/mocks"
	"github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	kP                      float64 = 10
	kI                      float64 = 0
	kD                      float64 = 0
	thermometerReadErrorMsg         = "could not read thermometer"
	actuatorOnErrorMsg              = "could not turn actuator on"
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
		actuatorOn  bool
	}{
		{name: "below", temperature: 10, setPoint: 15, actuatorOn: true},
		{name: "same", temperature: 15, setPoint: 15, actuatorOn: false},
		{name: "above", temperature: 20, setPoint: 15, actuatorOn: false},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l, _ := logtest.NewNullLogger()
			l.SetLevel(logrus.DebugLevel)

			var actuatorOn bool

			thermometerMock := &mocks.Thermometer{}
			thermometerMock.On("GetTemperature").Return(tc.temperature, nil)

			actuatorCh := make(chan struct{}, 1)
			actuatorMock := &mocks.Actuator{}
			actuatorMock.Mock.On("On").Return(nil).Run(
				func(args mock.Arguments) {
					actuatorCh <- struct{}{}
				})
			actuatorMock.Mock.On("Off").Return(nil)

			ctlr := pid.NewPIDTemperatureController(thermometerMock, actuatorMock, kP, kI, kD, l)

			ctx, stop := context.WithCancel(context.Background())

			if tc.actuatorOn {
				// wait for chillerFake.On to be called and we receive the elapsed time on channel
				go func() {
					<-actuatorCh
					actuatorOn = true
					stop()
				}()
			}

			// timeout
			go func() {
				<-time.After(100 * time.Millisecond)
				stop()
			}()

			// this will block until stop() is called
			err := ctlr.Run(ctx, tc.setPoint)
			assert.NoError(t, err)

			assert.Equal(t, tc.actuatorOn, actuatorOn, "expected actuator to be on")
		})
	}
}

func TestRunDutyCycle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		temperature         float64
		expectedElapsedTime time.Duration
	}{
		{name: "0% duty", temperature: 20, expectedElapsedTime: 0},
		{name: "50% duty", temperature: 15, expectedElapsedTime: 500 * time.Millisecond},
		{name: "100% duty", temperature: 10, expectedElapsedTime: 1 * time.Second},
	}

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

			actuatorMock := &mocks.Actuator{}
			actuatorMock.Mock.On("On").Return(nil).Run(on)
			actuatorMock.Mock.On("Off").Return(nil).Run(off)
			therm := pid.NewPIDTemperatureController(thermometerMock, actuatorMock, kP, kI, kD, l)

			var elapsedTime time.Duration
			var mu sync.RWMutex

			ctx, stop := context.WithCancel(context.Background())

			// wait for actuatorMock.Off to be called and we receive the elapsed time on channel
			go func() {
				mu.Lock()
				elapsedTime = <-doneCh
				stop()
				mu.Unlock()
			}()

			// timeout
			go func() {
				<-time.After(1000 * time.Millisecond)
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

	doneCh := make(chan struct{}, 1)
	actuatorMock := &mocks.Actuator{}
	actuatorMock.Mock.On("On").Return(nil)
	actuatorMock.Mock.On("Off").Return(nil).Run(
		func(args mock.Arguments) {
			doneCh <- struct{}{}
		})

	therm := pid.NewPIDTemperatureController(thermometerMock, actuatorMock, kP, kI, kD, l,
		pid.CyclePeriod(100*time.Millisecond))

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

	assert.True(t, logContains(hook.AllEntries(), logrus.DebugLevel,
		"Actuator current temperature is 25.0000°C, set point is 30.0000°C"))
	assert.True(t, logContains(hook.AllEntries(), logrus.DebugLevel,
		"Actuator dutyCycle is 50.00%, dutyTime is 50ms, waitTime is 50ms"))
	assert.True(t, logContains(hook.AllEntries(), logrus.DebugLevel, "Actuator acting for 50ms"))
	assert.True(t, logContains(hook.AllEntries(), logrus.DebugLevel, "Actuator acted for 50ms"))
	assert.True(t, logContains(hook.AllEntries(), logrus.DebugLevel, "Actuator waiting for 50ms"))
	assert.True(t, logContains(hook.AllEntries(), logrus.DebugLevel, "Actuator waited for 50ms"))
}

func TestRunAlreadyRunningError(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()
	ctx, stop := context.WithCancel(context.Background())

	thermometerMock := &mocks.Thermometer{}
	thermometerMock.On("GetTemperature").Return(20.0, nil)

	doneCh := make(chan struct{}, 1)
	actuatorMock := &mocks.Actuator{}
	actuatorMock.Mock.On("Off").Return(nil).Run(func(args mock.Arguments) {
		doneCh <- struct{}{}
	})

	ctrl := pid.NewPIDTemperatureController(thermometerMock, actuatorMock, kP, kI, kD, l)

	go func() {
		// wait until first therm.Run is called
		<-doneCh

		err := ctrl.Run(ctx, 66)
		assert.ErrorIs(t, err, pid.ErrAlreadyRunning)
		stop()
	}()

	// first therm.Run is called
	err := ctrl.Run(ctx, 20)
	assert.NoError(t, err)
}

func TestThermometerIsNilError(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()

	actuatorMock := &mocks.Actuator{}

	ctrl := pid.NewPIDTemperatureController(nil, actuatorMock, kP, kI, kD, l)

	// first therm.Run is called
	err := ctrl.Run(context.Background(), 20)
	assert.ErrorIs(t, err, pid.ErrThermometerIsNil)
}

func TestActuatorIsNilError(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()

	thermometerMock := &mocks.Thermometer{}
	thermometerMock.On("GetTemperature").Return(20.0, nil)

	ctrl := pid.NewPIDTemperatureController(thermometerMock, nil, kP, kI, kD, l)

	// first therm.Run is called
	err := ctrl.Run(context.Background(), 20)
	assert.ErrorIs(t, err, pid.ErrActuatorIsNil)
}

func TestThermometerError(t *testing.T) {
	t.Parallel()

	l, hook := logtest.NewNullLogger()

	doneCh := make(chan struct{}, 1)

	thermometerMock := &mocks.Thermometer{}
	thermometerMock.On("GetTemperature").Return(0.0, errDeadThermometer)

	actuatorMock := &mocks.Actuator{}

	ctrl := pid.NewPIDTemperatureController(thermometerMock, actuatorMock, kP, kI, kD, l)

	go func() {
		for {
			if logContains(hook.AllEntries(), logrus.ErrorLevel, thermometerReadErrorMsg) {
				doneCh <- struct{}{}

				return
			}

			<-time.After(100 * time.Millisecond)
		}
	}()

	go func() {
		_ = ctrl.Run(context.Background(), 15)
	}()

	select {
	case <-doneCh:
	case <-time.After(5 * time.Second):
		assert.Fail(t, "log should contain expected value by now")
	}
}

func TestActuatorOnError(t *testing.T) {
	t.Parallel()

	l, hook := logtest.NewNullLogger()

	thermometerMock := &mocks.Thermometer{}
	thermometerMock.On("GetTemperature").Return(20.0, nil)

	doneCh := make(chan struct{}, 1)

	actuatorMock := &mocks.Actuator{}
	actuatorMock.Mock.On("On").Return(errDeadActuator)

	ctrl := pid.NewPIDTemperatureController(thermometerMock, actuatorMock, kP, kI, kD, l)

	go func() {
		for {
			if logContains(hook.AllEntries(), logrus.ErrorLevel, actuatorOnErrorMsg) {
				doneCh <- struct{}{}

				return
			}

			<-time.After(100 * time.Millisecond)
		}
	}()

	go func() {
		_ = ctrl.Run(context.Background(), 25)
	}()

	select {
	case <-doneCh:
	case <-time.After(5 * time.Second):
		assert.Fail(t, "log should contain expected value by now")
	}
}

func TestActuatorOffErrorAfterDuty(t *testing.T) {
	t.Parallel()

	l, hook := logtest.NewNullLogger()

	thermometerMock := &mocks.Thermometer{}
	thermometerMock.On("GetTemperature").Return(20.0, nil)

	doneCh := make(chan struct{}, 1)

	actuatorMock := &mocks.Actuator{}
	actuatorMock.Mock.On("On").Return(nil)
	actuatorMock.Mock.On("Off").Return(errDeadActuator)

	ctrl := pid.NewPIDTemperatureController(thermometerMock, actuatorMock, kP, kI, kD, l)

	go func() {
		for {
			if logContains(hook.AllEntries(), logrus.ErrorLevel, actuatorOffError) {
				doneCh <- struct{}{}

				return
			}

			<-time.After(100 * time.Millisecond)
		}
	}()

	go func() {
		_ = ctrl.Run(context.Background(), 15)
	}()

	select {
	case <-doneCh:
	case <-time.After(5 * time.Second):
		assert.Fail(t, "log should contain expected value by now")
	}
}

func TestActuatorOffErrorOnQuit(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()

	thermometerMock := &mocks.Thermometer{}
	thermometerMock.On("GetTemperature").Return(20.0, nil)

	doneCh := make(chan struct{}, 1)
	actuatorMock := &mocks.Actuator{}
	actuatorMock.Mock.On("On").Return(nil).Run(
		func(args mock.Arguments) {
			doneCh <- struct{}{}
		})
	actuatorMock.Mock.On("Off").Return(errDeadActuator)

	ctlr := pid.NewPIDTemperatureController(thermometerMock, actuatorMock, kP, kI, kD, l)

	ctx, stop := context.WithCancel(context.Background())

	go func() {
		// wait until first therm.Run is called
		<-doneCh
		stop()
	}()

	err := ctlr.Run(ctx, 25)
	assert.Contains(t, err.Error(), actuatorQuitError)
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
