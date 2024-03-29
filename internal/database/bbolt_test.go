package database_test

import (
	"os"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/database"
	"go.etcd.io/bbolt"
)

// TestClient is a wrapper around the bbolt.Client.
type testDB struct {
	db           *bbolt.DB
	chamberRepo  *database.ChamberRepo
	settingsRepo *database.SettingsRepo
}

func createTestDB() *testDB {
	p := "zymurgauge-test-"

	f, err := os.CreateTemp("", p)
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

	settingsRepo, err := database.NewSettingsRepo(db)
	if err != nil {
		panic(err)
	}

	t := &testDB{
		db:           db,
		chamberRepo:  chamberRepo,
		settingsRepo: settingsRepo,
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
