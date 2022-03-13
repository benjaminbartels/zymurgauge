package hysteresis_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/temperaturecontrol/hysteresis"
	"github.com/benjaminbartels/zymurgauge/internal/test/mocks"
	"github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	hysteresisBand          = 1.0
	thermometerReadErrorMsg = "could not read thermometer"
	chillerOnErrorMsg       = "could not turn chiller actuator on"
	chillerOffErrorMsg      = "could not turn chiller actuator off"
	heaterOnErrorMsg        = "could not turn heater actuator on"
	heaterOffErrorMsg       = "could not turn heater actuator off"
	cooldownLogMsg          = "Cannot turn chiller on for another"
)

var (
	errDeadThermometer = errors.New("thermometer is dead")
	errDeadActuator    = errors.New("actuator is dead")
)

func TestRunFirstRead(t *testing.T) {
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
			chillerMock.Mock.On("On").Return(nil).Run(
				func(args mock.Arguments) {
					chillerCh <- struct{}{}
				})
			chillerMock.Mock.On("Off").Return(nil)

			heaterCh := make(chan struct{}, 1)
			heaterMock := &mocks.Actuator{}
			heaterMock.Mock.On("On").Return(nil).Run(
				func(args mock.Arguments) {
					heaterCh <- struct{}{}
				})
			heaterMock.Mock.On("Off").Return(nil)

			ctlr := hysteresis.NewController(thermometerMock, chillerMock, heaterMock, hysteresisBand, l)

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

func TestRunConsecutiveTemperatureChanges(t *testing.T) {
	t.Parallel()

	setPoint := 10.0

	//nolint: dupl
	tests := []struct {
		name     string
		temp1    float64
		temp2    float64
		chiller1 bool
		heater1  bool
		chiller2 bool
		heater2  bool
	}{
		{name: "betweenToChill", temp1: 10, temp2: 15, chiller1: false, heater1: false, chiller2: true, heater2: false},
		{name: "betweenToHeat", temp1: 10, temp2: 5, chiller1: false, heater1: false, chiller2: false, heater2: true},
		{name: "chillToBetween", temp1: 20, temp2: 10, chiller1: true, heater1: false, chiller2: true, heater2: false},
		{name: "chillToHeat", temp1: 20, temp2: 5, chiller1: true, heater1: false, chiller2: false, heater2: true},
		{name: "heatToBetween", temp1: 5, temp2: 10, chiller1: false, heater1: true, chiller2: false, heater2: true},
		{name: "heatToChill", temp1: 5, temp2: 20, chiller1: false, heater1: true, chiller2: true, heater2: false},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l, _ := logtest.NewNullLogger()

			var (
				chillerOn bool
				heaterOn  bool
			)

			doneCh := make(chan struct{}, 1)

			thermometerMock := &mocks.Thermometer{}
			thermometerMock.On("GetTemperature").Return(tc.temp1, nil).Once()

			thermometerMock.On("GetTemperature").Return(tc.temp2, nil).Once().Run(
				func(args mock.Arguments) {
					assert.Equal(t, tc.chiller1, chillerOn, "unexpected first chiller state")
					assert.Equal(t, tc.heater1, heaterOn, "unexpected first heater state")
				})

			thermometerMock.On("GetTemperature").Return(tc.temp2, nil).Once().Run(
				func(args mock.Arguments) {
					assert.Equal(t, tc.chiller2, chillerOn, "unexpected second chiller state")
					assert.Equal(t, tc.heater2, heaterOn, "unexpected second heater state")

					doneCh <- struct{}{}
				})

			chillerMock := &mocks.Actuator{}
			chillerMock.Mock.On("On").Return(nil).Run(
				func(args mock.Arguments) {
					chillerOn = true
				})
			chillerMock.Mock.On("Off").Return(nil).Run(
				func(args mock.Arguments) {
					chillerOn = false
				})

			heaterMock := &mocks.Actuator{}
			heaterMock.Mock.On("On").Return(nil).Run(
				func(args mock.Arguments) {
					heaterOn = true
				})
			heaterMock.Mock.On("Off").Return(nil).Run(
				func(args mock.Arguments) {
					heaterOn = false
				})

			opts := hysteresis.CyclePeriod(1 * time.Millisecond)

			ctlr := hysteresis.NewController(thermometerMock, chillerMock, heaterMock, hysteresisBand, l, opts)

			ctx, stop := context.WithCancel(context.Background())

			// timeout
			go func() {
				<-doneCh
				stop()
			}()

			// this will block until stop() is called
			err := ctlr.Run(ctx, setPoint)
			assert.NoError(t, err)
		})
	}
}

