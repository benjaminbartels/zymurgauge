package web_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/benjaminbartels/zymurgauge/internal/test/mocks"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//nolint:paralleltest // False positives with r.Run not in a loop
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

	w := httptest.NewRecorder()
	r := httptest.NewRequest("", "/", nil)

	r.URL.Path = "/somepath"

	fsMock := &mocks.FileReader{}
	fsMock.On("ReadFile", "build/index.html").Return(b, nil)

	l, _ := logtest.NewNullLogger()

	app := web.NewApp(nil, fsMock, l)

	app.ServeHTTP(w, r)

	resp := w.Result()
	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	assert.Equal(t, []byte(response), bodyBytes)
}

func getStatic(t *testing.T) {
	t.Parallel()

	type test struct {
		path string
	}

	testCases := []test{
		{path: "/static/js/somefile.js"},
		{path: "/somefile.png"},
		{path: "/somefile.ico"},
	}

	for _, tc := range testCases {
		tc := tc

		response := "some response"

		b := []byte(response)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("", "/", nil)

		r.URL.Path = tc.path

		fsMock := &mocks.FileReader{}
		fsMock.On("ReadFile", "build"+tc.path).Return(b, nil)

		l, _ := logtest.NewNullLogger()

		app := web.NewApp(nil, fsMock, l)

		app.ServeHTTP(w, r)

		resp := w.Result()
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		assert.Equal(t, []byte(response), bodyBytes)
	}
}

func getReadFileError(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("", "/", nil)

	fsMock := &mocks.FileReader{}
	fsMock.On("ReadFile", mock.Anything).Return(nil, errors.New("fileReaderMock error"))

	l, hook := logtest.NewNullLogger()

	app := web.NewApp(nil, fsMock, l)

	app.ServeHTTP(w, r)

	assert.True(t, logContains(hook.AllEntries(), logrus.ErrorLevel, "Could not read file"))
}

func getWriteError(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest("", "/", nil)
	w := mockResponseWriter{}

	fsMock := &mocks.FileReader{}
	fsMock.On("ReadFile", mock.Anything).Return(nil, nil)

	l, hook := logtest.NewNullLogger()

	app := web.NewApp(nil, fsMock, l)

	app.ServeHTTP(w, r)

	assert.True(t, logContains(hook.AllEntries(), logrus.ErrorLevel, "Could not write response"))
}

var _ http.ResponseWriter = (*mockResponseWriter)(nil)

type mockResponseWriter struct {
	http.ResponseWriter
}

func (m mockResponseWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("write error")
}

func (m mockResponseWriter) Header() http.Header {
	return make(map[string][]string)
}

func logContains(logs []*logrus.Entry, level logrus.Level, substr string) bool {
	found := false

	for _, v := range logs {
		if strings.Contains(v.Message, substr) && v.Level == level {
			found = true
		}
	}

	return found
}
