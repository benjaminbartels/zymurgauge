package simulation

import (
	"time"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"github.com/benjaminbartels/zymurgauge/internal/platform/temporal"
	"github.com/felixge/pidctrl"
)

type Test struct {
	Name       string
	Result     Result
	chamber    *Chamber
	speed      float64
	clock      temporal.Clock
	targetTemp float64
	start      time.Time
	logger     log.Logger
}

type Result struct {
	Durations []time.Duration
	Temps     []float64
	Targets   []float64
}

func NewTest(name string, initalTemp float64, targetTemp float64, interval time.Duration, minimumCooling time.Duration,
	minimumHeating time.Duration, p float64, i float64, d float64, speed float64,
	logger log.Logger) (*Test, error) {

	thermometer := NewThermometer(initalTemp)
	chiller := &Actuator{ActuatorType: Chiller}
	heater := &Actuator{ActuatorType: Heater}
	pidCtrl := pidctrl.NewPIDController(p, i, d)

	thermostat := &internal.Thermostat{
		ChillerPin:    "1",
		HeaterPin:     "2",
		ThermometerID: "abc123",
	}

	clock := temporal.NewDilatedClock(speed)

	if err := thermostat.Configure(pidCtrl,
		thermometer, chiller, heater,
		internal.MinimumCooling(minimumCooling),
		internal.MinimumHeating(minimumHeating),
		internal.Interval(interval),
		internal.Logger(logger),
		internal.Clock(clock)); err != nil {
		return nil, err
	}
	chamber := NewChamber(thermostat, thermometer, chiller, heater, speed, logger)

	t := &Test{
		Name:       name,
		chamber:    chamber,
		targetTemp: targetTemp,
		speed:      speed,
		clock:      clock,
		logger:     logger,
	}

	chamber.Thermostat.Subscribe(name, t.processStatus)

	return t, nil

}

func (t *Test) Run(duration time.Duration) Result {
	t.start = time.Now()

	t.chamber.Thermostat.Set(t.targetTemp)

	t.chamber.Thermostat.On()

	<-time.After(duration)

	t.chamber.Thermostat.Off()

	return t.Result
}

func (t *Test) processStatus(s internal.ThermostatStatus) {

	if s.Error != nil {
		t.logger.Fatal(s.Error)
	} else {
		t.Result.Durations = append(t.Result.Durations, time.Duration(float64(time.Since(t.start))*t.speed))
		t.Result.Temps = append(t.Result.Temps, *(s.CurrentTemperature))
		t.Result.Targets = append(t.Result.Targets, t.targetTemp)
	}
}
