package thermostat_test

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/test/mocks"
	"github.com/benjaminbartels/zymurgauge/internal/thermostat"
	"github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
)

const (
	chillerKp float64 = -10
	chillerKi float64 = 0
	chillerKd float64 = 0
	heaterKp  float64 = 10
	heaterKi  float64 = 0
	heaterKd  float64 = 0
)

var (
	errDeadThermometer = errors.New("thermometer is dead")
	errDeadActuator    = errors.New("actuator is dead")
)

func TestOnActuatorsOn(t *testing.T) {
	tests := map[string]struct {
		currentTemperature float64
		setPoint           float64
		chillerOn          bool
		heaterOn           bool
		waitCount          int
	}{
		"below": {currentTemperature: 10, setPoint: 15, chillerOn: false, heaterOn: true, waitCount: 1},
		"same":  {currentTemperature: 15, setPoint: 15, chillerOn: false, heaterOn: false, waitCount: 0},
		"above": {currentTemperature: 20, setPoint: 15, chillerOn: true, heaterOn: false, waitCount: 1},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			l, _ := logtest.NewNullLogger()

			var (
				wg        sync.WaitGroup
				chillerOn bool
				heaterOn  bool
			)

			thermometer := &mocks.Thermometer{
				ReadFn: func() (float64, error) { return tc.currentTemperature, nil },
			}

			chiller := &mocks.Actuator{
				OnFn: func() error {
					chillerOn = true
					wg.Done()

					return nil
				},
				OffFn: func() error { return nil },
			}

			heater := &mocks.Actuator{
				OnFn: func() error {
					heaterOn = true
					wg.Done()

					return nil
				},
				OffFn: func() error { return nil },
			}

			therm := thermostat.NewThermostat(thermometer, chiller, heater, chillerKp, chillerKi, chillerKd, heaterKp,
				heaterKi, heaterKd, l)

			wg.Add(tc.waitCount)

			go func() {
				if err := therm.On(tc.setPoint); err != nil {
					t.Errorf("Unexpected error. Got: %+v", err)
				}
			}()

			wg.Wait()

			if tc.chillerOn != chillerOn {
				t.Errorf("Expected chillerOn to be %t", tc.chillerOn)
			}

			if tc.heaterOn && !heaterOn {
				t.Errorf("Expected heaterOn to be %t", tc.heaterOn)
			}
		})
	}
}

func TestOnDutyCycle(t *testing.T) {
	tests := map[string]struct {
		currentTemperature float64
		dutyTime           time.Duration
		waitTime           time.Duration
	}{
		"0% duty":      {currentTemperature: 20, dutyTime: 0 * time.Millisecond, waitTime: 100 * time.Millisecond},
		"minimum duty": {currentTemperature: 21, dutyTime: 10 * time.Millisecond, waitTime: 90 * time.Millisecond},
		"50% duty":     {currentTemperature: 25, dutyTime: 50 * time.Millisecond, waitTime: 50 * time.Millisecond},
		"100% duty":    {currentTemperature: 30, dutyTime: 100 * time.Millisecond, waitTime: 0 * time.Millisecond},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			l, hook := logtest.NewNullLogger()
			l.Level = logrus.DebugLevel

			var (
				wg              sync.WaitGroup
				readCalledCount int
			)

			thermometer := &mocks.Thermometer{
				ReadFn: func() (float64, error) {
					readCalledCount++
					if readCalledCount == 3 {
						wg.Done()
					}

					return tc.currentTemperature, nil
				},
			}

			chiller := &mocks.Actuator{
				OnFn:  func() error { return nil },
				OffFn: func() error { return nil },
			}

			heater := &mocks.Actuator{
				OnFn:  func() error { return nil },
				OffFn: func() error { return nil },
			}

			therm := thermostat.NewThermostat(thermometer, chiller, heater, chillerKp, chillerKi, chillerKd, heaterKp,
				heaterKi, heaterKd, l, thermostat.SetChillingCyclePeriod(100*time.Millisecond),
				thermostat.SetChillingMinimum(10*time.Millisecond))

			wg.Add(1)

			go func() {
				if err := therm.On(20); err != nil {
					t.Errorf("Unexpected error. Got: %+v", err)
				}
			}()

			wg.Wait()

			if tc.dutyTime > 0 {
				if !logContains(hook.AllEntries(), logrus.DebugLevel, fmt.Sprintf("Actuator chiller acted for %s", tc.dutyTime)) {
					t.Errorf("Expected '%s' to be logged", fmt.Sprintf("Actuator chiller acted for %s", tc.dutyTime))
				}
			}

			if !logContains(hook.AllEntries(), logrus.DebugLevel, fmt.Sprintf("Actuator chiller waited for %s", tc.waitTime)) {
				t.Errorf("Expected '%s' to be logged", fmt.Sprintf("Actuator chiller waited for %s", tc.waitTime))
			}
		})
	}
}

