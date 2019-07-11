package boltdb_test

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/database/boltdb"
	"github.com/boltdb/bolt"
)

// TestClient is a wrapper around the bolt.Client.
type testDB struct {
	db               *bolt.DB
	chamberRepo      *boltdb.ChamberRepo
	beerRepo         *boltdb.BeerRepo
	fermentationRepo *boltdb.FermentationRepo
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

	chamberRepo, err := boltdb.NewChamberRepo(db)
	if err != nil {
		panic(err)
	}

	beerRepo, err := boltdb.NewBeerRepo(db)
	if err != nil {
		panic(err)
	}

	fermentationRepo, err := boltdb.NewFermentationRepo(db)
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
	p := t.db.Path()
	if err := t.db.Close(); err != nil {
		panic(err)
	}
	if err := os.Remove(p); err != nil {
		panic(err)
	}
}
