package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/benjaminbartels/zymurgauge/cmd/zymsim/simulator"
	"github.com/benjaminbartels/zymurgauge/internal/test/fakes"
	"github.com/benjaminbartels/zymurgauge/internal/thermostat"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	numOfJobs    = 1000
	numOfWorkers = 20
)

type simulationJob struct {
	*globals
	*simulation
}

type simulationResults struct {
	ChillerKP          float64
	ChillerKI          float64
	ChillerKD          float64
	HeaterKP           float64
	HeaterKI           float64
	HeaterKD           float64
	TemperatureAverage float64
}

type simulationSet struct {
	StartingTemp float64 `kong:"arg,help='Starting temperature.'"`
	TargetTemp   float64 `kong:"arg,help='Target temperature.'"`
	ChillerKPMin float64 `kong:"arg,help='Chiller proportional gain (kP) minimum.'"`
	ChillerKPMax float64 `kong:"arg,help='Chiller proportional gain (kP) maximum.'"`
	ChillerKIMin float64 `kong:"arg,help='Chiller integral gain (kI) minimum.'"`
	ChillerKIMax float64 `kong:"arg,help='Chiller integral gain (kI) maximum.'"`
	ChillerKDMin float64 `kong:"arg,help='Chiller derivative gain (kd) minimum.'"`
	ChillerKDMax float64 `kong:"arg,help='Chiller derivative gain (kd) maximum.'"`
	HeaterKPMin  float64 `kong:"arg,help='Heater proportional gain (kP) minimum.'"`
	HeaterKPMax  float64 `kong:"arg,help='Heater proportional gain (kP) maximum.'"`
	HeaterKIMin  float64 `kong:"arg,help='Heater integral gain (kI) minimum.'"`
	HeaterKIMax  float64 `kong:"arg,help='Heater integral gain (kI) maximum.'"`
	HeaterKDMin  float64 `kong:"arg,help='Heater derivative gain (kD) minimum.'"`
	HeaterKDMax  float64 `kong:"arg,help='Heater derivative gain (kD) maximum.'"`
	FileName     string  `kong:"arg,optional,help='Name of results file. Defaults to results_{timestamp}.csv.'"`
}

func (s *simulationSet) Run(g *globals) error {
	logger := logrus.New()

	if g.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	jobs := make(chan simulationJob, numOfJobs)
	results := make(chan simulationResults, numOfJobs)

	var wg sync.WaitGroup

	for w := 0; w < numOfWorkers; w++ {
		wg.Add(1)

		go worker(jobs, results, logger, &wg)
	}

	for chillerKP := s.ChillerKPMin; chillerKP <= s.ChillerKPMax; chillerKP++ {
		for chillerKI := s.ChillerKIMin; chillerKI <= s.ChillerKIMax; chillerKI++ {
			for chillerKD := s.ChillerKDMin; chillerKD <= s.ChillerKDMax; chillerKD++ {
				for heaterKP := s.HeaterKPMin; heaterKP <= s.HeaterKPMax; heaterKP++ {
					for heaterKI := s.HeaterKIMin; heaterKI <= s.HeaterKIMax; heaterKI++ {
						for heaterKD := s.HeaterKDMin; heaterKD <= s.HeaterKDMax; heaterKD++ {
							jobs <- simulationJob{
								globals: g,
								simulation: &simulation{
									StartingTemp: s.StartingTemp,
									TargetTemp:   s.TargetTemp,
									ChillerKP:    chillerKP,
									ChillerKI:    chillerKI,
									ChillerKD:    chillerKD,
									HeaterKP:     heaterKP,
									HeaterKI:     heaterKI,
									HeaterKD:     heaterKD,
								},
							}
						}
					}
				}
			}
		}
	}
	close(jobs)

	var simulationResults []simulationResults

	go func() {
		for {
			result := <-results
			fmt.Println(result)
			simulationResults = append(simulationResults, result)
		}
	}()

	wg.Wait()

	file, err := os.Create("result.csv")
	if err != nil {
		return errors.Wrap(err, "could not create csv file")
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, result := range simulationResults {
		if err := writer.Write([]string{
			fmt.Sprintf("%f", result.ChillerKP), fmt.Sprintf("%f", result.ChillerKI),
			fmt.Sprintf("%f", result.ChillerKD), fmt.Sprintf("%f", result.HeaterKP),
			fmt.Sprintf("%f", result.HeaterKI), fmt.Sprintf("%f", result.HeaterKD),
			fmt.Sprintf("%f", result.TemperatureAverage),
		}); err != nil {
			return errors.Wrap(err, "could not write results to csv file")
		}
	}

	return nil
}

func worker(jobs <-chan simulationJob, results chan<- simulationResults, logger *logrus.Logger, wg *sync.WaitGroup) {
	for j := range jobs {
		sim := simulator.New(j.StartingTemp)
		clock := fakes.NewDilatedClock(j.Multiplier)
		thermostat := thermostat.NewThermostat(sim.Thermometer, sim.Chiller, sim.Heater,
			j.ChillerKP, j.ChillerKI, j.ChillerKD, j.HeaterKP, j.HeaterKI, j.HeaterKD,
			logger, thermostat.SetClock(clock))
		ctx, stop := context.WithCancel(context.Background())
		startTime := time.Now()

		go runSimulator(ctx, sim, j.Multiplier)

		go func() {
			if err := thermostat.On(j.TargetTemp); err != nil {
				fmt.Println("Failed to turn thermostat on:", err)
				os.Exit(1)
			}
		}()

		readings := make(chan reading)

		go runTemperatureReader(ctx, sim.Thermometer, startTime, j.Multiplier, readings)

		go func() {
			<-time.After(j.Runtime)
			thermostat.Off()
			stop()
		}()

		var (
			temps float64
			count int
		)

		for reading := range readings {
			temps += reading.Temp
			count++
		}

		results <- simulationResults{
			ChillerKP:          j.ChillerKP,
			ChillerKI:          j.ChillerKI,
			ChillerKD:          j.ChillerKD,
			HeaterKP:           j.HeaterKP,
			HeaterKI:           j.HeaterKI,
			HeaterKD:           j.HeaterKD,
			TemperatureAverage: temps / float64(count),
		}
	}

	wg.Done()
}
