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

	pid := NewTemperatureController(nil, nil, nil, 1, 1, 1, 1, 1, 1, l,
		SetClock(clock),
		SetChillingCyclePeriod(expected),
		SetHeatingCyclePeriod(expected),
		SetChillingMinimum(expected),
		SetHeatingMinimum(expected))

	if pid.clock != clock {
		assert.Equal(t, expected, pid.clock)
	}

	if pid.chillingCyclePeriod != expected {
		assert.Equal(t, expected, pid.chillingCyclePeriod)
	}

	if pid.heatingCyclePeriod != expected {
		assert.Equal(t, expected, pid.heatingCyclePeriod)
	}

	if pid.chillingMinimum != expected {
		assert.Equal(t, expected, pid.chillingMinimum)
	}

	if pid.heatingMinimum != expected {
		assert.Equal(t, expected, pid.heatingMinimum)
	}
}
