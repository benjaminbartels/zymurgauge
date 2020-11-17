//nolint:gomnd
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/simulation"
	"github.com/sirupsen/logrus"

	chart "github.com/wcharczuk/go-chart"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	chillingMinimum := 1 * time.Minute
	heatingMinimum := 1 * time.Minute
	chillerCyclePeriod := 10 * time.Minute
	heaterCyclePeriod := 10 * time.Minute
	runTime := 5 * time.Second

	initialTemp := 20.0
	targetTemp := 18.0
	speed := 3600.0

	tests := []*simulation.Test{}

	for p := 1; p <= 1; p++ {
		test, err := simulation.NewTest(fmt.Sprintf("P = %d, %d", -p, p), chillingMinimum, heatingMinimum,
			chillerCyclePeriod, heaterCyclePeriod, float64(-p), 0, 0, float64(p), 0, 0, speed, initialTemp, targetTemp,
			logger)
		if err != nil {
			panic(err)
		}

		tests = append(tests, test)
	}

	results := make([]*simulation.Test, len(tests))

	var wg sync.WaitGroup

	for i, test := range tests {
		wg.Add(1)

		go func(i int, test *simulation.Test) {
			defer wg.Done()

			_ = test.Run(runTime)
			results[i] = test
		}(i, test)
	}

	wg.Wait()

	err := createGraph(results, targetTemp)
	if err != nil {
		panic(err)
	}

	fmt.Println("Bye!")
}

//nolint:funlen
func createGraph(tests []*simulation.Test, targetTemp float64) error {
	series := []chart.Series{}

	var maxTime float64

	for _, test := range tests {
		times := []float64{}

		for _, duration := range test.Result.Durations {
			time := float64(duration) / 3600000000000.0
			if time > maxTime {
				maxTime = time
			}

			times = append(times, time)
		}

		s := chart.ContinuousSeries{
			Name:    test.Name,
			XValues: times,
			YValues: test.Result.Temps,
		}

		series = append(series, s)
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:      "Hours",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			Ticks: []chart.Tick{
				{Value: 0.0, Label: "00:00"},
				{Value: 0.5, Label: "00:30"},
				{Value: 1.0, Label: "01:00"},
				{Value: 1.5, Label: "01:30"},
				{Value: 2.0, Label: "02:00"},
				{Value: 2.5, Label: "02:30"},
				{Value: 3.0, Label: "03:00"},
				{Value: 3.5, Label: "03:30"},
				{Value: 4.0, Label: "04:00"},
				{Value: 4.5, Label: "04:30"},
				{Value: 5.0, Label: "05:00"},
			},
		},
		YAxis: chart.YAxis{
			Name:      "Temperature",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			GridMajorStyle: chart.Style{
				Show:        true,
				StrokeColor: chart.ColorAlternateGray,
				StrokeWidth: 1.0,
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

	filename := "chart_" + time.Now().Format("20060102150405") + ".png"

	err = ioutil.WriteFile(filename, readBuf, 0600)

	return err
}
