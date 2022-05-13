package ibeacon

import "context"

type Discoverer interface {
	Discover(ctx context.Context) (chan Event, error)
}
