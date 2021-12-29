package device

type ThermometerType string

// Hydrometer represents a device that can read temperatures.
type Thermometer interface {
	GetTemperature() (float64, error)
}
