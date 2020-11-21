package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"sync"
	"time"

	"github.com/alecthomas/kong"
	"github.com/benjaminbartels/zymurgauge/internal/test/fakes"
	"github.com/benjaminbartels/zymurgauge/internal/thermostat"
	"github.com/sirupsen/logrus"
	"github.com/wcharczuk/go-chart"
)

const (
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

type CLI struct {
	Multiplier   float64       `kong:"default=6000.0,short=m,help='Time dilation multiplier. Defaults to 6000.'"`
	Runtime      time.Duration `kong:"default=5s,short=r,help='Runtime of simulation. Defaults to 5s.'"`
	Log          bool          `kong:"default=false,short=l,help='Enable logger. Default is false.'"`
	StartingTemp float64       `kong:"arg,help='Starting temperature.'"`           // 25.0
	TargetTemp   float64       `kong:"arg,help='Target temperature.'"`             // 20.0
	ChillerKp    float64       `kong:"arg,help='Chiller proportional gain (kP).'"` // -1.0
	ChillerKi    float64       `kong:"arg,help='Chiller integral gain (kI).'"`     // 0.0
	ChillerKd    float64       `kong:"arg,help='Chiller derivative gain (kd).'"`   // 0.0
	HeaterKp     float64       `kong:"arg,help='Heater proportional gain (kP).'"`  // 1.0
	HeaterKi     float64       `kong:"arg,help='Heater integral gain (kI).'"`      // 0.0
	HeaterKd     float64       `kong:"arg,help='Heater derivative gain (kD).'"`    // 0.0
	FileName     string        `kong:"arg,optional,help='Name of results file. Defaults to chart_{timestamp}.png.'"`
}

//nolint:funlen
func main() {
	cli := CLI{}
	kong.Parse(&cli,
		kong.Name("zymsim"),
		kong.Description("Zymurgauge Thermostat Simulator"),
		kong.UsageOnError(),
	)

	thermometer := &thermometer{currentTemp: cli.StartingTemp}
	chiller := &actuator{}
	heater := &actuator{}
	clock := fakes.NewDilatedClock(cli.Multiplier)
	logger := logrus.New()

	if cli.Log {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.Out = ioutil.Discard
	}

	thermostat := thermostat.NewThermostat(thermometer, chiller, heater, cli.ChillerKp, cli.ChillerKi, cli.ChillerKd,
		cli.HeaterKp, cli.HeaterKi, cli.HeaterKd, logger, thermostat.SetClock(clock))

	ctx, stop := context.WithCancel(context.Background())

	go run(ctx, thermometer, chiller, heater, cli.Multiplier)

	startTime := time.Now()

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		if err := thermostat.On(cli.TargetTemp); err != nil {
			logger.WithError(err).Warn("Thermostat failed to turn on")
		}

		wg.Done()
	}()

	durations := []time.Duration{}
	temps := []float64{}

	wg.Add(1)

	go func() {
		var err error

		durations, temps, err = read(ctx, thermometer, startTime, cli.Multiplier)

		if err != nil {
			logger.WithError(err).Warn("Failed to get thermometer readings")
		}

		wg.Done()
	}()

	<-time.After(cli.Runtime)
	thermostat.Off()
	stop()
	wg.Wait()

	if err := createGraph(durations, temps, cli.TargetTemp, cli.FileName); err != nil {
		logger.WithError(err).Warn("Failed to create graph")
	}
}

func read(ctx context.Context, thermometer thermostat.Thermometer, startTime time.Time,
	multiplier float64) ([]time.Duration, []float64, error) {
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

			d := time.Duration(float64(time.Since(startTime)) * multiplier)
			durations = append(durations, d)

			temps = append(temps, temp)
		case <-ctx.Done():
			return durations, temps, nil
		}
	}
}

func run(ctx context.Context, thermometer *thermometer, chiller, heater *actuator, multiplier float64) {
	updateCycle := time.Duration(int64(math.Round(1 / multiplier * (1e9))))
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
			return
		}
	}
}

func createGraph(durations []time.Duration, temps []float64, targetTemp float64, fileName string) error {
	times := []float64{}
	ticks := make([]chart.Tick, len(durations))
	maxDuration := durations[len(durations)-1]
	interval := maxDuration / graphInterval
	numOfTick := 10

	for _, duration := range durations {
		times = append(times, float64(duration))
	}

	if len(durations) < numOfTick {
		numOfTick = len(durations)
	}

	for i := 0; i < numOfTick; i++ {
		tickValue := int64(interval) * int64(i)
		d := time.Duration(tickValue).Round(time.Minute)
		hour := d / time.Hour
		d -= hour * time.Hour
		minute := d / time.Minute
		ticks[i] = chart.Tick{Value: float64(tickValue), Label: fmt.Sprintf("%02d:%02d", hour, minute)}
	}

	series := []chart.Series{chart.ContinuousSeries{
		XValues: times,
		YValues: temps,
	}}

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

	return writeChart(graph, fileName)
}

func writeChart(c chart.Chart, fileName string) error {
	buffer := bytes.NewBuffer([]byte{})
	if err := c.Render(chart.PNG, buffer); err != nil {
		return err
	}

	readBuf, err := ioutil.ReadAll(buffer)
	if err != nil {
		return err
	}

	if fileName == "" {
		fileName = "chart_" + time.Now().Format("20060102150405") + ".png"
	}

	return ioutil.WriteFile(fileName, readBuf, 0600)
}
