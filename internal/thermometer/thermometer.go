package thermometer

// Thermometer represents a device that and read temperatures.
type Thermometer interface {
	Read() (float64, error)
}
