//nolint:gomnd
package main

import (
	"net"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/benjaminbartels/zymurgauge/cmd/fermmon/controller"
	"github.com/benjaminbartels/zymurgauge/internal/client"
	"github.com/benjaminbartels/zymurgauge/internal/thermostat"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const clientID = "18CHgJa2D3GyxmZfKdF2uhmSv4aS78Xb"

type config struct {
	APIAddress          string `required:"true"`
	ClientSecret        string `required:"true"`
	Interface           string
	ChillingMinimum     time.Duration `default:"10m"`
	HeatingMinimum      time.Duration `default:"10s"`
	ChillingCyclePeriod time.Duration `default:"30m"`
	HeatingCyclePeriod  time.Duration `default:"10m"`
	ChillerP            float64       `default:"-1"`
	ChillerI            float64       `default:"0"`
	ChillerD            float64       `default:"0"`
	HeaterP             float64       `default:"1"`
	HeaterI             float64       `default:"0"`
	HeaterD             float64       `default:"0"`
}

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Process env variables
	var cfg config

	if err := envconfig.Process("fermmon", &cfg); err != nil {
		logger.Fatal(err.Error())
	}

	// Parse server url
	addr, err := url.Parse(cfg.APIAddress)
	if err != nil {
		logger.Fatal(err.Error())
	}

	client, err := client.NewClient(addr, "v1", clientID, cfg.ClientSecret, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	var mac string

	// Get mac address
	if cfg.Interface != "" {
		mac, err = getMacAddressByInterfaceName(cfg.Interface)
		if err != nil {
			logger.Fatal(err.Error())
		}
	} else {
		mac, err = getMacAddress()
		if err != nil {
			logger.Fatal(err.Error())
		}
	}

	ctl := controller.New(mac,
		cfg.ChillerP, cfg.ChillerI, cfg.ChillerD, cfg.HeaterP, cfg.HeaterI, cfg.HeaterD,
		client.ChamberProvider, client.FermentationProvider, logger,
		thermostat.SetChillingMinimum(cfg.ChillingMinimum),
		thermostat.SetHeatingMinimum(cfg.HeatingMinimum),
		thermostat.SetChillingCyclePeriod(cfg.ChillingCyclePeriod),
		thermostat.SetHeatingCyclePeriod(cfg.HeatingCyclePeriod),
	)

	var wg sync.WaitGroup

	wg.Add(1)

	ctl.Start(10 * time.Second)

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	ctl.Stop()

	wg.Wait()
	logger.Println("Bye!")
}

// getMacAddress returns the first Mac Address of the first network interface found.
func getMacAddress() (string, error) {
	mac := ""

	interfaces, err := net.Interfaces()
	if err != nil {
		return "", errors.New("failed to get host MAC address")
	}

	for _, iface := range interfaces {
		if len(iface.HardwareAddr.String()) > 0 {
			if iface.Name == "eth0" { // TODO: fix this
				mac = iface.HardwareAddr.String()
			}

			if iface.Name == "wlan0" { // TODO: fix this
				mac = iface.HardwareAddr.String()
			}

			if mac == "" && iface.Name == "en0" { // TODO: fix this
				mac = iface.HardwareAddr.String()
			}
		}
	}

	if mac == "" {
		return mac, errors.New("Failed to get host MAC address. No valid interfaces found")
	}

	return mac, nil
}

// getMacAddressByIFaceName returns the mac address of the given interface name.
func getMacAddressByInterfaceName(name string) (string, error) {
	var mac string

	interfaces, err := net.Interfaces()
	if err != nil {
		return "", errors.New("Failed to get host MAC address")
	}

	for _, iface := range interfaces {
		if iface.Name == name {
			mac = iface.HardwareAddr.String()
		}
	}

	if mac == "" {
		return mac, errors.Errorf("Failed to get host MAC address for interface %s.", name)
	}

	return mac, nil
}
