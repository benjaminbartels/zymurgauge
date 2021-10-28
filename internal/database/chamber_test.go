package database_test

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/benjaminbartels/zymurgauge/internal/chamber"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/bbolt"
)

//nolint: paralleltest // False positives with r.Run not in a loop
func TestGetAllChambers(t *testing.T) {
	t.Parallel()
	t.Run("saveNewChamber", saveNewChamber)
	t.Run("saveExistingChamber", saveExistingChamber)
	t.Run("saveChamberPutError", saveChamberPutError)
}

func saveNewChamber(t *testing.T) {
	t.Parallel()

	testDB := createTestDB()

	defer func() { testDB.Close() }()

	c := chamber.Chamber{
		Name: "My Chamber",
	}

	err := testDB.chamberRepo.Save(&c)

	assert.NoError(t, err)
}

func saveExistingChamber(t *testing.T) {
	t.Parallel()

	testDB := createTestDB()

	defer func() { testDB.Close() }()

	c1 := &chamber.Chamber{Name: "My Chamber 1", ID: "59679696-1263-4340-a256-6c46876b4a13"}
	c2 := &chamber.Chamber{Name: "My Chamber 2", ID: "d9d075b4-6b45-44cc-945b-c5b9ce13e442"}

	err := testDB.chamberRepo.Save(c1)
	assert.NoError(t, err)

	err = testDB.chamberRepo.Save(c2)
	assert.NoError(t, err)

	uc1, err := testDB.chamberRepo.Get(c1.ID)
	assert.NoError(t, err)
	assert.Equal(t, c1.ID, uc1.ID)

	uc2, err := testDB.chamberRepo.Get(c2.ID)
	assert.NoError(t, err)
	assert.Equal(t, c2.ID, uc2.ID)
}

func saveChamberPutError(t *testing.T) {
	t.Parallel()

	testDB := createTestDB()

	defer func() { testDB.Close() }()

	c := chamber.Chamber{
		ID:   generateRandomString(bbolt.MaxKeySize + 1),
		Name: "My Chamber",
	}

	err := testDB.chamberRepo.Save(&c)
	// TODO: Waiting on PR for ErrorContains(): https://github.com/stretchr/testify/pull/1022
	assert.Contains(t, err.Error(), "could not execute update transaction: could not put Chamber")
}

// func randSeq(n int) string {
// 	letters := []rune("0123456789abcdf")
// 	b := make([]rune, n)

// 	for i := range b {
// 		b[i] = letters[rand.Intn(len(letters))]
// 	}

// 	return string(b)
// }

func generateRandomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	b := make([]byte, n)

	for i := 0; i < n; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		b[i] = letters[num.Int64()]
	}

	return string(b)
}
