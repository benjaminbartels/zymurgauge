package bolt_test

import (
	"reflect"
	"testing"

	"time"

	"github.com/orangesword/zymurgauge"
)

func TestFermentationService_Save_New(t *testing.T) {

	c := MustOpenClient()
	defer func() { _ = c.Close() }()

	now := time.Now()

	s := c.FermentationService()

	f := zymurgauge.Fermentation{
		Beer: zymurgauge.Beer{
			Name:    "My Beer",
			ModTime: now,
		},
		StartTime:     now,
		CompletedTime: now,
	}

	if err := s.Save(&f); err != nil {
		t.Fatal(err)
	} else if f.ID != 1 {
		t.Fatalf("unexpected id: %d", f.ID)
	}

	other, err := s.Get(1)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&f, other) {
		t.Fatalf("unexpected fermentation: %#v", other)
	}
}

func TestFermentationService_Save_Existing(t *testing.T) {

	c := MustOpenClient()
	defer func() { _ = c.Close() }()

	s := c.FermentationService()

	now := time.Now()

	f1 := &zymurgauge.Fermentation{
		Beer: zymurgauge.Beer{
			Name:    "My Beer 1",
			ModTime: now,
		},
		StartTime:     now,
		CompletedTime: now,
	}

	f2 := &zymurgauge.Fermentation{
		Beer: zymurgauge.Beer{
			Name:    "My Beer 2",
			ModTime: now,
		},
		StartTime:     now,
		CompletedTime: now,
	}

	if err := s.Save(f1); err != nil {
		t.Fatal(err)
	} else if err := s.Save(f2); err != nil {
		t.Fatal(err)
	}

	f1.CurrentStep = 2
	f2.CurrentStep = 3

	if err := s.Save(f1); err != nil {
		t.Fatal(err)
	} else if err := s.Save(f2); err != nil {
		t.Fatal(err)
	}

	if uf1, err := s.Get(f1.ID); err != nil {
		t.Fatal(err)
	} else if uf1.CurrentStep != f1.CurrentStep {
		t.Fatalf("unexpected fermentation #1 CurrentStep: %f", uf1.CurrentStep)
	}

	if uf2, err := s.Get(f2.ID); err != nil {
		t.Fatal(err)
	} else if uf2.CurrentStep != f2.CurrentStep {
		t.Fatalf("unexpected fermentation #2 CurrentStep: %f", uf2.CurrentStep)
	}
}

// ToDo: Move to its own service
// func TestFermentationService_LogEvent(t *testing.T) {
// 	c := MustOpenClient()
// 	defer func() { _ = c.Close() }()

// 	now := time.Now()

// 	s := c.FermentationService()

// 	f := zymurgauge.Fermentation{
// 		Beer: zymurgauge.Beer{
// 			Name:    "My Beer",
// 			ModTime: now,
// 		},
// 		StartTime:     now,
// 		CompletedTime: now,
// 	}

// 	if err := s.Save(&f); err != nil {
// 		t.Fatal(err)
// 	} else if f.ID != 1 {
// 		t.Fatalf("unexpected id: %d", f.ID)
// 	}

// 	e1 := zymurgauge.FermentationEvent{
// 		Time:  now,
// 		Type:  "SomeEvent1",
// 		Value: "SomeValue1",
// 	}

// 	e2 := zymurgauge.FermentationEvent{
// 		Time:  now,
// 		Type:  "SomeEvent2",
// 		Value: "SomeValue2",
// 	}

// 	_ = s.LogEvent(f.ID, e1)
// 	_ = s.LogEvent(f.ID, e2)

// 	other, err := s.Get(f.ID)
// 	if err != nil {
// 		t.Fatal(err)
// 	} else if len(other.Events) != 2 {
// 		t.Fatalf("unexpected count of events: %#v", len(other.Events))
// 	}
// }

// ToDo: Move to its own service
// func TestFermentationService_LogEvent_NotFound(t *testing.T) {
// 	c := MustOpenClient()
// 	defer func() { _ = c.Close() }()

// 	s := c.FermentationService()

// 	e := zymurgauge.FermentationEvent{
// 		Time:  time.Now(),
// 		Type:  "SomeEvent",
// 		Value: "SomeValue",
// 	}

// 	if err := s.LogEvent(0, e); err != zymurgauge.ErrNotFound {
// 		t.Fatal(err)
// 	}
// }
