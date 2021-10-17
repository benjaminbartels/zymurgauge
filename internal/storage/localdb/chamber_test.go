package localdb_test

import (
	"testing"

	"github.com/benjaminbartels/zymurgauge/internal/chamber"
)

func TestChamberServiceSaveNew(t *testing.T) {
	t.Parallel()

	id := "59679696-1263-4340-a256-6c46876b4a13"

	testDB := createTestDB()

	defer func() { testDB.Close() }()

	c := chamber.Chamber{
		Name: "My Chamber",
		ID:   id,
	}

	if err := testDB.chamberRepo.Save(&c); err != nil {
		t.Fatal(err)
	} else if c.ID != id {
		t.Fatalf("unexpected id: %s", c.ID)
	}
}

func TestChamberServiceSaveExisting(t *testing.T) {
	t.Parallel()

	testDB := createTestDB()

	defer func() { testDB.Close() }()

	c1 := &chamber.Chamber{Name: "My Chamber 1", ID: "59679696-1263-4340-a256-6c46876b4a13"}
	c2 := &chamber.Chamber{Name: "My Chamber 2", ID: "d9d075b4-6b45-44cc-945b-c5b9ce13e442"}

	if err := testDB.chamberRepo.Save(c1); err != nil {
		t.Fatal(err)
	} else if err := testDB.chamberRepo.Save(c2); err != nil {
		t.Fatal(err)
	}

	if err := testDB.chamberRepo.Save(c1); err != nil {
		t.Fatal(err)
	} else if err := testDB.chamberRepo.Save(c2); err != nil {
		t.Fatal(err)
	}

	if uc1, err := testDB.chamberRepo.Get(c1.ID); err != nil {
		t.Fatal(err)
	} else if uc1.ID != c1.ID {
		t.Fatalf("unexpected controller #1 ID: %s", uc1.ID)
	}

	if uc2, err := testDB.chamberRepo.Get(c2.ID); err != nil {
		t.Fatal(err)
	} else if uc2.ID != c2.ID {
		t.Fatalf("unexpected controller #2 ID: %s", uc2.ID)
	}
}
