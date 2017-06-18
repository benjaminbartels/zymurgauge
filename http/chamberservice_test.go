package http_test

import (
	"errors"
	"reflect"
	"testing"

	"time"

	"github.com/benjaminbartels/zymurgauge"
	"github.com/benjaminbartels/zymurgauge/gpio"
)

func TestChamberService_Get(t *testing.T) {
	t.Run("OK", testChamberService_Get)
	t.Run("NotFound", testChamberService_Get_NotFound)
	t.Run("ErrInternal", testChamberService_Get_ErrInternal)
}

func testChamberService_Get(t *testing.T) {
	s, c := MustOpenServerAndClient()
	defer func() { _ = s.Close() }()

	m := "00:11:22:33:44:55"
	var now = time.Now()

	ctl := &zymurgauge.Chamber{
		MacAddress: m, Name: "NAME",
		ModTime: now,
		Controller: &gpio.Thermostat{
			ThermometerID: "blah",
		}}

	s.ChamberServiceMock.GetFn = func(mac string) (*zymurgauge.Chamber, error) {
		if m != mac {
			t.Fatalf("unexpected mac: %s - %s", m, mac)
		}
		return ctl, nil
	}

	r, err := c.ChamberService().Get(m)

	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(ctl, r) {
		t.Fatalf("unexpected controller: %#v", r)
	}
}

func testChamberService_Get_NotFound(t *testing.T) {
	s, c := MustOpenServerAndClient()
	defer func() { _ = s.Close() }()

	s.ChamberServiceMock.GetFn = func(mac string) (*zymurgauge.Chamber, error) {
		return nil, nil
	}

	if r, err := c.ChamberService().Get("00:11:22:33:44:55"); err != nil {
		t.Fatal(err)
	} else if r != nil {
		t.Fatal("expected nil controller")
	}
}

func testChamberService_Get_ErrInternal(t *testing.T) {
	s, c := MustOpenServerAndClient()
	defer func() { _ = s.Close() }()

	s.ChamberServiceMock.GetFn = func(mac string) (*zymurgauge.Chamber, error) {
		return nil, errors.New("some internal error")
	}

	if _, err := c.ChamberService().Get("00:11:22:33:44:55"); err != zymurgauge.ErrInternal {
		t.Fatal(err)
	}

}

func TestChamberService_Save(t *testing.T) {
	t.Run("OK", testChamberService_Save)
	t.Run("ErrChamberRequired", testChamberService_Save_ErrChamberRequired)
	t.Run("ErrInternal", testChamberService_Save_MacAddressRequired)
}

func testChamberService_Save(t *testing.T) {
	s, c := MustOpenServerAndClient()
	defer func() { _ = s.Close() }()

	var now = time.Now()
	mac := "00:11:22:33:44:55"

	n := &zymurgauge.Chamber{
		Name:       "NAME",
		MacAddress: mac,
		ModTime:    now,
		Controller: &gpio.Thermostat{
			ThermometerID: "blah",
		},
		CurrentFermentation: &zymurgauge.Fermentation{
			ModTime: now,
		},
	}

	after := &zymurgauge.Chamber{}

	s.ChamberServiceMock.SaveFn = func(c *zymurgauge.Chamber) error {

		c.ModTime = now
		c.MacAddress = mac
		after = c

		return nil
	}

	err := c.ChamberService().Save(n)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(n, after) {
		t.Fatalf("unexpected controller: %#v", n)
	}
}

func testChamberService_Save_ErrChamberRequired(t *testing.T) {
	s, c := MustOpenServerAndClient()
	defer func() { _ = s.Close() }()

	if err := c.ChamberService().Save(nil); err != zymurgauge.ErrChamberRequired {
		t.Fatal(err)
	}
}

func testChamberService_Save_MacAddressRequired(t *testing.T) {
	s, c := MustOpenServerAndClient()
	defer func() { _ = s.Close() }()

	s.ChamberServiceMock.SaveFn = func(b *zymurgauge.Chamber) error {
		return errors.New("some internal error")
	}

	if err := c.ChamberService().Save(&zymurgauge.Chamber{}); err !=
		zymurgauge.ErrMacAddressRequired {
		t.Fatal(err)
	}
}
