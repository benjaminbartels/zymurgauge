package internal

// TemperatureController device that controls temperature //ToDo: Pending removal?
type TemperatureController interface {
	// On turns the Thermostat on and allows to being monitoring
	On()
	// Off turns the Thermostat Off
	Off()
	// Set sets TemperatureController to the specified temperature
	Set(temp float64)
}

// Thermometer represents a device that and read temperatures
type Thermometer interface {
	Read() (*float64, error)
}

// Actuator represents a device that can be switched on and off
type Actuator interface {
	On() error
	Off() error
}
