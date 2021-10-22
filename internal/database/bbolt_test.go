package database_test

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/database"
	"go.etcd.io/bbolt"
)

// TestClient is a wrapper around the bbolt.Client.
type testDB struct {
	db          *bbolt.DB
	chamberRepo *database.ChamberRepo
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

	db, err := bbolt.Open(f.Name(), 0o666, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}

	chamberRepo, err := database.NewChamberRepo(db)
	if err != nil {
		panic(err)
	}

	t := &testDB{
		db:          db,
		chamberRepo: chamberRepo,
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