package database_test

import (
	"reflect"
	"testing"

	"github.com/benjaminbartels/zymurgauge/internal"

	"time"
)

func TestFermentationService_Save_New(t *testing.T) {

	testDB := createTestDB()
	defer func() { testDB.Close() }()

	now := time.Now()

	f := internal.Fermentation{
		Beer: internal.Beer{
			Name:    "My Beer",
			ModTime: now,
		},
		StartTime:     now,
		CompletedTime: now,
	}

	if err := testDB.fermentationRepo.Save(&f); err != nil {
		t.Fatal(err)
	} else if f.ID != 0x1 {
		t.Fatalf("unexpected id: %d", f.ID)
	}

	other, err := testDB.fermentationRepo.Get(1)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&f, other) {
		t.Fatalf("unexpected fermentation: %#v", other)
	}
}

func TestFermentationService_Save_Existing(t *testing.T) {

	testDB := createTestDB()
	defer func() { testDB.Close() }()

	now := time.Now()

	f1 := &internal.Fermentation{
		Beer: internal.Beer{
			Name:    "My Beer 1",
			ModTime: now,
		},
		StartTime:     now,
		CompletedTime: now,
	}

	f2 := &internal.Fermentation{
		Beer: internal.Beer{
			Name:    "My Beer 2",
			ModTime: now,
		},
		StartTime:     now,
		CompletedTime: now,
	}

	if err := testDB.fermentationRepo.Save(f1); err != nil {
		t.Fatal(err)
	} else if err := testDB.fermentationRepo.Save(f2); err != nil {
		t.Fatal(err)
	}

	f1.CurrentStep = 2
	f2.CurrentStep = 3

	if err := testDB.fermentationRepo.Save(f1); err != nil {
		t.Fatal(err)
	} else if err := testDB.fermentationRepo.Save(f2); err != nil {
		t.Fatal(err)
	}

	if uf1, err := testDB.fermentationRepo.Get(f1.ID); err != nil {
		t.Fatal(err)
	} else if uf1.CurrentStep != f1.CurrentStep {
		t.Fatalf("unexpected fermentation #1 CurrentStep: %s", uf1.CurrentStep)
	}

	if uf2, err := testDB.fermentationRepo.Get(f2.ID); err != nil {
		t.Fatal(err)
	} else if uf2.CurrentStep != f2.CurrentStep {
		t.Fatalf("unexpected fermentation #2 CurrentStep: %s", uf2.CurrentStep)
	}
}