func TestRunConsecutiveSetPointChanges(t *testing.T) {
	t.Parallel()

	temperature := 20.0

	//nolint: dupl
	tests := []struct {
		name     string
		sp1      float64
		sp2      float64
		chiller1 bool
		heater1  bool
		chiller2 bool
		heater2  bool
	}{
		{name: "betweenToChill", sp1: 20, sp2: 15, chiller1: false, heater1: false, chiller2: true, heater2: false},
		{name: "betweenToHeat", sp1: 20, sp2: 25, chiller1: false, heater1: false, chiller2: false, heater2: true},
		{name: "chillToBetween", sp1: 15, sp2: 20, chiller1: true, heater1: false, chiller2: true, heater2: false},
		{name: "chillToHeat", sp1: 15, sp2: 25, chiller1: true, heater1: false, chiller2: false, heater2: true},
		{name: "heatToBetween", sp1: 25, sp2: 20, chiller1: false, heater1: true, chiller2: false, heater2: true},
		{name: "heatToChill", sp1: 25, sp2: 15, chiller1: false, heater1: true, chiller2: true, heater2: false},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l, _ := logtest.NewNullLogger()

			var (
				chillerOn bool
				heaterOn  bool
			)

			doneCh := make(chan struct{}, 1)

			var ctlr *hysteresis.Controller

			thermometerMock := &mocks.Thermometer{}
			thermometerMock.On("GetTemperature").Return(temperature, nil).Once().Run( // do nothing on 1st run
				func(args mock.Arguments) {
					go ctlr.SetTemperature(tc.sp2) // change set point for for 2nd run
				})

			thermometerMock.On("GetTemperature").Return(temperature, nil).Once().Run( // check 1st run states on 2nd run
				func(args mock.Arguments) {
					assert.Equal(t, tc.chiller1, chillerOn, "unexpected first chiller state")
					assert.Equal(t, tc.heater1, heaterOn, "unexpected first heater state")
				})

			thermometerMock.On("GetTemperature").Return(temperature, nil).Once().Run( // check 2nd run states on 3rd run
				func(args mock.Arguments) {
					assert.Equal(t, tc.chiller2, chillerOn, "unexpected second chiller state")
					assert.Equal(t, tc.heater2, heaterOn, "unexpected second heater state")
					doneCh <- struct{}{}
				})

			chillerMock := &mocks.Actuator{}
			chillerMock.Mock.On("On").Return(nil).Run(
				func(args mock.Arguments) {
					chillerOn = true
				})
			chillerMock.Mock.On("Off").Return(nil).Run(
				func(args mock.Arguments) {
					chillerOn = false
				})

			heaterMock := &mocks.Actuator{}
			heaterMock.Mock.On("On").Return(nil).Run(
				func(args mock.Arguments) {
					heaterOn = true
				})
			heaterMock.Mock.On("Off").Return(nil).Run(
				func(args mock.Arguments) {
					heaterOn = false
				})

			opts := hysteresis.CyclePeriod(1 * time.Millisecond)

			ctlr = hysteresis.NewController(thermometerMock, chillerMock, heaterMock, hysteresisBand, l, opts)

			ctx, stop := context.WithCancel(context.Background())

			// timeout
			go func() {
				<-doneCh
				stop()
			}()

			// this will block until stop() is called
			err := ctlr.Run(ctx, tc.sp1)
			assert.NoError(t, err)
		})
	}
}

func TestChillerCooldown(t *testing.T) {
	t.Parallel()

	l, hook := logtest.NewNullLogger()
	l.SetLevel(logrus.DebugLevel)

	doneCh := make(chan struct{}, 1)

	thermometerMock := &mocks.Thermometer{}
	thermometerMock.On("GetTemperature").Once().Return(20.0, nil)
	thermometerMock.On("GetTemperature").Once().Return(5.0, nil)
	thermometerMock.On("GetTemperature").Once().Return(20.0, nil)

	chillerMock := &mocks.Actuator{}
	chillerMock.Mock.On("On").Return(nil).Once()  // 1st loop
	chillerMock.Mock.On("Off").Return(nil).Once() // 2nd loop
	chillerMock.Mock.On("Off").Return(nil).Once() // stop

	heaterMock := &mocks.Actuator{}
	heaterMock.Mock.On("Off").Return(nil).Once()      // 1st loop
	heaterMock.Mock.On("On").Return(nil).Once()       // 2nd loop
	heaterMock.Mock.On("Off").Return(nil).Once().Run( // 3rd loop
		func(args mock.Arguments) {
			doneCh <- struct{}{}
		})
	heaterMock.Mock.On("Off").Return(nil).Once() // stop

	cycleOpts := hysteresis.CyclePeriod(100 * time.Millisecond)
	cooldownOpts := hysteresis.ChillerCooldown(10 * time.Second)

	ctlr := hysteresis.NewController(thermometerMock, chillerMock, heaterMock, hysteresisBand, l, cycleOpts, cooldownOpts)

	ctx, stop := context.WithCancel(context.Background())

	// timeout
	go func() {
		<-doneCh
		stop()
	}()

	// this will block until stop() is called
	err := ctlr.Run(ctx, 10.0)
	assert.NoError(t, err)

	assert.True(t, logContains(hook.AllEntries(), logrus.DebugLevel, cooldownLogMsg))
}

