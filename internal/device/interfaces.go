package device

import "context"

// Thermometer represents a device that can read temperatures.
type Thermometer interface {
	GetTemperature() (float64, error)
}

// Hydrometer represents a device that can read specific gravity.
type Hydrometer interface {
	GetGravity() (float64, error)
}

type ThermometerAndHydrometer interface {
	Thermometer
	Hydrometer
}

// Actuator represents a device that can be switched on and off.
type Actuator interface {
	On() error
	Off() error
}

type TemperatureController interface {
	Run(ctx context.Context, setPoint float64) error
}
