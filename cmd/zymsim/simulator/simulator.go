package simulator

const (
	// Borrowed from https://github.com/BrewPi/firmware/blob/0.5.10/lib/test/SimulationTest.cpp#L115
	beerCapacity = 4.2 * 1.0 * 20        // heat capacity water * density of water * 20L volume (in kJ per kelvin).
	airCapacity  = 1.005 * 1.225 * 0.200 // heat capacity of dry air * density of air * 200L volume (in kJ per kelvin).
	// Moist air has only slightly higher heat capacity, 1.02 when saturated at 20C.
	wallCapacity            = 5.0 // just a guess
	heaterCapacity          = 1.0 // also a guess, to simulate that heater first heats itself, then starts heating the air
	heaterPower             = 0.1 // 100W, in kW.
	coolerPower             = 0.1 // 100W, in kW. Assuming 200W at 50% efficiency
	airBeerTransfer         = 1.0 / 300
	wallAirTransfer         = 1.0 / 300
	heaterAirTransfer       = 1.0 / 30
	environmentWallTransfer = 0.001 // losses to environment
	// heaterToBeer            = 0.0   // ratio of heater transferred directly to beer instead of fridge air
	// heaterToAir             = 1.0 - heaterToBeer.
	initialWallTemp        = 20.0
	initialAirTemp         = 20.0
	initialHeaterTemp      = 20.0
	initialEnvironmentTemp = 20.0
)

type Simulator struct {
	wallTemp        float64
	airTemp         float64
	beerTemp        float64
	heaterTemp      float64
	environmentTemp float64
	Thermometer     *Thermometer
	Chiller         *Actuator
	Heater          *Actuator
}

type Actuator struct {
	isOn bool
}

func (a *Actuator) On() error {
	if !a.isOn {
		a.isOn = true
	}

	return nil
}

func (a *Actuator) Off() error {
	if a.isOn {
		a.isOn = false
	}

	return nil
}

type Thermometer struct {
	currentTemp float64
}

func (t *Thermometer) Read() (float64, error) {
	return t.currentTemp, nil
}

func New(initialBeerTemp float64) *Simulator {
	return &Simulator{
		beerTemp:        initialBeerTemp,
		wallTemp:        initialWallTemp,
		airTemp:         initialAirTemp,
		heaterTemp:      initialHeaterTemp,
		environmentTemp: initialEnvironmentTemp,
		Thermometer:     &Thermometer{currentTemp: initialBeerTemp},
		Chiller:         &Actuator{},
		Heater:          &Actuator{},
	}
}

func (s *Simulator) Update() {
	beerTempNew := s.beerTemp
	airTempNew := s.airTemp
	wallTempNew := s.wallTemp
	heaterTempNew := s.heaterTemp

	beerTempNew += (s.airTemp - s.beerTemp) * airBeerTransfer / beerCapacity

	if s.Chiller.isOn {
		wallTempNew -= coolerPower / wallCapacity
	} else if s.Heater.isOn {
		heaterTempNew += heaterPower / heaterCapacity
	}

	airTempNew += (s.heaterTemp - s.airTemp) * heaterAirTransfer / airCapacity
	airTempNew += (s.wallTemp - s.airTemp) * wallAirTransfer / airCapacity
	airTempNew += (s.beerTemp - s.airTemp) * airBeerTransfer / airCapacity

	beerTempNew += (s.airTemp - s.beerTemp) * airBeerTransfer / beerCapacity

	heaterTempNew += (s.airTemp - s.heaterTemp) * heaterAirTransfer / heaterCapacity

	wallTempNew += (s.environmentTemp - s.wallTemp) * environmentWallTransfer / wallCapacity
	wallTempNew += (s.airTemp - s.wallTemp) * wallAirTransfer / wallCapacity

	s.airTemp = airTempNew
	s.beerTemp = beerTempNew
	s.wallTemp = wallTempNew
	s.heaterTemp = heaterTempNew
	s.Thermometer.currentTemp = s.beerTemp
}
