package thermostat

import (
	"testing"
	"time"

	logtest "github.com/sirupsen/logrus/hooks/test"
)

func TestOptions(t *testing.T) {
	l, _ := logtest.NewNullLogger()
	expected := time.Nanosecond
	clock := NewRealClock()

	thermostat := NewThermostat(nil, nil, nil, 1, 1, 1, 1, 1, 1, l,
		SetClock(clock),
		SetChillingCyclePeriod(expected),
		SetHeatingCyclePeriod(expected),
		SetChillingMinimum(expected),
		SetHeatingMinimum(expected))

	if thermostat.clock != clock {
		t.Errorf("Unexpected t.clock. Want: '%s', Got: '%s'", expected, thermostat.clock)
	}

	if thermostat.chillingCyclePeriod != expected {
		t.Errorf("Unexpected t.chillingCyclePeriod. Want: '%s', Got: '%s'", expected, thermostat.chillingCyclePeriod)
	}

	if thermostat.heatingCyclePeriod != expected {
		t.Errorf("Unexpected t.heatingCyclePeriod. Want: '%s', Got: '%s'", expected, thermostat.heatingCyclePeriod)
	}

	if thermostat.chillingMinimum != expected {
		t.Errorf("Unexpected t.chillingMinimum. Want: '%s', Got: '%s'", expected, thermostat.chillingMinimum)
	}

	if thermostat.heatingMinimum != expected {
		t.Errorf("Unexpected t.heatingMinimum. Want: '%s', Got: '%s'", expected, thermostat.heatingMinimum)
	}
}
