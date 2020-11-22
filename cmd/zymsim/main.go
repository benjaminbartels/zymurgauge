package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sync"
	"time"

	"github.com/alecthomas/kong"
	"github.com/benjaminbartels/zymurgauge/cmd/zymsim/simulator"
	"github.com/benjaminbartels/zymurgauge/internal/test/fakes"
	"github.com/benjaminbartels/zymurgauge/internal/thermostat"
	"github.com/sirupsen/logrus"
	"github.com/wcharczuk/go-chart"
)

const (
	graphInterval    = 9
	graphStrokeWidth = 1.0
)

//nolint:gochecknoglobals
var (
	readInterval = 100 * time.Millisecond
)

type CLI struct {
	Multiplier   float64       `kong:"default=6000.0,short=m,help='Time dilation multiplier. Defaults to 6000.'"`
	Runtime      time.Duration `kong:"default=5s,short=r,help='Runtime of simulation. Defaults to 5s.'"`
	Debug        bool          `kong:"default=false,short=l,help='Enable debug logging. Default is false.'"`
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

	clock := fakes.NewDilatedClock(cli.Multiplier)
	logger := logrus.New()

	if cli.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	sim := simulator.New(cli.StartingTemp)
	thermostat := thermostat.NewThermostat(sim.Thermometer, sim.Chiller, sim.Heater,
		cli.ChillerKp, cli.ChillerKi, cli.ChillerKd, cli.HeaterKp, cli.HeaterKi, cli.HeaterKd,
		logger, thermostat.SetClock(clock))

	ctx, stop := context.WithCancel(context.Background())

	go run(ctx, sim, cli.Multiplier)

	startTime := time.Now()

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		if err := thermostat.On(cli.TargetTemp); err != nil {
			fmt.Println("Failed to turn thermostat on:", err)
			os.Exit(1)
		}

		wg.Done()
	}()

	durations := []time.Duration{}
	temps := []float64{}

	wg.Add(1)

	go func() {
		var err error

		durations, temps, err = read(ctx, sim.Thermometer, startTime, cli.Multiplier)

		if err != nil {
			fmt.Println("Failed to get thermometer readings:", err)
			os.Exit(1)
		}

		wg.Done()
	}()

	<-time.After(cli.Runtime)
	thermostat.Off()
	stop()
	wg.Wait()

	if err := createGraph(durations, temps, cli.TargetTemp, cli.FileName); err != nil {
		fmt.Println("Failed to create graph", err)
		os.Exit(1)
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

func run(ctx context.Context, simulator *simulator.Simulator, multiplier float64) {
	updateCycle := time.Duration(int64(math.Round(1 / multiplier * (1e9))))
	tick := time.Tick(updateCycle)

	for {
		select {
		case <-tick:
			simulator.Update()
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
