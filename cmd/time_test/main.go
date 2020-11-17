package main

import (
	"fmt"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/platform/temporal"
)

func main() {

	// start := time.Now()
	// c := clock.NewDilatedClock(2)
	// time.Sleep(1 * time.Second)
	// end := c.Now()
	// fmt.Println("Start:", start)
	// fmt.Println("End:", end)
	// fmt.Println("Diff:", end.Sub(start))
	// fmt.Println("Done.")

	start := time.Now()
	c := temporal.NewDilatedClock(600)

	<-c.After(10 * time.Minute)

	end := c.Now()
	fmt.Println("Start:", start)
	fmt.Println("End:", end)
	fmt.Println("Diff:", end.Sub(start))
	fmt.Println("Done.")
}
