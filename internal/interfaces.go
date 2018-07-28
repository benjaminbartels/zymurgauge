package internal

type Thermometer interface {
	Read() (*float64, error)
}

type Actuator interface {
	On() error
	Off() error
}
