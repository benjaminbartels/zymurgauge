package chamber

const (
	ErrNotFound       = Error("chamber not found")
	ErrNoCurrentBatch = Error("chamber does not have a current batch")
	ErrInvalidStep    = Error("invalid step")
	ErrNotFermenting  = Error("fermentation has not started")
	ErrFermenting     = Error("fermentation has started")
)

type Error string

func (e Error) Error() string {
	return string(e)
}

type InvalidConfigurationError struct {
	configErrors []error
}

func (e InvalidConfigurationError) Error() string {
	return "configuration is invalid"
}

func (e InvalidConfigurationError) Problems() []error {
	return e.configErrors
}
