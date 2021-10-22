package device

// Actuator represents a device that can be switched on and off.
type Actuator interface {
	On() error
	Off() error
}
