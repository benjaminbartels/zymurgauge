package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/mock"
	"github.com/felixge/pidctrl"
	chart "github.com/wcharczuk/go-chart"
)

var logger *log.Logger
var times []time.Time
var temps, targets []float64

const target = 18.00

func main() {

	logger = log.New(os.Stderr, "", log.LstdFlags)

	thermometer := mock.NewThermometer(20)
	chiller := &mock.Actuator{ActuatorType: mock.Chiller}
	heater := &mock.Actuator{ActuatorType: mock.Heater}
	pidCtrl := pidctrl.NewPIDController(20, 0, 0)
	pidCtrl.SetOutputLimits(-10, 10)

	thermostat, err := internal.NewThermostat(pidCtrl, thermometer, chiller, heater,
		internal.MinimumChill(1*time.Second),
		internal.MinimumHeat(1*time.Second),
		internal.Interval(1*time.Second), // 1sec = 10min
		internal.Factor(600),
		// internal.Interval(10*time.Second), // 10sec = 10min
		// internal.Factor(60),               // 10sec = 10min
		internal.Logger(logger),
	)

	if err != nil {
		panic(err)
	}

	chamber := mock.NewChamber(thermostat, thermometer, chiller, heater, 600, log.New(os.Stderr, "", log.LstdFlags))
	chamber.Thermostat.Subscribe("test", processStatus)
	chamber.Thermostat.Set(target)
	chamber.Thermostat.On()

	fmt.Println("Sleeping...")
	time.Sleep(20 * time.Second)

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
		logger.Println(s.State, *(s.CurrentTemperature))
		times = append(times, time.Now())
		temps = append(temps, *(s.CurrentTemperature))
		targets = append(targets, target)
	}
}

func createGraph(x []time.Time, y []float64, targets []float64) error {

	for i, _ := range x {
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

	err = ioutil.WriteFile("chart.png", readBuf, 0644)
	if err != nil {
		return err
	}

	return nil
}
