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
	speed := 3600.0

	tests := []*simulation.Test{}
	for i := 1; i <= 1; i++ {
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
			_ = test.Run(5 * time.Second)
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
				{0.0, "00:00"},
				{0.5, "00:30"},
				{1.0, "01:00"},
				{1.5, "01:30"},
				{2.0, "02:00"},
				{2.5, "02:30"},
				{3.0, "03:00"},
				{3.5, "03:30"},
				{4.0, "04:00"},
				{4.5, "04:30"},
				{5.0, "05:00"},
			},
			//	ValueFormatter: chart.TimeMinuteValueFormatter,
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
