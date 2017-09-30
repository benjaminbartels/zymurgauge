package ds18b20_test

import (
	"fmt"
	"testing"

	"github.com/benjaminbartels/zymurgauge/internal/ds18b20"
	"github.com/spf13/afero"
)

func TestGetThermometer(t *testing.T) {
	fs := afero.NewMemMapFs()
	id := "28-000123456789"
	mockData := "af 01 4b 46 7f ff 01 10 bc : crc=bc YES\naf 01 4b 46 7f ff 01 10 bc t=%d\n"
	val := 26937

	file, err := fs.Create(fmt.Sprintf("/sys/bus/w1/devices/%s/w1_slave", id))
	if err != nil {
		panic(err)
	}

	_, err = file.Write([]byte(fmt.Sprintf(mockData, val)))
	if err != nil {
		panic(err)
	}

	err = file.Close()
	if err != nil {
		panic(err)
	}

	grp := ds18b20.New(fs)

	therm, err := grp.GetThermometer(id)
	if err != nil {
		panic(err)
	}

	temp, err := therm.ReadTemperature()
	if err != nil {
		panic(err)
	}

	if *temp != 26.937 {
		t.Fatalf("unexpected celsius reading: %f", *temp)
	}

}
