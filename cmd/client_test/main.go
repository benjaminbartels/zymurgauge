package main

import (
	"errors"
	"net"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"time"

	"fmt"

	"github.com/benjaminbartels/zymurgauge"
	"github.com/benjaminbartels/zymurgauge/http"
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
	logger.Level = logrus.DebugLevel

	addr, err := url.Parse("http://localhost:3000")
	if err != nil {
		panic(err)
	}
	c := http.NewClient(*addr, logger)

	for {

		err := updateChamber(c.ChamberService())
		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(5 * time.Second)
	}

}

func updateChamber(c zymurgauge.ChamberService) error {

	fmt.Println("Saving Chamber...")

	mac, err := getMacAddress()
	if err != nil {
		panic(err)
	}

	coolerGPIO := 17
	heaterGPIO := 0

	err = c.Save(&zymurgauge.Chamber{
		MacAddress: mac,
		Name:       "Chamber 1",
		Controller: &zymurgauge.TemperatureController{
			ThermometerID: "28-000006285484",
			CoolerGPIO:    &coolerGPIO,
			HeaterGPIO:    &heaterGPIO,
			Interval:      5 * time.Second,
		},
		CurrentFermentation: &zymurgauge.Fermentation{
			ID: 1,
			Beer: zymurgauge.Beer{
				ID:    1,
				Name:  "My Stout",
				Style: "Stout",
				Schedule: []zymurgauge.FermentationStep{
					zymurgauge.FermentationStep{
						Order:      1,
						TargetTemp: 25.0,
						Duration:   9999999,
					},
				},
			},
			CurrentStep: 1,
		},
	})

	return err

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
