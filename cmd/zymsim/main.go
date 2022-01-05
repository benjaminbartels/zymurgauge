package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/benjaminbartels/zymurgauge/cmd/zymsim/simulator"
	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/benjaminbartels/zymurgauge/internal/device/pid"
	"github.com/benjaminbartels/zymurgauge/internal/test/fakes"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wcharczuk/go-chart"
)

const (
	graphInterval            = 9
	graphStrokeWidth         = 1.0
	chartFilePermissions     = 0o600
	readIntervalMilliseconds = 100
)

type reading struct {
	Duration time.Duration
	Temp     float64
}

type cli struct {
	Multiplier   float64       `kong:"default=6000.0,short=m,help='Time dilation multiplier. Defaults to 6000.'"`
	Runtime      time.Duration `kong:"default=5s,short=r,help='Runtime of simulation. Defaults to 5s.'"`
	Debug        bool          `kong:"default=false,short=d,help='Enable debug logging. Default is false.'"`
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

func main() {
	logger := logrus.New()
	if err := run(logger); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}

func run(logger *logrus.Logger) error {
	cli := cli{}
	kong.Parse(&cli,
		kong.Name("zymsim"),
		kong.Description("Zymurgauge Thermostat Simulator"),
		kong.UsageOnError(),
	)

	if cli.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	sim := simulator.New(cli.StartingTemp)
	clock := fakes.NewDilatedClock(cli.Multiplier)
	pid := pid.NewPIDTemperatureController(sim.Thermometer, sim.Chiller, sim.Heater,
		cli.ChillerKp, cli.ChillerKi, cli.ChillerKd, cli.HeaterKp, cli.HeaterKi, cli.HeaterKd,
		logger, pid.SetClock(clock))
	ctx, stop := context.WithCancel(context.Background())
	startTime := time.Now()

	go runSimulator(ctx, sim, cli.Multiplier)

	go func() {
		if err := pid.Run(ctx, cli.TargetTemp); err != nil {
			logger.Error("Failed to turn pid on:", err)
			os.Exit(1)
		}
	}()

	readings := make(chan reading)

	go runTemperatureReader(ctx, sim.Thermometer, startTime, cli.Multiplier, readings)

	go func() {
		<-time.After(cli.Runtime)
		stop()
		close(readings)
	}()

	durations := []time.Duration{}
	temps := []float64{}

	for reading := range readings {
		durations = append(durations, reading.Duration)
		temps = append(temps, reading.Temp)
	}

	if err := createGraph(durations, temps, cli.TargetTemp, cli.FileName); err != nil {
		return errors.Wrap(err, "could not create graph")
	}

	return nil
}

func runSimulator(ctx context.Context, simulator *simulator.Simulator, multiplier float64) {
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

func runTemperatureReader(ctx context.Context, thermometer device.Thermometer, startTime time.Time,
	multiplier float64, readings chan reading) {
	tick := time.Tick(readIntervalMilliseconds * time.Millisecond)

	for {
		select {
		case <-tick:
			temp, err := thermometer.GetTemperature()
			if err != nil {
				os.Exit(1)
			}
			readings <- reading{
				Duration: time.Duration(float64(time.Since(startTime)) * multiplier),
				Temp:     temp,
			}
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
		d -= hour * time.Hour //nolint: durationcheck // TODO: fix durationcheck later
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
		return errors.Wrap(err, "could not render chart")
	}

	readBuf, err := ioutil.ReadAll(buffer)
	if err != nil {
		return errors.Wrap(err, "could not real buffer")
	}

	if fileName == "" {
		fileName = "chart_" + time.Now().Format("20060102150405") + ".png"
	}

	if err := ioutil.WriteFile(fileName, readBuf, chartFilePermissions); err != nil {
		return errors.Wrap(err, "could not write file")
	}

	return nil
}
