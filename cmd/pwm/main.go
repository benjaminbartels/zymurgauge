package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/benjaminbartels/zymurgauge/internal/device/gpio"
	"github.com/benjaminbartels/zymurgauge/internal/device/onewire"
	"github.com/benjaminbartels/zymurgauge/internal/temperaturecontrol/pid"
	"github.com/sirupsen/logrus"
	"periph.io/x/host/v3"
)

func main() {
	logger := logrus.New()

	if _, err := host.Init(); err != nil {
		panic(err)
	}

	therm, err := onewire.New("28-000006285484")
	if err != nil {
		panic(err)
	}

	actuator, err := gpio.NewGPIOActuator("19")
	if err != nil {
		panic(err)
	}

	pid := pid.NewController(therm, actuator, 1, 0, 0, logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := pid.Run(ctx, 100); err != nil {
			logger.Error(err)
		}
	}()

	<-ctx.Done()
	logger.Info("Bye!")
}
