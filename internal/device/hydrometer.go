package device

// Hydrometer represents a device that can read gravity.
type Hydrometer interface {
	GetGravity() (float64, error)
}
