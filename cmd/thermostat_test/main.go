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

	initialTemp := 20.0
	targetTemp := 18.0
	speed := 43200.0 * 2

	tests := []*simulation.Test{}
	for i := 1; i <= 10; i++ {
		test, err := simulation.NewTest(fmt.Sprintf("P = %d", i), initialTemp, targetTemp, 10*time.Minute, 1*time.Minute, 1*time.Minute,
			float64(i), 0, 0, speed, logger)
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
			_ = test.Run(1 * time.Second)
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

func createGraph(tests []*simulation.Test, targetTemp float64) error {

	series := []chart.Series{}

	for _, test := range tests {

		times := []time.Time{}

		for _, duration := range test.Result.Durations {
			fmt.Println(duration, time.Time{}.Add(duration))
			times = append(times, time.Time{}.Add(duration))
		}

		s := chart.TimeSeries{
			Name:    test.Name,
			XValues: times,
			YValues: test.Result.Temps,
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
				{Value: targetTemp},
			},
		},
		Series: series,
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
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
