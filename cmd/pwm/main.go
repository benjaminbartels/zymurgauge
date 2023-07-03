package main

import (
	"context"
	"os"
	"os/signal"

	// "time"

	"github.com/benjaminbartels/zymurgauge/internal/device/gpio"
	"github.com/benjaminbartels/zymurgauge/internal/device/onewire"
	"github.com/benjaminbartels/zymurgauge/internal/temperaturecontrol/pid"
	"github.com/sirupsen/logrus"
	"periph.io/x/host/v3"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	if _, err := host.Init(); err != nil {
		panic(err)
	}

	actuator, err := gpio.NewGPIOActuator("22")
	if err != nil {
		panic(err)
	}

	// if err := actuator.PWMOn(0.5); err != nil {
	// 	panic(err)
	// }

	// <-time.After(10 * time.Second)

	// if err := actuator.Off(); err != nil {
	// 	panic(err)
	// }

	therm, err := onewire.New("28-000006285484")
	if err != nil {
		panic(err)
	}

	pid := pid.NewController(therm, actuator, 1, 0, 0, logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		logger.Info("running....")
		if err := pid.Run(ctx, 37); err != nil {
			logger.Error(err)
		}
	}()

	<-ctx.Done()
	logger.Info("Bye!")
}
