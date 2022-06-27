package pid

import (
	"testing"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/platform/clock"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()
	expected := time.Nanosecond
	clock := clock.NewRealClock()

	pid := NewPIDTemperatureController(nil, nil, 1, 1, 1, l, SetClock(clock), Period(expected))

	if pid.clock != clock {
		assert.Equal(t, expected, pid.clock)
	}

	if pid.period != expected {
		assert.Equal(t, expected, pid.period)
	}
}
