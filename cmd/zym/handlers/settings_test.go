package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/benjaminbartels/zymurgauge/cmd/zym/handlers"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/benjaminbartels/zymurgauge/internal/settings"
	"github.com/benjaminbartels/zymurgauge/internal/test/mocks"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func getTestSettings() *settings.Settings {
	return &settings.Settings{
		AppSettings: settings.AppSettings{
			BrewfatherAPIUserID: "someID",
			BrewfatherAPIKey:    "someKey",
			BrewfatherLogURL:    "https://someurl.com",
			TemperatureUnits:    "Celsius",
		},
	}
}

//nolint: paralleltest // False positives with r.Run not in a loop
func TestGetSettings(t *testing.T) {
	t.Parallel()
	t.Run("getSettingsFound", getSettingsFound)
	t.Run("getSettingsNotFoundError", getSettingsNotFoundError)
	t.Run("getSettingsOtherError", getSettingsOtherError)
	t.Run("getSettingsRespondError", getSettingsRespondError)
}

func getSettingsFound(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)

	s := getTestSettings()
	settingsMock := &mocks.SettingsRepo{}
	settingsMock.On("Get").Return(s, nil)

	ch := make(chan settings.Settings)

	handler := &handlers.SettingsHandler{SettingsRepo: settingsMock, UpdateChan: ch}

	err := handler.Get(ctx, w, r, httprouter.Params{})
	assert.NoError(t, err)

	resp := w.Result()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	result := &settings.Settings{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assert.Equal(t, s, result)
}

func getSettingsNotFoundError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)

	settingsMock := &mocks.SettingsRepo{}
	settingsMock.On("Get").Return(nil, nil)

	ch := make(chan settings.Settings)

	handler := &handlers.SettingsHandler{SettingsRepo: settingsMock, UpdateChan: ch}

	err := handler.Get(ctx, w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), settingsNotFound)

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusNotFound)
}

func getSettingsOtherError(t *testing.T) {
	t.Parallel()

	s := getTestSettings()
	jsonBytes, _ := json.Marshal(s)

	w, r, ctx := setupHandlerTest("", bytes.NewBuffer(jsonBytes))

	settingsMock := &mocks.SettingsRepo{}
	settingsMock.On("Get").Return(nil, errors.New("settingsMock error"))

	ch := make(chan settings.Settings)

	handler := &handlers.SettingsHandler{SettingsRepo: settingsMock, UpdateChan: ch}

	err := handler.Get(ctx, w, r, httprouter.Params{httprouter.Param{}})
	assert.Contains(t, err.Error(), fmt.Sprintf(repoErrMsg, "get settings from"))
}

func getSettingsRespondError(t *testing.T) {
	t.Parallel()

	s := getTestSettings()
	jsonBytes, _ := json.Marshal(s)

	w, r, _ := setupHandlerTest("", bytes.NewBuffer(jsonBytes))

	settingsMock := &mocks.SettingsRepo{}
	settingsMock.On("Get").Return(s, nil)

	ch := make(chan settings.Settings)

	handler := &handlers.SettingsHandler{SettingsRepo: settingsMock, UpdateChan: ch}

	// use new ctx to force error
	err := handler.Get(context.Background(), w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), respondErrMsg)
}

//nolint: paralleltest // False positives with r.Run not in a loop
func TestSaveSettings(t *testing.T) {
	t.Parallel()
	t.Run("saveSettings", saveSettings)
	t.Run("saveSettingsParseError", saveSettingsParseError)
	t.Run("saveSettingsOtherError", saveSettingsOtherError)
	t.Run("saveSettingsRespondError", saveSettingsRespondError)
}

func saveSettings(t *testing.T) {
	t.Parallel()

	s := getTestSettings()
	jsonBytes, _ := json.Marshal(s)

	w, r, ctx := setupHandlerTest("", bytes.NewBuffer(jsonBytes))

	settingsMock := &mocks.SettingsRepo{}
	settingsMock.On("Get").Return(s, nil)
	settingsMock.On("Save", s).Return(nil)

	ch := make(chan settings.Settings)

	go func() {
		<-ch
	}()

	handler := &handlers.SettingsHandler{SettingsRepo: settingsMock, UpdateChan: ch}

	err := handler.Save(ctx, w, r, httprouter.Params{})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	result := &settings.Settings{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
}

func saveSettingsParseError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)

	settingsMock := &mocks.SettingsRepo{}

	ch := make(chan settings.Settings)

	handler := &handlers.SettingsHandler{SettingsRepo: settingsMock, UpdateChan: ch}

	err := handler.Save(ctx, w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), parseErrorMsg)
}

func saveSettingsOtherError(t *testing.T) {
	t.Parallel()

	s := getTestSettings()
	jsonBytes, _ := json.Marshal(s)

	w, r, ctx := setupHandlerTest("", bytes.NewBuffer(jsonBytes))

	settingsMock := &mocks.SettingsRepo{}
	settingsMock.On("Get").Return(s, nil)
	settingsMock.On("Save", s).Return(errors.New("controllerMock error"))

	ch := make(chan settings.Settings)

	handler := &handlers.SettingsHandler{SettingsRepo: settingsMock, UpdateChan: ch}

	err := handler.Save(ctx, w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), fmt.Sprintf(repoErrMsg, "save settings to"))
}

func saveSettingsRespondError(t *testing.T) {
	t.Parallel()

	s := getTestSettings()
	jsonBytes, _ := json.Marshal(s)

	w, r, _ := setupHandlerTest("", bytes.NewBuffer(jsonBytes))

	settingsMock := &mocks.SettingsRepo{}
	settingsMock.On("Get").Return(s, nil)
	settingsMock.On("Save", s).Return(nil)

	ch := make(chan settings.Settings)

	handler := &handlers.SettingsHandler{SettingsRepo: settingsMock, UpdateChan: ch}

	// use new ctx to force error
	err := handler.Save(context.Background(), w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), respondErrMsg)
}
