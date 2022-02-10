package handlers_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/benjaminbartels/zymurgauge/cmd/zym/handlers"
	"github.com/benjaminbartels/zymurgauge/internal/test/mocks"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//nolint: paralleltest // False positives with r.Run not in a loop
func TestGet(t *testing.T) {
	t.Parallel()
	t.Run("getIndex", getIndex)
	t.Run("getStatic", getStatic)
	t.Run("getReadFileError", getReadFileError)
	t.Run("getWriteError", getWriteError)
}

func getIndex(t *testing.T) {
	t.Parallel()

	response := "some response"

	b := []byte(response)

	w, r, ctx := setupHandlerTest("", nil)

	r.URL.Path = "/ui/somepath"

	fsMock := &mocks.FileReader{}
	fsMock.On("ReadFile", "build/index.html").Return(b, nil)

	handler := &handlers.UIHandler{FileReader: fsMock}
	err := handler.Get(ctx, w, r, httprouter.Params{})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	assert.Equal(t, []byte(response), bodyBytes)
}

func getStatic(t *testing.T) {
	t.Parallel()

	response := "some response"

	b := []byte(response)

	w, r, ctx := setupHandlerTest("", nil)

	r.URL.Path = "/ui/static/js/somefile.js"

	fsMock := &mocks.FileReader{}
	fsMock.On("ReadFile", "build/static/js/somefile.js").Return(b, nil)

	handler := &handlers.UIHandler{FileReader: fsMock}
	err := handler.Get(ctx, w, r, httprouter.Params{})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	assert.Equal(t, []byte(response), bodyBytes)
}

func getReadFileError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)

	fsMock := &mocks.FileReader{}
	fsMock.On("ReadFile", mock.Anything).Return(nil, errors.New("fileReaderMock error"))

	handler := &handlers.UIHandler{FileReader: fsMock}
	err := handler.Get(ctx, w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), "could not read file: fileReaderMock error")
}

func getWriteError(t *testing.T) {
	t.Parallel()

	_, r, ctx := setupHandlerTest("", nil)

	fsMock := &mocks.FileReader{}
	fsMock.On("ReadFile", mock.Anything).Return(nil, nil)

	w := mockResponseWriter{}

	handler := &handlers.UIHandler{FileReader: fsMock}
	err := handler.Get(ctx, w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), "could not write response: write error")
}

var _ http.ResponseWriter = (*mockResponseWriter)(nil)

type mockResponseWriter struct {
	http.ResponseWriter
}

func (m mockResponseWriter) Write(bytes []byte) (int, error) {
	return 0, errors.New("write error")
}

func (m mockResponseWriter) Header() http.Header {
	return make(map[string][]string)
}
