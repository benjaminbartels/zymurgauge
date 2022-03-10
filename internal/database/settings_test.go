package database_test

import (
	"testing"

	"github.com/benjaminbartels/zymurgauge/internal/settings"
	"github.com/stretchr/testify/assert"
)

func TestSaveSettings(t *testing.T) {
	t.Parallel()
	t.Run("saveSettings", saveSettings)
}

func saveSettings(t *testing.T) {
	t.Parallel()

	testDB := createTestDB()

	defer func() { testDB.Close() }()

	s := settings.Settings{
		BrewfatherAPIUserID: "someID",
		BrewfatherAPIKey:    "someKey",
		BrewfatherLogURL:    "https://someurl.com",
		TemperatureUnits:    "C",
	}

	err := testDB.settingsRepo.Save(&s)

	assert.NoError(t, err)
}
