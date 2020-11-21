package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"sync"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/test/fakes"
	"github.com/benjaminbartels/zymurgauge/internal/thermostat"
	"github.com/sirupsen/logrus"
	"github.com/wcharczuk/go-chart"
)

const (
	startingTemp = 26
	targetTemp   = 20
	speed        = 6000.0 // 100ms = 10m
	chillerKp    = -1.0
	chillerKi    = 0.0
	chillerKd    = 0.0
	heaterKp     = 1.0
	heaterKi     = 0.0
	heaterKd     = 0.0
	beerCapacity = 4.2 * 1.0 * 20        // heat capacity water * density of water * 20L volume (in kJ per kelvin).
	airCapacity  = 1.005 * 1.225 * 0.200 // heat capacity of dry air * density of air * 200L volume (in kJ per kelvin).
	// Moist air has only slightly higher heat capacity, 1.02 when saturated at 20C.
	wallCapacity            = 5.0 // just a guess
	heaterCapacity          = 1.0 // also a guess, to simulate that heater first heats itself, then starts heating the air
	heaterPower             = 0.1 // 100W, in kW.
	coolerPower             = 0.1 // 100W, in kW. Assuming 200W at 50% efficiency
	airBeerTransfer         = 1.0 / 300
	wallAirTransfer         = 1.0 / 300
	heaterAirTransfer       = 1.0 / 30
	environmentWallTransfer = 0.001 // losses to environment
	// heaterToBeer            = 0.0   // ratio of heater transferred directly to beer instead of fridge air
	// heaterToAir             = 1.0 - heaterToBeer.
	graphInterval    = 9
	graphStrokeWidth = 1.0
)

//nolint:gochecknoglobals
var (
	runtime      = 5 * time.Second
	readInterval = 100 * time.Millisecond
)

type thermometer struct {
	currentTemp float64
}

func (t *thermometer) Read() (float64, error) {
	return t.currentTemp, nil
}

type actuator struct {
	isOn bool
}

func (a *actuator) On() error {
	if !a.isOn {
		a.isOn = true
	}

	return nil
}

func (a *actuator) Off() error {
	if a.isOn {
		a.isOn = false
	}

	return nil
}

func main() {
	thermometer := &thermometer{currentTemp: startingTemp}
	chiller := &actuator{}
	heater := &actuator{}
	clock := fakes.NewDilatedClock(speed)
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	thermostat := thermostat.NewThermostat(thermometer, chiller, heater,
		chillerKp, chillerKi, chillerKd, heaterKp, heaterKi, heaterKd, logger, thermostat.SetClock(clock))

	ctx, stop := context.WithCancel(context.Background())

	go run(ctx, thermometer, chiller, heater)

	startTime := time.Now()

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		if err := thermostat.On(targetTemp); err != nil {
			logger.WithError(err).Warn("Thermostat failed to turn on")
		}

		wg.Done()
	}()

	durations := []time.Duration{}
	temps := []float64{}

	wg.Add(1)

	go func() {
		var err error

		durations, temps, err = read(ctx, thermometer, startTime)

		if err != nil {
			logger.WithError(err).Warn("Failed to get thermometer readings")
		}

		wg.Done()
	}()

	<-time.After(runtime)

	thermostat.Off()

	stop()

	wg.Wait()

	if err := createGraph(durations, temps, targetTemp); err != nil {
		logger.WithError(err).Warn("Failed to create graph")
	}
}

func read(ctx context.Context, thermometer thermostat.Thermometer,
	startTime time.Time) ([]time.Duration, []float64, error) {
	durations := []time.Duration{}
	temps := []float64{}
	tick := time.Tick(readInterval)

	for {
		select {
		case <-tick:
			temp, err := thermometer.Read()
			if err != nil {
				return durations, temps, err
			}

			d := time.Duration(float64(time.Since(startTime)) * speed)
			durations = append(durations, d)

			temps = append(temps, temp)
		case <-ctx.Done():
			return durations, temps, nil
		}
	}
}

func run(ctx context.Context, thermometer *thermometer, chiller, heater *actuator) {
	updateCycle := time.Duration(int64(math.Round(1 / speed * (1e9))))
	tick := time.Tick(updateCycle)

	var (
		wallTemp        = 20.0
		airTemp         = 20.0
		beerTemp        = thermometer.currentTemp
		heaterTemp      = 20.0
		environmentTemp = 20.0
	)

	var ctr int

	for {
		select {
		case <-tick:
			ctr++

			beerTempNew := beerTemp
			airTempNew := airTemp
			wallTempNew := wallTemp
			heaterTempNew := heaterTemp

			beerTempNew += (airTemp - beerTemp) * airBeerTransfer / beerCapacity

			if chiller.isOn {
				wallTempNew -= coolerPower / wallCapacity
			} else if heater.isOn {
				heaterTempNew += heaterPower / heaterCapacity
			}

			airTempNew += (heaterTemp - airTemp) * heaterAirTransfer / airCapacity
			airTempNew += (wallTemp - airTemp) * wallAirTransfer / airCapacity
			airTempNew += (beerTemp - airTemp) * airBeerTransfer / airCapacity

			beerTempNew += (airTemp - beerTemp) * airBeerTransfer / beerCapacity

			heaterTempNew += (airTemp - heaterTemp) * heaterAirTransfer / heaterCapacity

			wallTempNew += (environmentTemp - wallTemp) * environmentWallTransfer / wallCapacity
			wallTempNew += (airTemp - wallTemp) * wallAirTransfer / wallCapacity

			airTemp = airTempNew
			beerTemp = beerTempNew
			wallTemp = wallTempNew
			heaterTemp = heaterTempNew
			thermometer.currentTemp = beerTemp
		case <-ctx.Done():
			fmt.Println("COUNT", ctr)

			return
		}
	}
}

func createGraph(durations []time.Duration, temps []float64, targetTemp float64) error {
	series := []chart.Series{}
	times := []float64{}
	ticks := make([]chart.Tick, len(durations))
	maxDuration := durations[len(durations)-1]
	interval := maxDuration / graphInterval

	fmt.Println("MAX", maxDuration)

	for _, duration := range durations {
		times = append(times, float64(duration))
	}

	for i := 0; i < 10; i++ {
		tickValue := int64(interval) * int64(i)
		d := time.Duration(tickValue).Round(time.Minute)
		hour := d / time.Hour
		d -= hour * time.Hour
		minute := d / time.Minute
		ticks[i] = chart.Tick{Value: float64(tickValue), Label: fmt.Sprintf("%02d:%02d", hour, minute)}
	}

	s := chart.ContinuousSeries{
		Name:    "myTest",
		XValues: times,
		YValues: temps,
	}

	series = append(series, s)

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:      "Time (hours)",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			Ticks:     ticks,
		},
		YAxis: chart.YAxis{
			Name:      "Temperature (C)",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			GridMajorStyle: chart.Style{
				Show:        true,
				StrokeColor: chart.ColorAlternateGray,
				StrokeWidth: graphStrokeWidth,
			},
			GridLines: []chart.GridLine{
				{Value: targetTemp},
			},
		},
		Series: series,
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	buffer := bytes.NewBuffer([]byte{})
	if err := graph.Render(chart.PNG, buffer); err != nil {
		return err
	}

	readBuf, err := ioutil.ReadAll(buffer)
	if err != nil {
		return err
	}

	return ioutil.WriteFile("chart_"+time.Now().Format("20060102150405")+".png", readBuf, 0600)
}
