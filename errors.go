package zymurgauge

// General errors.
const (
	ErrInternal    = Error("internal error")
	ErrInvalidJSON = Error("invalid json")
)

// Errors
const (
	ErrChamberRequired = Error("Chamber required")

	ErrFermentationRequired = Error("Fermentation required")
	ErrBeerRequired         = Error("Beer required")
	ErrNotFound             = Error("Not found")
	ErrExists               = Error("Already exists")
	ErrIDRequired           = Error("ID required")
	ErrMacAddressRequired   = Error("Mac Address required")
)

// Error represents a zymurgauge error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
