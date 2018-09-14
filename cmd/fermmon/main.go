package main

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"

	"github.com/benjaminbartels/zymurgauge/cmd/fermmon/controller"
	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/client"
	"github.com/felixge/pidctrl"
	"github.com/kelseyhightower/envconfig"
)

const clientID = "18CHgJa2D3GyxmZfKdF2uhmSv4aS78Xb"

type config struct {
	APIAddress   string `required:"true"`
	ClientSecret string `required:"true"`
	Interface    string
	Interval     time.Duration `default:"10m"`
	MinimumChill time.Duration `default:"3m"`
	MinimumHeat  time.Duration `default:"3m"`
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

	// Create PID Controller
	pid := pidctrl.NewPIDController(1, 1, 0) // ToDo: get from env vars

	ctl := controller.New(mac, pid, client.ChamberProvider, client.FermentationProvider, logger,
		internal.MinimumChill(cfg.MinimumChill),
		internal.MinimumHeat(cfg.MinimumHeat),
		internal.Interval(cfg.Interval),
		internal.Logger(logger))

	var wg sync.WaitGroup
	wg.Add(1)

	// Start polling
	ctl.Start(10 * time.Second)

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	// Stop polling
	ctl.Stop()

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
			if iface.Name == "eth0" { //ToDo: fix this
				mac = iface.HardwareAddr.String()
			}
			if iface.Name == "wlan0" { //ToDo: fix this
				mac = iface.HardwareAddr.String()
			}
			if mac == "" && iface.Name == "en0" { //ToDo: fix this
				mac = iface.HardwareAddr.String()
			}
		}
	}

	if mac == "" {
		return mac, errors.New("Failed to get host MAC address. No valid interfaces found")
	}

	return mac, nil

}

// getMacAddressByIFaceName returns the mac address of the given interface name
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
