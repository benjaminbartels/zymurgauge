package main

import (
	"errors"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal"

	"github.com/benjaminbartels/zymurgauge/internal/client"
	"github.com/felixge/pidctrl"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	APIAddress string `required:"true"`
}

func main() {

	logger := log.New(os.Stderr, "", log.LstdFlags)

	// Process env variables
	var cfg config
	err := envconfig.Process("fermmon", &cfg)
	if err != nil {
		logger.Fatal(err.Error())
	}

	// Parse server url
	addr, err := url.Parse(cfg.APIAddress)
	if err != nil {
		logger.Fatal(err.Error())
	}

	// Create a new client
	client, err := client.NewClient(addr, "v1", logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	// Set mac address
	macAddress, err := getMacAddress()
	if err != nil {
		logger.Fatal(err.Error())
	}

	// Create PID Controller
	pid := pidctrl.NewPIDController(1, 1, 0) // ToDo: get from env vars

	var ctl *chamberCtl

	for {
		var err error
		chamber, err := client.ChamberResource.Get(macAddress)
		if err != nil {
			logger.Fatal(err)
		}

		if chamber != nil {

			ctl, err = newChamberCtl(chamber, pid, client, logger,
				internal.MinimumChill(1*time.Second), // ToDO: env vars
				internal.MinimumHeat(1*time.Second),
				internal.Interval(1*time.Second),
				internal.Logger(logger))
			if err != nil {
				logger.Fatal(err)
			}
			break

		}

		logger.Printf("No Chamber found for Mac: %s, retrying in 5 seconds\n", macAddress)

		time.Sleep(5 * time.Second)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	// Start listening
	ctl.Listen()

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	// Stop listening
	ctl.Close()

	wg.Wait()
	logger.Println("Bye!")
}

// getMacAddress returns the first Mac Address of the first network interface found
func getMacAddress() (string, error) {

	mac := ""

	interfaces, err := net.Interfaces()
	if err != nil {
		return "", errors.New("failed to get host MAC address")
	}
	for _, iface := range interfaces {
		if len(iface.HardwareAddr.String()) > 0 {
			if iface.Name == "wlan0" { //ToDo: fix this
				mac = iface.HardwareAddr.String()
			}
			if iface.Name == "en0" && mac == "" { //ToDo: fix this
				mac = iface.HardwareAddr.String()
			}
		}
	}

	if mac == "" {
		return mac, errors.New("Failed to get host MAC address. No valid interfaces found")
	}

	return mac, nil

}
