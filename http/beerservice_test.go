package http_test

import (
	"errors"
	"reflect"
	"testing"

	"time"

	"github.com/orangesword/zymurgauge"
)

func TestBeerService_Get(t *testing.T) {
	t.Run("OK", testBeerService_Get)
	t.Run("NotFound", testBeerService_Get_NotFound)
	t.Run("ErrInternal", testBeerService_Get_ErrInternal)
}

func testBeerService_Get(t *testing.T) {
	s, c := MustOpenServerAndClient()
	defer func() { _ = s.Close() }()

	var beerID uint64 = 88
	var now = time.Now()

	b := &zymurgauge.Beer{ID: beerID, Name: "NAME", ModTime: now}

	s.BeerServiceMock.GetFn = func(id uint64) (*zymurgauge.Beer, error) {
		if id != beerID {
			t.Fatalf("unexpected id: %d", beerID)
		}
		return b, nil
	}

	r, err := c.BeerService().Get(beerID)

	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(b, r) {
		t.Fatalf("unexpected beer: %#v", r)
	}
}

func testBeerService_Get_NotFound(t *testing.T) {
	s, c := MustOpenServerAndClient()
	defer func() { _ = s.Close() }()

	s.BeerServiceMock.GetFn = func(id uint64) (*zymurgauge.Beer, error) {
		return nil, nil
	}

	if r, err := c.BeerService().Get(0); err != nil {
		t.Fatal(err)
	} else if r != nil {
		t.Fatal("expected nil beer")
	}
}

func testBeerService_Get_ErrInternal(t *testing.T) {
	s, c := MustOpenServerAndClient()
	defer func() { _ = s.Close() }()

	s.BeerServiceMock.GetFn = func(id uint64) (*zymurgauge.Beer, error) {
		return nil, errors.New("some internal error")
	}

	if _, err := c.BeerService().Get(0); err != zymurgauge.ErrInternal {
		t.Fatal(err)
	}

}

func TestBeerService_Save(t *testing.T) {
	t.Run("OK", testBeerService_Save)
	t.Run("ErrBeerRequired", testBeerService_Save_ErrBeerRequired)
	t.Run("ErrInternal", testBeerService_Save_ErrInternal)
}

func testBeerService_Save(t *testing.T) {
	s, c := MustOpenServerAndClient()
	defer func() { _ = s.Close() }()

	var now = time.Now()

	n := &zymurgauge.Beer{Name: "NAME"}

	s.BeerServiceMock.SaveFn = func(b *zymurgauge.Beer) error {

		b.ModTime = now
		b.ID = 1

		return nil
	}

	err := c.BeerService().Save(n)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(n, &zymurgauge.Beer{ID: 1, Name: "NAME", ModTime: now}) {
		t.Fatalf("unexpected beer: %#v", n)
	}
}

func testBeerService_Save_ErrBeerRequired(t *testing.T) {
	s, c := MustOpenServerAndClient()
	defer func() { _ = s.Close() }()

	if err := c.BeerService().Save(nil); err != zymurgauge.ErrBeerRequired {
		t.Fatal(err)
	}
}

func testBeerService_Save_ErrInternal(t *testing.T) {
	s, c := MustOpenServerAndClient()
	defer func() { _ = s.Close() }()

	s.BeerServiceMock.SaveFn = func(b *zymurgauge.Beer) error {
		return errors.New("some internal error")
	}

	if err := c.BeerService().Save(&zymurgauge.Beer{ID: 0}); err != zymurgauge.ErrInternal {
		t.Fatal(err)
	}
}
