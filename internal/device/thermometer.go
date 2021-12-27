package device

// Hydrometer represents a device that can read temperatures.
type Thermometer interface {
	GetTemperature() (float64, error)
}

// type ThermometerRepo interface {
// 	GetThermometerIDs() ([]string, error)
// }
