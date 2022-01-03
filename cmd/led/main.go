package main

import (
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/device/gpio"
	"periph.io/x/host/v3"
)

func main() {
	_, err := host.Init()
	if err != nil {
		panic(err)
	}

	red, err := gpio.NewGPIOActuator("17")
	if err != nil {
		panic(err)
	}

	green, err := gpio.NewGPIOActuator("22")
	if err != nil {
		panic(err)
	}

	for {
		if err = red.On(); err != nil {
			panic(err)
		}

		if err = green.Off(); err != nil {
			panic(err)
		}

		<-time.After(50 * time.Millisecond)

		if err = red.Off(); err != nil {
			panic(err)
		}

		if err = green.On(); err != nil {
			panic(err)
		}

		<-time.After(50 * time.Millisecond)
	}
}
