package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kong"
)

type cli struct {
	globals

	Once simulation    `kong:"cmd,help='Run a single simulation.'"`
	Many simulationSet `kong:"cmd,help='Run a set of simulations.'"`
}

type globals struct {
	Multiplier float64       `kong:"default=6000.0,short=m,help='Time dilation multiplier. Defaults to 6000.'"`
	Runtime    time.Duration `kong:"default=5s,short=r,help='Runtime of simulation. Defaults to 5s.'"`
	Debug      bool          `kong:"default=false,short=d,help='Enable debug logging. Default is false.'"`
}

func main() {
	cli := cli{}
	ctx := kong.Parse(&cli,
		kong.Name("zymsim"),
		kong.Description("Zymurgauge Thermostat Simulator"),
		kong.UsageOnError(),
	)

	if err := ctx.Run(&cli.globals); err != nil {
		fmt.Println("Failed to run command:", err)
		os.Exit(1)
	}
}
