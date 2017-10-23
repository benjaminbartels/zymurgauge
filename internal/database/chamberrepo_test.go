package database_test

import (
	"testing"

	"github.com/benjaminbartels/zymurgauge/internal"
)

func TestChamberService_Save_New(t *testing.T) {

	mac := "00:11:22:33:44:55"

	testDB := createTestDB()
	defer func() { testDB.Close() }()

	c := internal.Chamber{
		Name:       "My Chamber",
		MacAddress: mac,
		Thermostat: &internal.Thermostat{
			ThermometerID: "blah",
		},
	}

	if err := testDB.chamberRepo.Save(&c); err != nil {
		t.Fatal(err)
	} else if c.MacAddress != mac {
		t.Fatalf("unexpected mac: %s", c.MacAddress)
	}

	// other, err := testDB.chamberRepo.Get(mac)
	// if err != nil {
	// 	t.Fatal(err)
	// } else if !reflect.DeepEqual(&c, other) {
	// 	t.Fatalf("unexpected controller: %#v vs %#v", &c, other)
	// }
}

func TestChamberService_Save_Existing(t *testing.T) {

	testDB := createTestDB()
	defer func() { testDB.Close() }()

	c1 := &internal.Chamber{Name: "My Chamber 1", MacAddress: "00:11:22:33:44:55"}
	c2 := &internal.Chamber{Name: "My Chamber 2", MacAddress: "aa:bb:cc:dd:ee:ff"}

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

	if uc1, err := testDB.chamberRepo.Get(c1.MacAddress); err != nil {
		t.Fatal(err)
	} else if uc1.MacAddress != c1.MacAddress {
		t.Fatalf("unexpected controller #1 MacAddress: %s", uc1.MacAddress)
	}

	if uc2, err := testDB.chamberRepo.Get(c2.MacAddress); err != nil {
		t.Fatal(err)
	} else if uc2.MacAddress != c2.MacAddress {
		t.Fatalf("unexpected controller #2 MacAddress: %s", uc2.MacAddress)
	}
}
