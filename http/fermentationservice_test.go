package http_test

import (
	"errors"
	"reflect"
	"testing"

	"time"

	"github.com/orangesword/zymurgauge"
)

func TestFermentationService_Get(t *testing.T) {
	t.Run("OK", testFermentationService_Get)
	t.Run("NotFound", testFermentationService_Get_NotFound)
	t.Run("ErrInternal", testFermentationService_Get_ErrInternal)
}

func testFermentationService_Get(t *testing.T) {
	s, c := MustOpenServerAndClient()
	defer func() { _ = s.Close() }()

	var fermentationID uint64 = 88
	var now = time.Now()

	f := &zymurgauge.Fermentation{
		Beer: zymurgauge.Beer{
			Name:    "My Beer",
			ModTime: now,
		},
		ID:            fermentationID,
		CurrentStep:   1,
		StartTime:     now,
		CompletedTime: now,
		ModTime:       now}

	s.FermentationServiceMock.GetFn = func(id uint64) (*zymurgauge.Fermentation, error) {
		if id != fermentationID {
			t.Fatalf("unexpected id: %d", fermentationID)
		}
		return f, nil
	}

	r, err := c.FermentationService().Get(fermentationID)

	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(f, r) {
		t.Fatalf("unexpected fermentation: %#v", r)
	}
}

func testFermentationService_Get_NotFound(t *testing.T) {
	s, c := MustOpenServerAndClient()
	defer func() { _ = s.Close() }()

	s.FermentationServiceMock.GetFn = func(id uint64) (*zymurgauge.Fermentation, error) {
		return nil, nil
	}

	if r, err := c.FermentationService().Get(0); err != nil {
		t.Fatal(err)
	} else if r != nil {
		t.Fatal("expected nil fermentation")
	}
}

func testFermentationService_Get_ErrInternal(t *testing.T) {
	s, c := MustOpenServerAndClient()
	defer func() { _ = s.Close() }()

	s.FermentationServiceMock.GetFn = func(id uint64) (*zymurgauge.Fermentation, error) {
		return nil, errors.New("some internal error")
	}

	if _, err := c.FermentationService().Get(0); err != zymurgauge.ErrInternal {
		t.Fatal(err)
	}

}

func TestFermentationService_Save(t *testing.T) {
	t.Run("OK", testFermentationService_Save)
	t.Run("ErrFermentationRequired", testFermentationService_Save_ErrFermentationRequired)
	t.Run("ErrInternal", testFermentationService_Save_ErrInternal)
}

func testFermentationService_Save(t *testing.T) {
	s, c := MustOpenServerAndClient()
	defer func() { _ = s.Close() }()

	var now = time.Now()

	n := &zymurgauge.Fermentation{
		Beer: zymurgauge.Beer{
			Name:    "My Beer",
			ModTime: now,
		},
		CurrentStep:   1,
		StartTime:     now,
		CompletedTime: now}

	s.FermentationServiceMock.SaveFn = func(f *zymurgauge.Fermentation) error {

		f.ModTime = now
		f.ID = 1

		return nil
	}

	exp := &zymurgauge.Fermentation{
		Beer: zymurgauge.Beer{
			Name:    "My Beer",
			ModTime: now,
		},
		ID:            1,
		ModTime:       now,
		CurrentStep:   1,
		StartTime:     now,
		CompletedTime: now}

	err := c.FermentationService().Save(n)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(n, exp) {
		t.Fatalf("unexpected fermentation: %#v", n)
	}
}

func testFermentationService_Save_ErrFermentationRequired(t *testing.T) {
	s, c := MustOpenServerAndClient()
	defer func() { _ = s.Close() }()

	if err := c.FermentationService().Save(nil); err != zymurgauge.ErrFermentationRequired {
		t.Fatal(err)
	}
}

func testFermentationService_Save_ErrInternal(t *testing.T) {
	s, c := MustOpenServerAndClient()
	defer func() { _ = s.Close() }()

	s.FermentationServiceMock.SaveFn = func(b *zymurgauge.Fermentation) error {
		return errors.New("some internal error")
	}

	if err := c.FermentationService().Save(&zymurgauge.Fermentation{ID: 0}); err != zymurgauge.ErrInternal {
		t.Fatal(err)
	}
}

// ToDo: Move to its own service
// func TestFermentationService_LogEvent(t *testing.T) {
// 	t.Run("OK", testFermentationService_LogEvent)
// 	// t.Run("NotFound", testFermentationService_LogEvent_NotFound)
// 	// t.Run("ErrInternal", testFermentationService_LogEvent_ErrInternal)
// }

// func testFermentationService_LogEvent(t *testing.T) {
// 	s, c := MustOpenServerAndClient()
// 	defer func() { _ = s.Close() }()

// 	var now = time.Now()

// 	s.FermentationServiceMock.LogEventFn = func(fermentationID uint64, event zymurgauge.FermentationEvent) error {
// 		return nil
// 	}

// 	e := zymurgauge.FermentationEvent{Time: now, Type: "SomeType", Value: "SomeValue"}

// 	err := c.FermentationService().LogEvent(1, e)

// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// func testFermentationService_LogEvent_NotFound(t *testing.T) {
// 	s, c := MustOpenServerAndClient()
// 	defer s.Close()

// 	var now = time.Now()

// 	s.FermentationServiceMock.LogEventFn = func(fermentationID uint64, event zymurgauge.FermentationEvent) error {
// 		return zymurgauge.ErrNotFound
// 	}

// 	e := zymurgauge.FermentationEvent{Time: now, Type: "SomeType", Value: "SomeValue"}

// 	if err := c.FermentationService().LogEvent(88, e); err != zymurgauge.ErrNotFound {
// 		t.Fatal(err)
// 	}
// }

// func testFermentationService_LogEvent_ErrInternal(t *testing.T) {
// 	s, c := MustOpenServerAndClient()
// 	defer s.Close()

// 	var now = time.Now()

// 	s.FermentationServiceMock.LogEventFn = func(fermentationID uint64, event zymurgauge.FermentationEvent) error {
// 		return errors.New("some internal error")
// 	}

// 	e := zymurgauge.FermentationEvent{Time: now, Type: "SomeType", Value: "SomeValue"}

// 	if err := c.FermentationService().LogEvent(88, e); err != zymurgauge.ErrInternal {
// 		t.Fatal(err)
// 	}

// }
