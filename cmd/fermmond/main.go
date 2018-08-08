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

	"gobot.io/x/gobot/platforms/raspi"

	"github.com/alecthomas/kingpin"
	"github.com/benjaminbartels/zymurgauge/cmd/fermmond/client"
	"github.com/benjaminbartels/zymurgauge/internal"
)

// ToDo: Refactor
var (
	logger *log.Logger
	//basePath string
	chamber *internal.Chamber
	//status    *internal.ThermostatStatus
	zymClient *client.Client
	last      = 0.0
)

func main() {

	logger = log.New(os.Stderr, "", log.LstdFlags)

	// var err error
	// basePath, err = ioutil.TempDir("", "zymurgauge")
	// if err != nil {
	// 	logger.Fatal(err)
	// }
	// defer os.RemoveAll(basePath)

	// fmt.Println(basePath)

	// id := "28-000006285484"
	// mockData := "af 01 4b 46 7f ff 01 10 bc : crc=bc YES\naf 01 4b 46 7f ff 01 10 bc t=26937\n"

	// err = os.MkdirAll(filepath.Join(basePath, id), 0777)
	// if err != nil {
	// 	logger.Fatal(err)
	// }

	// if err = ioutil.WriteFile(filepath.Join(basePath, id, "w1_slave"), []byte(mockData), 0777); err != nil {
	// 	log.Fatal(err)
	// }

	// Setup graceful exit
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		os.Exit(1)
	}()

	address := kingpin.Flag("address", "Url of Zymurgauge server").Default("http://localhost:3000").Short('a').String()

	kingpin.Parse()

	var err error
	addr, err := url.Parse(*address)
	if err != nil {
		logger.Fatal(err)
	}

	zymClient, err = client.NewClient(addr, "v1", logger)
	if err != nil {
		logger.Fatal(err)
	}

	ch := make(chan internal.Chamber)

	mac, err := getMacAddress()
	if err != nil {
		logger.Fatal(err)
	}

	for {

		var cham *internal.Chamber
		cham, err = zymClient.ChamberResource.Get(mac)
		if err != nil {
			logger.Fatal(err)
		}

		if cham != nil {
			if err = processChamber(cham); err != nil {
				logger.Fatal(err)
			}
			break
		}

		logger.Printf("No Chamber found for Mac: %s, retrying in 5 seconds\n", mac)

		time.Sleep(5 * time.Second)
	}

	err = zymClient.ChamberResource.Subscribe(mac, ch)
	if err != nil {
		logger.Fatal(err)
	}

	for {
		logger.Println("Waiting for ChamberService updates")
		c := <-ch
		if err = processChamber(&c); err != nil {
			logger.Fatal(err)
		}
	}
}

func processChamber(c *internal.Chamber) error {
	logger.Println("Processing Chamber:", c.Name, c.MacAddress)

	if chamber == nil {
		chamber = c
	} else {
		c.Thermostat.Off()
	}

	chamber.Thermostat.SetLogger(logger)

	//ds18b20.DevicePath = basePath
	if err := c.Thermostat.InitThermometer(); err != nil {
		return err
	}

	chamber.Thermostat.InitActuators(raspi.NewAdaptor())

	//Check for updated fermentation
	if chamber.CurrentFermentationID != nil {

		ferm, err := zymClient.FermentationResource.Get(*chamber.CurrentFermentationID)
		if err != nil {
			return err
		}

		if ferm != nil {

			chamber.Thermostat.Subscribe("test_key", func(s internal.ThermostatStatus) {
				logger.Printf("State: %v, Error: %s\n", s.State, s.Error)

				if s.CurrentTemperature != nil && *s.CurrentTemperature != last {
					last = *s.CurrentTemperature
					change := &internal.TemperatureChange{
						FermentationID: ferm.ID,
						InsertTime:     time.Now(),
						Chamber:        chamber.Name,
						Beer:           ferm.Beer.Name,
						Thermometer:    chamber.Thermostat.ThermometerID,
						Temperature:    *s.CurrentTemperature,
					}
					if err := zymClient.FermentationResource.SaveTemperatureChange(change); err != nil {
						logger.Println(err)
					}
				}
			})

			chamber.Thermostat.Set(ferm.Beer.Schedule[0].TargetTemp)
			chamber.Thermostat.On()

		} else {
			logger.Println("Could not find Fermentation")
		}

	} else {
		logger.Println("No Current Fermentation")
	}

	return nil
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