func TestRunAlreadyRunningError(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()
	ctx, stop := context.WithCancel(context.Background())

	thermometerMock := &mocks.Thermometer{}
	thermometerMock.On("GetTemperature").Return(25.0, nil)

	doneCh := make(chan struct{}, 1)
	chillerMock := &mocks.Actuator{}
	chillerMock.Mock.On("On").Return(nil).Run(func(args mock.Arguments) {
		doneCh <- struct{}{}
	})
	chillerMock.Mock.On("Off").Return(nil)

	heaterMock := &mocks.Actuator{}
	heaterMock.Mock.On("Off").Return(nil)

	ctrl := hysteresis.NewController(thermometerMock, chillerMock, heaterMock, hysteresisBand, l)

	go func() {
		// wait until first therm.Run is called
		<-doneCh

		err := ctrl.Run(ctx, 66)
		assert.ErrorIs(t, err, hysteresis.ErrAlreadyRunning)
		stop()
	}()

	// first therm.Run is called
	err := ctrl.Run(ctx, 20)
	assert.NoError(t, err)
}

func TestThermometerIsNilError(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()

	chillerMock := &mocks.Actuator{}

	heaterMock := &mocks.Actuator{}

	ctrl := hysteresis.NewController(nil, chillerMock, heaterMock, hysteresisBand, l)

	// first therm.Run is called
	err := ctrl.Run(context.Background(), 20)
	assert.ErrorIs(t, err, hysteresis.ErrThermometerIsNil)
}

func TestActuatorIsNilError(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()

	thermometerMock := &mocks.Thermometer{}
	thermometerMock.On("GetTemperature").Return(20.0, nil)

	ctrl := hysteresis.NewController(thermometerMock, nil, nil, hysteresisBand, l)

	// first therm.Run is called
	err := ctrl.Run(context.Background(), 20)
	assert.ErrorIs(t, err, hysteresis.ErrActuatorIsNil)
	assert.Contains(t, err.Error(), hysteresis.ErrActuatorIsNil.Error()) // For coverage
}

func TestThermometerError(t *testing.T) {
	t.Parallel()

	l, hook := logtest.NewNullLogger()

	doneCh := make(chan struct{}, 1)

	thermometerMock := &mocks.Thermometer{}
	thermometerMock.On("GetTemperature").Return(0.0, errDeadThermometer)

	chillerMock := &mocks.Actuator{}
	heaterMock := &mocks.Actuator{}

	ctrl := hysteresis.NewController(thermometerMock, chillerMock, heaterMock, hysteresisBand, l)

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

func TestActuatorErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		temperature     float64
		chillerErrorMsg string
		heaterErrorMsg  string
	}{
		{name: "chillingErrors", temperature: 15, chillerErrorMsg: chillerOnErrorMsg, heaterErrorMsg: heaterOffErrorMsg},
		{name: "heatingErrors", temperature: 25, chillerErrorMsg: chillerOffErrorMsg, heaterErrorMsg: heaterOnErrorMsg},
	}

	for _, tc := range tests {
		tc := tc

		l, hook := logtest.NewNullLogger()

		thermometerMock := &mocks.Thermometer{}
		thermometerMock.On("GetTemperature").Return(20.0, nil)

		doneCh := make(chan struct{}, 1)

		chillerMock := &mocks.Actuator{}
		chillerMock.Mock.On("On").Return(errDeadActuator)
		chillerMock.Mock.On("Off").Return(errDeadActuator)

		heaterMock := &mocks.Actuator{}
		heaterMock.Mock.On("On").Return(errDeadActuator)
		heaterMock.Mock.On("Off").Return(errDeadActuator)

		ctrl := hysteresis.NewController(thermometerMock, chillerMock, heaterMock, hysteresisBand, l)

		go func() {
			for {
				if logContains(hook.AllEntries(), logrus.ErrorLevel, tc.chillerErrorMsg) &&
					logContains(hook.AllEntries(), logrus.ErrorLevel, tc.heaterErrorMsg) {
					doneCh <- struct{}{}

					return
				}

				<-time.After(100 * time.Millisecond)
			}
		}()

		go func() {
			_ = ctrl.Run(context.Background(), tc.temperature)
		}()

		select {
		case <-doneCh:
		case <-time.After(5 * time.Second):
			assert.Fail(t, "log should contain expected value by now")
		}
	}
}

func TestActuatorOffErrorOnQuit(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()

	doneCh := make(chan struct{}, 1)

	thermometerMock := &mocks.Thermometer{}
	thermometerMock.On("GetTemperature").Return(20.0, nil).Run(
		func(args mock.Arguments) {
			doneCh <- struct{}{}
		})

	chillerMock := &mocks.Actuator{}
	chillerMock.Mock.On("On").Return(nil)
	chillerMock.Mock.On("Off").Return(errDeadActuator)

	heaterMock := &mocks.Actuator{}
	heaterMock.Mock.On("Off").Return(nil).Once()
	heaterMock.Mock.On("Off").Return(errDeadActuator)

	ctlr := hysteresis.NewController(thermometerMock, chillerMock, heaterMock, hysteresisBand, l)

	ctx, stop := context.WithCancel(context.Background())

	go func() {
		// wait until first therm.Run is called
		<-doneCh
		stop()
	}()

	err := ctlr.Run(ctx, 15)
	assert.Contains(t, err.Error(), chillerOffErrorMsg)
	assert.Contains(t, err.Error(), heaterOffErrorMsg)
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
