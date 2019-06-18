package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/simulation"

	chart "github.com/wcharczuk/go-chart"
)

func main() {

	logger := log.New(os.Stderr, "", log.LstdFlags)

	target := 18.0

	multiplier := 43200.0 * 2

	tests := []*simulation.Test{}
	for i := 1; i <= 5; i++ {
		test, err := simulation.NewTest("Test 1", 20, target, 10*time.Minute, 1*time.Minute, 1*time.Minute,
			float64(i), 0, 0, multiplier, logger)
		if err != nil {
			panic(err)
		}
		tests = append(tests, test)
	}

	results := make([]simulation.Results, len(tests))

	var wg sync.WaitGroup

	for i, test := range tests {
		wg.Add(1)
		go func(i int, test *simulation.Test) {
			defer wg.Done()
			results[i] = test.Run(1 * time.Second)
		}(i, test)
	}

	wg.Wait()

	err := createGraph(results, target)
	if err != nil {
		panic(err)
	}

	fmt.Println("Bye!")
}

func createGraph(results []simulation.Results, target float64) error {

	series := []chart.Series{}

	for _, result := range results {

		times := []time.Time{}

		// for i := range result.Times {
		// 	fmt.Println(result.Times[i], result.Temps[i])
		// }

		for _, duration := range result.Durations {
			times = append(times, time.Time{}.Add(duration))
		}

		s := chart.TimeSeries{
			XValues: times,
			YValues: result.Temps,
		}
		series = append(series, s)
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
			GridMajorStyle: chart.Style{
				Show:        true,
				StrokeColor: chart.ColorAlternateGray,
				StrokeWidth: 1.0,
			},
			GridLines: []chart.GridLine{
				{Value: target},
			},
		},
		Series: series,
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
