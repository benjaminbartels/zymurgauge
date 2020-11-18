package storage_test

import (
	"testing"

	"github.com/benjaminbartels/zymurgauge/internal/storage"
)

func TestBeerServiceSaveNew(t *testing.T) {
	testDB := createTestDB()

	defer func() { testDB.Close() }()

	b := storage.Beer{
		Name: "My Beer",
	}

	if err := testDB.beerRepo.Save(&b); err != nil {
		t.Fatal(err)
	} else if b.ID != 1 {
		t.Fatalf("unexpected id: %d", b.ID)
	}
}

func TestBeerServiceSaveExisting(t *testing.T) {
	testDB := createTestDB()

	defer func() { testDB.Close() }()

	b1 := &storage.Beer{Name: "My Beer 1"}
	b2 := &storage.Beer{Name: "My Beer 2"}

	if err := testDB.beerRepo.Save(b1); err != nil {
		t.Fatal(err)
	} else if err := testDB.beerRepo.Save(b2); err != nil {
		t.Fatal(err)
	}

	b1.Style = "Lager"
	b2.Style = "Porter"

	if err := testDB.beerRepo.Save(b1); err != nil {
		t.Fatal(err)
	} else if err := testDB.beerRepo.Save(b2); err != nil {
		t.Fatal(err)
	}

	if ub1, err := testDB.beerRepo.Get(b1.ID); err != nil {
		t.Fatal(err)
	} else if ub1.Style != b1.Style {
		t.Fatalf("unexpected beer #1 style: %s", ub1.Style)
	}

	if ub2, err := testDB.beerRepo.Get(b2.ID); err != nil {
		t.Fatal(err)
	} else if ub2.Style != b2.Style {
		t.Fatalf("unexpected beer #2 style: %s", ub2.Style)
	}
}