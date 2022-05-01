package bluetooth

const (
	ErrInvalidManufacturerDataLength = Error("manufacurerData length is less that 25")
)

type Error string

func (e Error) Error() string {
	return string(e)
}