func TestOff(t *testing.T) {
	l, _ := logtest.NewNullLogger()

	var wg sync.WaitGroup

	thermometer := &mocks.Thermometer{
		ReadFn: func() (float64, error) { return 20, nil },
	}

	chiller := &mocks.Actuator{
		OnFn: func() error {
			wg.Done()

			return nil
		},
		OffFn: func() error { return nil },
	}

	heater := &mocks.Actuator{
		OnFn:  func() error { return nil },
		OffFn: func() error { return nil },
	}

	therm := thermostat.NewThermostat(thermometer, chiller, heater, chillerKp, chillerKi, chillerKd, heaterKp,
		heaterKi, heaterKd, l)

	wg.Add(1)

	go func() {
		if err := therm.On(15); err != nil {
			t.Errorf("Unexpected error. Got: %+v", err)
		}

		wg.Done()
	}()

	wg.Wait() // wait for chiller.On to be called

	wg.Add(1)

	therm.Off()

	wg.Wait() // wait for chiller.On to be called
}

func TestDutyTimeLessThanMinimum(t *testing.T) {
	l, hook := logtest.NewNullLogger()
	l.Level = logrus.DebugLevel

	var wg sync.WaitGroup

	thermometer := &mocks.Thermometer{
		ReadFn: func() (float64, error) { return 16, nil },
	}

	chiller := &mocks.Actuator{
		OnFn: func() error {
			wg.Done()

			return nil
		},
		OffFn: func() error { return nil },
	}

	heater := &mocks.Actuator{
		OnFn:  func() error { return nil },
		OffFn: func() error { return nil },
	}

	therm := thermostat.NewThermostat(thermometer, chiller, heater, chillerKp, chillerKi, chillerKd, heaterKp,
		heaterKi, heaterKd, l)

	wg.Add(1)

	go func() {
		if err := therm.On(15); err != nil {
			t.Errorf("Unexpected error. Got: %+v", err)
		}
	}()

	wg.Wait()

	if !logContains(hook.AllEntries(), logrus.DebugLevel, "Forcing chiller actuator to a run for a minimum") {
		t.Errorf("Expected '%s' to be logged", "Forcing chiller actuator to a run for a minimum")
	}
}

func TestOnAlreadyOnError(t *testing.T) {
	l, _ := logtest.NewNullLogger()

	var wg sync.WaitGroup

	thermometer := &mocks.Thermometer{
		ReadFn: func() (float64, error) { return 20, nil },
	}

	chiller := &mocks.Actuator{
		OnFn: func() error {
			wg.Done()

			return nil
		},
		OffFn: func() error { return nil },
	}

	heater := &mocks.Actuator{
		OnFn:  func() error { return nil },
		OffFn: func() error { return nil },
	}

	therm := thermostat.NewThermostat(thermometer, chiller, heater, chillerKp, chillerKi, chillerKd, heaterKp,
		heaterKi, heaterKd, l)

	wg.Add(1)

	go func() {
		if err := therm.On(15); err != nil {
			t.Errorf("Unexpected error. Got: %+v", err)
		}
	}()

	wg.Wait()

	err := therm.On(66)
	if !errors.Is(err, thermostat.ErrAlreadyOn) {
		t.Errorf("Unexpected error. Want: '%s', Got: '%s'", thermostat.ErrAlreadyOn, err)
	}
}

