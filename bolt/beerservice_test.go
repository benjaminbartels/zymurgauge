package bolt_test

import (
	"reflect"
	"testing"

	"github.com/orangesword/zymurgauge"
)

func TestBeerService_Save_New(t *testing.T) {

	c := MustOpenClient()
	defer func() { _ = c.Close() }()

	s := c.BeerService()

	b := zymurgauge.Beer{
		Name: "My Beer",
	}

	if err := s.Save(&b); err != nil {
		t.Fatal(err)
	} else if b.ID != 1 {
		t.Fatalf("unexpected id: %d", b.ID)
	}

	other, err := s.Get(1)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&b, other) {
		t.Fatalf("unexpected beer: %#v", other)
	}
}

func TestBeerService_Save_Existing(t *testing.T) {

	c := MustOpenClient()
	defer func() { _ = c.Close() }()

	s := c.BeerService()

	b1 := &zymurgauge.Beer{Name: "My Beer 1"}
	b2 := &zymurgauge.Beer{Name: "My Beer 2"}

	if err := s.Save(b1); err != nil {
		t.Fatal(err)
	} else if err := s.Save(b2); err != nil {
		t.Fatal(err)
	}

	b1.Style = "Lager"
	b2.Style = "Porter"

	if err := s.Save(b1); err != nil {
		t.Fatal(err)
	} else if err := s.Save(b2); err != nil {
		t.Fatal(err)
	}

	if ub1, err := s.Get(b1.ID); err != nil {
		t.Fatal(err)
	} else if ub1.Style != b1.Style {
		t.Fatalf("unexpected beer #1 style: %f", ub1.Style)
	}

	if ub2, err := s.Get(b2.ID); err != nil {
		t.Fatal(err)
	} else if ub2.Style != b2.Style {
		t.Fatalf("unexpected beer #2 style: %f", ub2.Style)
	}
}
