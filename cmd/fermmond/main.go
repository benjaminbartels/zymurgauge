package main

import (
	"errors"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/client"
)

// ToDo: Remove global vars
var (
	currentChamber *internal.Chamber
	logger         *log.Logger
)

func main() {

	logger = log.New(os.Stderr, "", log.LstdFlags)

	// Setup graceful exit
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		os.Exit(1)
	}()

	address := kingpin.Flag("address", "Url of Zymurgauge server").Default("http://localhost:3000").Short('u').String()

	kingpin.Parse()

	addr, err := url.Parse(*address)
	if err != nil {
		logger.Fatal(err)
	}

	c, err := client.NewClient(*addr, "v1", logger)
	if err != nil {
		logger.Fatal(err)
	}

	ch := make(chan internal.Chamber)

	mac, err := getMacAddress()
	if err != nil {
		logger.Fatal(err)
	}

	for {

		var chamber *internal.Chamber
		chamber, err = c.ChamberResource.Get(mac)
		if err != nil {
			logger.Fatal(err)
		}

		if chamber != nil {
			processChamber(chamber)
			break
		}

		logger.Printf("No Chamber found for Mac: %s, retrying in 5 seconds\n", mac)

		time.Sleep(5 * time.Second)
	}

	err = c.ChamberResource.Subscribe(mac, ch)
	if err != nil {
		logger.Fatal(err)
	}

	for {
		logger.Println("Waiting for ChamberService updates")
		c := <-ch
		processChamber(&c)
	}

}

func processChamber(c *internal.Chamber) {
	logger.Println("processChamber called")

	currentChamber = c

	logger.Println(currentChamber)

	//Check for updated fermentation
	if currentChamber.CurrentFermentation != nil {
		currentChamber.Thermostat.Set(c.CurrentFermentation.Beer.Schedule[0].TargetTemp)
	}
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
