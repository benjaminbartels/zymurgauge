package tilt

const (
	ErrNotFound                      = Error("chamber not found")
	ErrAlreadyRunning                = Error("monitor is already running")
	ErrInvalidManufacturerDataLength = Error("manufacurerData length is less that 25")
)

type Error string

func (e Error) Error() string {
	return string(e)
}
