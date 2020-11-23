add simulatorpackage simulator_test

import (
	"testing"

	"github.com/benjaminbartels/zymurgauge/cmd/zymsim/simulator"
)

const initialBeerTemp = 25.0

func TestChill(t *testing.T) {
	expected := 24.556124476072753
	sim := simulator.New(initialBeerTemp)

	if err := sim.Chiller.On(); err != nil {
		t.Errorf("Unexpected error. Got: %+v", err)
	}

	for i := 0; i < 1000; i++ {
		sim.Update()
	}

	temp, err := sim.Thermometer.Read()
	if err != nil {
		t.Errorf("Unexpected error. Got: %+v", err)
	}

	if expected != temp {
		t.Errorf("Unexpected temp. Want: '%f', Got: '%f'", expected, temp)
	}
}

func TestHeat(t *testing.T) {
	expected := 25.771302060855625
	sim := simulator.New(initialBeerTemp)

	if err := sim.Heater.On(); err != nil {
		t.Errorf("Unexpected error. Got: %+v", err)
	}

	for i := 0; i < 1000; i++ {
		sim.Update()
	}

	temp, err := sim.Thermometer.Read()
	if err != nil {
		t.Errorf("Unexpected error. Got: %+v", err)
	}

	if expected != temp {
		t.Errorf("Unexpected temp. Want: '%f', Got: '%f'", expected, temp)
	}
}

func TestOnOf(t *testing.T) {
	expected := 23.964487886702734
	sim := simulator.New(initialBeerTemp)

	if err := sim.Chiller.On(); err != nil {
		t.Errorf("Unexpected error. Got: %+v", err)
	}

	for i := 0; i < 1000; i++ {
		sim.Update()
	}

	if err := sim.Chiller.Off(); err != nil {
		t.Errorf("Unexpected error. Got: %+v", err)
	}

	for i := 0; i < 1000; i++ {
		sim.Update()
	}

	temp, err := sim.Thermometer.Read()
	if err != nil {
		t.Errorf("Unexpected error. Got: %+v", err)
	}

	if expected != temp {
		t.Errorf("Unexpected temp. Want: '%f', Got: '%f'", expected, temp)
	}
}
