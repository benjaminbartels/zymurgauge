package web_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/pkg/errors"
)

func TestRespond(t *testing.T) {
	t.Run("OK", testRespondOK)
	t.Run("StatusNoContent", testRespondStatusNoContent)
	t.Run("WriteError", testRespondWriteError)
	t.Run("MarshalError", testRespondMarshalError)
}

func testRespondOK(t *testing.T) {
	rw := &responseWriterMock{}

	var code int

	rw.HeaderFn = func() http.Header {
		return make(map[string][]string)
	}

	rw.WriteFn = func(bytes []byte) (int, error) {
		return len(bytes), nil
	}

	rw.WriteHeaderFn = func(statusCode int) {
		code = statusCode
	}

	b := &internal.Beer{
		Name: "Golden Stout",
	}

	r, _ := http.NewRequest(http.MethodGet, "/beer/1", nil)
	ctx := context.WithValue(r.Context(), web.CtxValuesKey, &web.CtxValues{StartTime: time.Now()})
	err := web.Respond(ctx, rw, b, http.StatusOK)

	if err != nil {
		t.Fatal(err)
	} else if !rw.HeaderInvoked {
		t.Fatal("Header not invoked")
	} else if !rw.WriteInvoked {
		t.Fatal("Write not invoked")
	} else if !rw.WriteHeaderInvoked {
		t.Fatal("WriteHeader not invoked")
	} else if code != http.StatusOK {
		t.Fatalf("Unexpected StatusCode %d", code)
	}

}

func testRespondStatusNoContent(t *testing.T) {
	rw := &responseWriterMock{}

	var code int

	rw.WriteHeaderFn = func(statusCode int) {
		code = statusCode
	}

	r, _ := http.NewRequest(http.MethodGet, "/beer", nil)
	ctx := context.WithValue(r.Context(), web.CtxValuesKey, &web.CtxValues{StartTime: time.Now()})
	err := web.Respond(ctx, rw, nil, http.StatusNoContent)

	if err != nil {
		t.Fatal(err)
	} else if rw.HeaderInvoked {
		t.Fatal("Header invoked")
	} else if rw.WriteInvoked {
		t.Fatal("Write invoked")
	} else if !rw.WriteHeaderInvoked {
		t.Fatal("WriteHeader not invoked")
	} else if code != http.StatusNoContent {
		t.Fatalf("Unexpected StatusCode %d", code)
	}

}

func testRespondWriteError(t *testing.T) {
	rw := &responseWriterMock{}

	var writeError = errors.New("write error occurred")

	rw.HeaderFn = func() http.Header {
		return make(map[string][]string)
	}

	rw.WriteFn = func(bytes []byte) (int, error) {
		return 0, writeError
	}

	rw.WriteHeaderFn = func(statusCode int) {}

	b := &internal.Beer{
		Name: "Golden Stout",
	}

	r, _ := http.NewRequest(http.MethodGet, "/beer/1", nil)
	ctx := context.WithValue(r.Context(), web.CtxValuesKey, &web.CtxValues{StartTime: time.Now()})
	err := web.Respond(ctx, rw, b, http.StatusOK)

	if err != writeError {
		t.Fatalf("Expected Error %v was %v", writeError, err)
	} else if !rw.HeaderInvoked {
		t.Fatal("Header not invoked")
	} else if !rw.WriteInvoked {
		t.Fatal("Write not invoked")
	} else if !rw.WriteHeaderInvoked {
		t.Fatal("WriteHeader not invoked")
	}
}

func testRespondMarshalError(t *testing.T) {
	rw := &responseWriterMock{}

	bogusData := make(chan int)

	r, _ := http.NewRequest(http.MethodGet, "/beer/1", nil)
	ctx := context.WithValue(r.Context(), web.CtxValuesKey, &web.CtxValues{StartTime: time.Now()})
	err := web.Respond(ctx, rw, bogusData, http.StatusOK)

	if _, ok := err.(*json.UnsupportedTypeError); !ok {
		t.Errorf("Unexpected Error %v", err)
	} else if rw.HeaderInvoked {
		t.Fatal("Header invoked")
	} else if rw.WriteInvoked {
		t.Fatal("Write invoked")
	} else if rw.WriteHeaderInvoked {
		t.Fatal("WriteHeader invoked")
	}
}

