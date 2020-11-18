package simulation

import (
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/test/fakes"
	"github.com/benjaminbartels/zymurgauge/internal/thermostat"
	"github.com/sirupsen/logrus"
)

// TODO: make this a unit test

type Test struct {
	Name       string
	Result     Result
	chamber    *Chamber
	speed      float64
	clock      thermostat.Clock
	targetTemp float64
	start      time.Time
	logger     *logrus.Logger
}

type Result struct {
	Durations []time.Duration
	Temps     []float64
	Targets   []float64
}

func NewTest(name string, chillingMinimum, heatingMinimum, chillerCyclePeriod, heaterCyclePeriod time.Duration,
	chillerKp, chillerKi, chillerKd, heaterKp, heaterKi, heaterKd, speed, initalTemp, targetTemp float64,
	logger *logrus.Logger) (*Test, error) {
	thermometer := NewThermometer(initalTemp)
	chiller := &Actuator{ActuatorType: Chiller}
	heater := &Actuator{ActuatorType: Heater}
	clock := fakes.NewDilatedClock(speed)

	thermostat := thermostat.NewThermostat(thermometer, chiller, heater,
		chillerKp, chillerKi, chillerKd, heaterKp, heaterKi, heaterKd, logger,
		thermostat.SetChillingMinimum(chillingMinimum),
		thermostat.SetHeatingMinimum(heatingMinimum),
		thermostat.SetChillingCyclePeriod(chillerCyclePeriod),
		thermostat.SetHeatingCyclePeriod(heaterCyclePeriod),
		thermostat.SetClock(clock),
	)

	chamber := NewChamber(thermostat, thermometer, chiller, heater, speed, logger)

	t := &Test{
		Name:       name,
		chamber:    chamber,
		targetTemp: targetTemp,
		speed:      speed,
		clock:      clock,
		logger:     logger,
	}

	return t, nil
}

func (t *Test) Run(runTime time.Duration) Result {
	t.start = time.Now()

	go func() {
		if err := t.chamber.Thermostat.On(t.targetTemp); err != nil {
			panic(err)
		}
	}()

	<-time.After(runTime)

	t.chamber.Thermostat.Off()

	return t.Result
}
