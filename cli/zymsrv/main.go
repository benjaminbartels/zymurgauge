package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/orangesword/zymurgauge"
	"github.com/orangesword/zymurgauge/bolt"
	"github.com/orangesword/zymurgauge/http"
	"github.com/sirupsen/logrus"
)

func main() {

	// Setup graceful exit
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		os.Exit(1)
	}()

	logger := logrus.New()
	logger.Level = logrus.DebugLevel // ToDo: set to InfoLevel
	//d.logger.Formatter = new(logrus.JSONFormatter)

	// if d.debug {
	// 	logger.Level = logrus.DebugLevel
	// }

	c := bolt.NewClient("zymurgaugedb", logger)
	err := c.Open()
	if err != nil {
		panic(err)
	}

	d := Daemon{
		server: *http.NewServer(logger),
		client: c,
	}

	d.Run(logger)

}

// Daemon is the container for the application:
type Daemon struct {
	server http.Server
	client zymurgauge.Client
}

// Run starts the Daemon
func (d Daemon) Run(logger *logrus.Logger) {

	d.server.Handler = &http.Handler{
		BeerHandler:         http.NewBeerHandler(logger),
		FermentationHandler: http.NewFermentationHandler(logger),
		ChamberHandler:      http.NewChamberHandler(logger),
	}
	d.server.Handler.BeerHandler.BeerService = d.client.BeerService()
	d.server.Handler.FermentationHandler.FermentationService = d.client.FermentationService()
	d.server.Handler.ChamberHandler.ChamberService = d.client.ChamberService()

	err := d.server.Open()
	if err != nil {
		logger.Panic(err)
	}

	logger.Infof("Listening on port %d", d.server.Port())

	d.server.Handler.ChamberHandler.Start() // ToDo: move go routine inside?

}
