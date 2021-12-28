package device

// Hydrometer represents a device that can read specific gravity.
type Hydrometer interface {
	GetSpecificGravity() (float64, error)
}
