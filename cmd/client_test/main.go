package main

import (
	"errors"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"

	"time"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/client"
	"github.com/benjaminbartels/zymurgauge/internal/ds18b20"
)

var logger *log.Logger

func main() {

	logger = log.New(os.Stderr, "", log.LstdFlags)

	// Setup graceful exit
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		os.Exit(1)
	}()

	// ToDo: Don't hardcode
	addr, err := url.Parse("http://192.168.0.10:3000")
	if err != nil {
		logger.Fatal(err)
	}
	c, err := client.NewClient(*addr, "v1", logger)
	if err != nil {
		logger.Fatal(err)
	}

	for {

		err := updateChamber(c.ChamberResource)
		if err != nil {
			logger.Fatal(err)
		}

		time.Sleep(5 * time.Second)
	}
}

func updateChamber(c *client.ChamberResource) error {

	logger.Println("Saving Chamber...")

	// mac, err := getMacAddress()
	// if err != nil {
	// 	panic(err)
	// }

	mac := "b8:27:eb:8e:d1:75"

	r := raspi.NewAdaptor()

	err := c.Save(&internal.Chamber{
		MacAddress: mac,
		Name:       "Chamber 1",
		Thermostat: &internal.Thermostat{
			Thermometer: &ds18b20.Thermometer{ID: "28-000006285484"},
			Chiller:     gpio.NewRelayDriver(r, "17"),
		},

		CurrentFermentation: &internal.Fermentation{
			ID: 1,
			Beer: internal.Beer{
				ID:    1,
				Name:  "My Stout",
				Style: "Stout",
				Schedule: []internal.FermentationStep{
					internal.FermentationStep{
						Order:      1,
						TargetTemp: 20.0,
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
