package database_test

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/database"
	"github.com/boltdb/bolt"
)

// TestClient is a wrapper around the bolt.Client.
type testDB struct {
	db               *bolt.DB
	chamberRepo      *database.ChamberRepo
	beerRepo         *database.BeerRepo
	fermentationRepo *database.FermentationRepo
}

func createTestDB() *testDB {

	p := "zymurgauge-test-"
	f, err := ioutil.TempFile("", p)
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}

	db, err := bolt.Open(f.Name(), 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}

	chamberRepo, err := database.NewChamberRepo(db)
	if err != nil {
		panic(err)
	}

	beerRepo, err := database.NewBeerRepo(db)
	if err != nil {
		panic(err)
	}

	fermentationRepo, err := database.NewFermentationRepo(db)
	if err != nil {
		panic(err)
	}

	t := &testDB{
		db:               db,
		chamberRepo:      chamberRepo,
		beerRepo:         beerRepo,
		fermentationRepo: fermentationRepo,
	}

	return t
}

func (t *testDB) Close() {
	t.db.Close()
	os.Remove(t.db.Path())
}
