package tilt

const (
	ErrNotFound       = Error("tilt not found")
	ErrAlreadyRunning = Error("monitor is already running")
)

type Error string

func (e Error) Error() string {
	return string(e)
}
