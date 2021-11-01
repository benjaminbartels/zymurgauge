package device

// Thermometer represents a device that and read temperatures.
type Thermometer interface {
	GetTemperature() (float64, error)
}

type ThermometerRepo interface {
	GetThermometerIDs() ([]string, error)
}
