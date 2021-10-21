package thermometer

// Thermometer represents a device that and read temperatures.
type Thermometer interface {
	GetTemperature() (float64, error)
}