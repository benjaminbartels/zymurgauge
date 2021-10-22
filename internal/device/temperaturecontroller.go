package device

import "context"

type TemperatureController interface {
	Run(ctx context.Context, setPoint float64) error
}
