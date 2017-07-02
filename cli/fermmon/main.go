package main

import (
	"errors"
	"net"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/benjaminbartels/zymurgauge"
	"github.com/benjaminbartels/zymurgauge/gpio"
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
	logger.Level = logrus.DebugLevel // ToDo: set to InfoLevel
	//d.logger.Formatter = new(logrus.JSONFormatter)

	// if d.debug {
	// 	logger.Level = logrus.DebugLevel
	// }

	address := kingpin.Flag("address", "Url of Zymurgauge server").Default("http://192.168.0.12:3000").Short('u').String()

	kingpin.Parse()

	addr, err := url.Parse(*address)
	if err != nil {
		panic(err)
	}

	c := http.NewClient(*addr, logger)

	// ToDo: Try function parameters
	m := Daemon{
		client: c,
		logger: logger,
	}

	m.Run()

}

type Daemon struct {
	client  zymurgauge.Client
	chamber *zymurgauge.Chamber
	logger  *logrus.Logger
}

func (d Daemon) Run() {

	ch := make(chan zymurgauge.Chamber)

	mac, err := getMacAddress()
	if err != nil {
		panic(err)
	}

	for {
		d.chamber, err = d.client.ChamberService().Get(mac)
		if err != nil {
			panic(err)
		}

		if d.chamber != nil {
			d.processChamber(d.chamber)
			break
		}

		d.logger.Infof("No Chamber found for Mac: %s, retrying in 5 seconds", mac)

		time.Sleep(5 * time.Second)
	}

	err = d.client.ChamberService().Subscribe(mac, ch)
	if err != nil {
		panic(err)
	}

	for {
		d.logger.Debug("Waiting for ChamberService updates")
		c := <-ch
		d.processChamber(&c)
	}

}

func (d Daemon) processChamber(c *zymurgauge.Chamber) {

	if new, ok := c.Controller.(*gpio.Thermostat); ok {

		cur := d.chamber.Controller.(*gpio.Thermostat)
		cur.Logger = d.logger // Todo: maybe

		if cur.ThermometerID != new.ThermometerID ||
			cur.CoolerGPIO != new.CoolerGPIO ||
			cur.HeaterGPIO != new.HeaterGPIO ||
			cur.Interval != new.Interval {

			err := d.chamber.Controller.SetTemperature(nil)
			if err != nil {
				d.logger.Error(err)
			}
			d.chamber.Controller = c.Controller
		}
	}

	d.chamber.CurrentFermentation = c.CurrentFermentation

	// Check for updated fermentation
	if d.chamber.CurrentFermentation != nil {
		t := c.CurrentFermentation.Beer.Schedule[0].TargetTemp
		err := d.chamber.Controller.SetTemperature(&t)
		if err != nil {
			panic(err)
		}
	} else {
		err := d.chamber.Controller.SetTemperature(nil)
		if err != nil {
			panic(err)
		}
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
