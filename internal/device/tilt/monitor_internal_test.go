package tilt

import (
	"testing"
	"time"

	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()
	expected := time.Nanosecond

	monitor := NewMonitor(nil, l,
		SetTimeout(expected),
		SetInterval(expected))

	if monitor.timeout != expected {
		assert.Equal(t, expected, monitor.timeout)
	}

	if monitor.interval != expected {
		assert.Equal(t, expected, monitor.timeout)
	}
}
