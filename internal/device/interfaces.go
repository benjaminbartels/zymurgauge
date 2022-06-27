package device

import "context"

type Sensor interface {
	GetID() string
}

// Thermometer represents a device that can read temperatures.
type Thermometer interface {
	Sensor
	GetTemperature() (float64, error)
}

// Hydrometer represents a device that can read specific gravity.
type Hydrometer interface {
	Sensor
	GetGravity() (float64, error)
}

type ThermometerAndHydrometer interface {
	Sensor
	Thermometer
	Hydrometer
}

// Actuator represents a device that can be switched on and off.
type Actuator interface {
	On() error
	Off() error
	PWMOn(duty float64) error
}

type TemperatureController interface {
	Run(ctx context.Context, setPoint float64) error
}
