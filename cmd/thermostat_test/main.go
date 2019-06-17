package main

// ToDo: make this a unit test or part of fermmon
import (
	"bytes"
	"fmt"
	"io/ioutil"
	golog "log"
	"os"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"github.com/benjaminbartels/zymurgauge/internal/platform/temporal"
	"github.com/benjaminbartels/zymurgauge/internal/simulation"
	"github.com/felixge/pidctrl"
	chart "github.com/wcharczuk/go-chart"
)

var logger log.Logger
var times []time.Time
var temps, targets []float64
var dilatedClock temporal.Clock

const target = 18.00
const factor = 6000
const interval = 10 * time.Minute
const minimun = 1 * time.Minute
const testDuration = 10 * time.Second

func main() {

	logger = golog.New(os.Stderr, "", golog.LstdFlags)
	thermometer := simulation.NewThermometer(20)
	chiller := &simulation.Actuator{ActuatorType: simulation.Chiller}
	heater := &simulation.Actuator{ActuatorType: simulation.Heater}
	pidCtrl := pidctrl.NewPIDController(1, 0, 0)
	pidCtrl.SetOutputLimits(-10, 10)

	thermostat := &internal.Thermostat{
		ChillerPin:    "1",
		HeaterPin:     "2",
		ThermometerID: "test",
	}

	dilatedClock = temporal.NewDilatedClock(factor)

	err := thermostat.Configure(pidCtrl, thermometer, chiller, heater,
		internal.MinimumCooling(minimun),
		internal.MinimumHeating(minimun),
		internal.Interval(interval),
		internal.Logger(logger),
		internal.Clock(dilatedClock),
	)
	if err != nil {
		panic(err)
	}

	thermostat.Subscribe("test", processStatus)

	chamber := simulation.NewChamber(thermostat, thermometer, chiller, heater, factor, logger)
	chamber.Thermostat.Set(target)
	chamber.Thermostat.On()

	time.Sleep(testDuration)

	chamber.Thermostat.Off()

	err = createGraph(times, temps, targets)
	if err != nil {
		panic(err)
	}

	fmt.Println("Bye!")
}

func processStatus(s internal.ThermostatStatus) {

	if s.Error != nil {
		logger.Fatal(s.Error)
	} else {
		logger.Println("Event:", s.State, *(s.CurrentTemperature))

		times = append(times, dilatedClock.Now())
		temps = append(temps, *(s.CurrentTemperature))
		targets = append(targets, target)
	}
}

func createGraph(x []time.Time, y []float64, targets []float64) error {

	for i := range x {
		fmt.Println(x[i], y[i])
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
			ValueFormatter: chart.TimeMinuteValueFormatter,
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		Series: []chart.Series{
			chart.TimeSeries{
				XValues: x,
				YValues: y,
			},
			chart.TimeSeries{
				XValues: x,
				YValues: targets,
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		return err
	}

	readBuf, err := ioutil.ReadAll(buffer)
	if err != nil {
		return err
	}

	filename := "chart_" + time.Now().Format("20060102150405") + ".png"

	err = ioutil.WriteFile(filename, readBuf, 0644)

	return err
}
