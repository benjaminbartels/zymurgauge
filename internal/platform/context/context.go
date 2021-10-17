package context

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func WithInterrupt(parent context.Context) (context.Context, func()) {
	return WithSignal(parent, syscall.SIGINT, syscall.SIGTERM)
}

func WithSignal(parent context.Context, signals ...os.Signal) (context.Context, func()) {
	ctx, closer := context.WithCancel(parent)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)

	go func() {
		select {
		case <-ch:
			closer()
		case <-ctx.Done():
		}
	}()

	return ctx, closer
}
