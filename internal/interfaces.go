package internal

// Thermometer represents a device that and read temperatures
type Thermometer interface {
	Read() (*float64, error)
}

// Actuator represents a device that can be switched on and off
type Actuator interface {
	On() error
	Off() error
}