func TestThermometerError(t *testing.T) {
	l, _ := logtest.NewNullLogger()

	thermometer := &mocks.Thermometer{
		ReadFn: func() (float64, error) { return 0, errDeadThermometer },
	}

	chiller := &mocks.Actuator{
		OnFn:  func() error { return nil },
		OffFn: func() error { return nil },
	}

	heater := &mocks.Actuator{
		OnFn:  func() error { return nil },
		OffFn: func() error { return nil },
	}

	therm := thermostat.NewThermostat(thermometer, chiller, heater, chillerKp, chillerKi, chillerKd, heaterKp,
		heaterKi, heaterKd, l)

	err := therm.On(15)
	if !errors.Is(err, errDeadThermometer) {
		t.Errorf("Unexpected error. Want: '%s', Got: '%s'", errDeadThermometer, err)
	}
}

func TestActuatorOnError(t *testing.T) {
	l, _ := logtest.NewNullLogger()

	thermometer := &mocks.Thermometer{
		ReadFn: func() (float64, error) { return 20, nil },
	}

	chiller := &mocks.Actuator{
		OnFn:  func() error { return errDeadActuator },
		OffFn: func() error { return nil },
	}

	heater := &mocks.Actuator{
		OnFn:  func() error { return nil },
		OffFn: func() error { return nil },
	}

	therm := thermostat.NewThermostat(thermometer, chiller, heater, chillerKp, chillerKi, chillerKd, heaterKp,
		heaterKi, heaterKd, l)

	err := therm.On(15)
	if !errors.Is(err, errDeadActuator) {
		t.Errorf("Unexpected error. Want: '%s', Got: '%s'", errDeadActuator, err)
	}
}

func TestActuatorOffErrorAfterDuty(t *testing.T) {
	l, _ := logtest.NewNullLogger()

	thermometer := &mocks.Thermometer{
		ReadFn: func() (float64, error) { return 20, nil },
	}

	chiller := &mocks.Actuator{
		OnFn:  func() error { return nil },
		OffFn: func() error { return errDeadActuator },
	}

	heater := &mocks.Actuator{
		OnFn:  func() error { return nil },
		OffFn: func() error { return nil },
	}

	therm := thermostat.NewThermostat(thermometer, chiller, heater, chillerKp, chillerKi, chillerKd, heaterKp,
		heaterKi, heaterKd, l, thermostat.SetChillingCyclePeriod(100*time.Millisecond),
		thermostat.SetChillingMinimum(10*time.Millisecond))

	err := therm.On(15)
	if !errors.Is(err, errDeadActuator) {
		t.Errorf("Unexpected error. Want: '%s', Got: '%s'", errDeadActuator, err)
	}
}

func TestActuatorOffErrorOnQuit(t *testing.T) {
	l, _ := logtest.NewNullLogger()

	var wg sync.WaitGroup

	thermometer := &mocks.Thermometer{
		ReadFn: func() (float64, error) { return 20, nil },
	}

	chiller := &mocks.Actuator{
		OnFn: func() error {
			wg.Done()

			return nil
		},
		OffFn: func() error { return errDeadActuator },
	}

	heater := &mocks.Actuator{
		OnFn:  func() error { return nil },
		OffFn: func() error { return nil },
	}

	therm := thermostat.NewThermostat(thermometer, chiller, heater, chillerKp, chillerKi, chillerKd, heaterKp,
		heaterKi, heaterKd, l)

	wg.Add(1)

	go func() {
		err := therm.On(15)
		if !errors.Is(err, errDeadActuator) {
			t.Errorf("Unexpected error. Want: '%s', Got: '%s'", errDeadActuator, err)
		}

		wg.Done()
	}()

	wg.Wait() // wait for chiller.On to be called

	wg.Add(1)

	therm.Off()

	wg.Wait() // wait for chiller.On to be called
}

func logContains(logs []*logrus.Entry, level logrus.Level, substr string) bool {
	found := false

	for _, v := range logs {
		// fmt.Println(v.Level, v.Message)
		if strings.Contains(v.Message, substr) && v.Level == level {
			found = true
		}
	}

	return found
}
