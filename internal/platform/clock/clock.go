package clock

import "time"

type Clock interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
}
