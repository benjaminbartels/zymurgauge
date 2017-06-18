package bolt_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/benjaminbartels/zymurgauge"
	"github.com/benjaminbartels/zymurgauge/gpio"
)

func TestChamberService_Save_New(t *testing.T) {

	mac := "00:11:22:33:44:55"
	now := time.Now()

	cl := MustOpenClient()
	defer func() { _ = cl.Close() }()

	s := cl.ChamberService()

	c := zymurgauge.Chamber{
		Name:       "My Chamber",
		MacAddress: mac,
		Controller: &gpio.Thermostat{
			ThermometerID: "blah",
		},
		CurrentFermentation: &zymurgauge.Fermentation{
			Beer: zymurgauge.Beer{
				Name:    "My Beer",
				ModTime: now,
			},
			StartTime:     now,
			CompletedTime: now,
			ModTime:       now,
		},
	}

	if err := s.Save(&c); err != nil {
		t.Fatal(err)
	} else if c.MacAddress != mac {
		t.Fatalf("unexpected mac: %s", c.MacAddress)
	}

	other, err := s.Get(mac)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&c, other) {
		t.Fatalf("unexpected controller: %#v", other)
	}
}

func TestChamberService_Save_Existing(t *testing.T) {

	cl := MustOpenClient()
	defer func() { _ = cl.Close() }()

	s := cl.ChamberService()

	c1 := &zymurgauge.Chamber{Name: "My Chamber 1", MacAddress: "00:11:22:33:44:55"}
	c2 := &zymurgauge.Chamber{Name: "My Chamber 2", MacAddress: "aa:bb:cc:dd:ee:ff"}

	if err := s.Save(c1); err != nil {
		t.Fatal(err)
	} else if err := s.Save(c2); err != nil {
		t.Fatal(err)
	}

	if err := s.Save(c1); err != nil {
		t.Fatal(err)
	} else if err := s.Save(c2); err != nil {
		t.Fatal(err)
	}

	if uc1, err := s.Get(c1.MacAddress); err != nil {
		t.Fatal(err)
	} else if uc1.MacAddress != c1.MacAddress {
		t.Fatalf("unexpected controller #1 MacAddress: %f", uc1.MacAddress)
	}

	if uc2, err := s.Get(c2.MacAddress); err != nil {
		t.Fatal(err)
	} else if uc2.MacAddress != c2.MacAddress {
		t.Fatalf("unexpected controller #2 MacAddress: %f", uc2.MacAddress)
	}
}

func TestChamberService_Subscribe_And_Receive(t *testing.T) {

	cl := MustOpenClient()
	defer func() { _ = cl.Close() }()

	s := cl.ChamberService()

	c := &zymurgauge.Chamber{Name: "My Chamber 1", MacAddress: "00:11:22:33:44:55"}

	if err := s.Save(c); err != nil {
		t.Fatal(err)
	}

	ch := make(chan zymurgauge.Chamber)
	done := make(chan bool)

	if err := s.Subscribe(c.MacAddress, ch); err != nil {
		t.Fatal(err)
	}

	var r zymurgauge.Chamber

	go func() {
		r = <-ch
		done <- true
	}()

	c.Name = "My Chamber 1 Updated"
	if err := s.Save(c); err != nil {
		t.Fatal(err)
	}

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		panic("timeout")
	}

	if r.Name != c.Name {
		t.Fatalf("unexpected controller: %#v", r)
	}
}