func TestError(t *testing.T) {
	t.Run("ErrNotFound", testErrorErrNotFound)
	t.Run("CatchAll", testErrorCatchAll)
	t.Run("ErrorFromRespond", testErrorErrorFromRespond)
}

func testErrorErrNotFound(t *testing.T) {
	rw := &responseWriterMock{}

	var code int

	rw.HeaderFn = func() http.Header {
		return make(map[string][]string)
	}

	rw.WriteFn = func(bytes []byte) (int, error) {
		return len(bytes), nil
	}

	rw.WriteHeaderFn = func(statusCode int) {
		code = statusCode
	}

	r, _ := http.NewRequest(http.MethodGet, "/beer/1", nil)
	ctx := context.WithValue(r.Context(), web.CtxValuesKey, &web.CtxValues{StartTime: time.Now()})
	err := web.Error(ctx, rw, web.ErrNotFound)

	if err != nil {
		t.Fatal(err)
	} else if !rw.HeaderInvoked {
		t.Fatal("Header not invoked")
	} else if !rw.WriteInvoked {
		t.Fatal("Write not invoked")
	} else if !rw.WriteHeaderInvoked {
		t.Fatal("WriteHeader not invoked")
	} else if code != http.StatusNotFound {
		t.Fatalf("Unexpected StatusCode %d", code)
	}

}

func testErrorCatchAll(t *testing.T) {
	rw := &responseWriterMock{}

	var code int

	var someError = errors.New("some error")

	rw.HeaderFn = func() http.Header {
		return make(map[string][]string)
	}

	rw.WriteFn = func(bytes []byte) (int, error) {
		return len(bytes), nil
	}

	rw.WriteHeaderFn = func(statusCode int) {
		code = statusCode
	}

	r, _ := http.NewRequest(http.MethodGet, "/beer/1", nil)
	ctx := context.WithValue(r.Context(), web.CtxValuesKey, &web.CtxValues{StartTime: time.Now()})
	err := web.Error(ctx, rw, someError)

	if err != nil {
		t.Fatal(err)
	} else if !rw.HeaderInvoked {
		t.Fatal("Header not invoked")
	} else if !rw.WriteInvoked {
		t.Fatal("Write not invoked")
	} else if !rw.WriteHeaderInvoked {
		t.Fatal("WriteHeader not invoked")
	} else if code != http.StatusInternalServerError {
		t.Fatalf("Unexpected StatusCode %d", code)
	}

}

func testErrorErrorFromRespond(t *testing.T) {
	rw := &responseWriterMock{}

	var code int

	var writeError = errors.New("write error occurred")

	rw.HeaderFn = func() http.Header {
		return make(map[string][]string)
	}

	rw.WriteFn = func(bytes []byte) (int, error) {
		return 0, writeError
	}

	rw.WriteHeaderFn = func(statusCode int) {
		code = statusCode
	}

	r, _ := http.NewRequest(http.MethodGet, "/beer/1", nil)
	ctx := context.WithValue(r.Context(), web.CtxValuesKey, &web.CtxValues{StartTime: time.Now()})
	err := web.Error(ctx, rw, web.ErrNotFound)

	if err != writeError {
		t.Fatalf("Expected Error %v was %v", writeError, err)
	} else if !rw.HeaderInvoked {
		t.Fatal("Header not invoked")
	} else if !rw.WriteInvoked {
		t.Fatal("Write not invoked")
	} else if !rw.WriteHeaderInvoked {
		t.Fatal("WriteHeader not invoked")
	} else if code != http.StatusNotFound {
		t.Fatalf("Unexpected StatusCode %d", code)
	}

}

type responseWriterMock struct {
	HeaderFn           func() http.Header
	WriteFn            func(bytes []byte) (int, error)
	WriteHeaderFn      func(statusCode int)
	HeaderInvoked      bool
	WriteInvoked       bool
	WriteHeaderInvoked bool
}

func (r *responseWriterMock) Header() http.Header {
	r.HeaderInvoked = true
	return r.HeaderFn()
}

func (r *responseWriterMock) Write(bytes []byte) (int, error) {
	r.WriteInvoked = true
	return r.WriteFn(bytes)
}

func (r *responseWriterMock) WriteHeader(statusCode int) {
	r.WriteHeaderInvoked = true
	r.WriteHeaderFn(statusCode)
}
