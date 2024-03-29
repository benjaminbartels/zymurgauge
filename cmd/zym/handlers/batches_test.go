package handlers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/benjaminbartels/zymurgauge/cmd/zym/handlers"
	"github.com/benjaminbartels/zymurgauge/internal/batch"
	"github.com/benjaminbartels/zymurgauge/internal/brewfather"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/benjaminbartels/zymurgauge/internal/test/mocks"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

//nolint:paralleltest // False positives with r.Run not in a loop
func TestGetAllBatches(t *testing.T) {
	t.Parallel()
	t.Run("getAllBatches", getAllBatches)
	t.Run("getAllBatchesEmpty", getAllBatchesEmpty)
	t.Run("getAllServiceError", getAllServiceError)
	t.Run("getAllBatchesRespondError", getAllBatchesRespondError)
}

func getAllBatches(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)

	expected := []batch.Summary{
		{ID: batchID},
		{ID: "f4ce0e05-1ada-42b8-8fc4-fb3482525d0d"},
	}
	serviceMock := &mocks.Service{}
	serviceMock.On("GetAllBatchSummaries", ctx).Return([]brewfather.BatchSummary{
		{ID: batchID},
		{ID: "f4ce0e05-1ada-42b8-8fc4-fb3482525d0d"},
	}, nil)

	handler := &handlers.BatchesHandler{Service: serviceMock}
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	result := []batch.Summary{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func getAllBatchesEmpty(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)

	expected := []batch.Summary{}
	serviceMock := &mocks.Service{}
	serviceMock.On("GetAllBatchSummaries", ctx).Return([]brewfather.BatchSummary{}, nil)

	handler := &handlers.BatchesHandler{Service: serviceMock}
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	result := []batch.Summary{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func getAllServiceError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)

	serviceMock := &mocks.Service{}
	serviceMock.On("GetAllBatchSummaries", ctx).Return([]brewfather.BatchSummary{}, errSomeError)

	handler := &handlers.BatchesHandler{Service: serviceMock}
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	// TODO: Waiting on PR for ErrorContains(): https://github.com/stretchr/testify/pull/1022
	assert.Contains(t, err.Error(), fmt.Sprintf(repoErrMsg, "get all batches from"))
}

func getAllBatchesRespondError(t *testing.T) {
	t.Parallel()

	w, r, _ := setupHandlerTest("", nil)
	ctx := context.Background()

	serviceMock := &mocks.Service{}
	serviceMock.On("GetAllBatchSummaries", ctx).Return([]brewfather.BatchSummary{}, nil)

	handler := &handlers.BatchesHandler{Service: serviceMock}
	// use new ctx to force error
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), respondErrMsg)
}

//nolint:paralleltest // False positives with r.Run not in a loop
func TestGetBatch(t *testing.T) {
	t.Parallel()
	t.Run("getBatchFound", getBatchFound)
	t.Run("getBatchNotFoundError", getBatchNotFoundError)
	t.Run("getBatchIsNil", getBatchIsNil)
	t.Run("getServiceError", getServiceError)
	t.Run("getBatchRespondError", getBatchRespondError)
}

func getBatchFound(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)

	expected := batch.Detail{
		ID: batchID,
		Recipe: batch.Recipe{
			Fermentation: batch.Fermentation{
				Steps: []batch.FermentationStep{},
			},
		},
	}
	serviceMock := &mocks.Service{}
	serviceMock.On("GetBatchDetail", ctx, batchID).Return(&brewfather.BatchDetail{ID: batchID}, nil)

	handler := &handlers.BatchesHandler{Service: serviceMock}
	err := handler.Get(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: batchID}})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	batch := batch.Detail{}
	err = json.Unmarshal(bodyBytes, &batch)
	assert.NoError(t, err)
	assert.Equal(t, expected, batch)
}

func getBatchNotFoundError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)

	serviceMock := &mocks.Service{}
	serviceMock.On("GetBatchDetail", ctx, batchID).Return(nil, brewfather.ErrNotFound)

	handler := &handlers.BatchesHandler{Service: serviceMock}
	err := handler.Get(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: batchID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(notFoundErrorMsg, "batch", batchID))

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusNotFound)
}

func getBatchIsNil(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)

	serviceMock := &mocks.Service{}
	serviceMock.On("GetBatchDetail", ctx, batchID).Return(nil, nil)

	handler := &handlers.BatchesHandler{Service: serviceMock}
	err := handler.Get(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: batchID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(notFoundErrorMsg, "batch", batchID))

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusNotFound)
}

func getServiceError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)

	serviceMock := &mocks.Service{}
	serviceMock.On("GetBatchDetail", ctx, batchID).Return(&brewfather.BatchDetail{}, errSomeError)

	handler := &handlers.BatchesHandler{Service: serviceMock}
	err := handler.Get(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: batchID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(repoErrMsg, "get batch from"))
}

func getBatchRespondError(t *testing.T) {
	t.Parallel()

	w, r, _ := setupHandlerTest("", nil)
	ctx := context.Background()

	serviceMock := &mocks.Service{}
	serviceMock.On("GetBatchDetail", ctx, batchID).Return(&brewfather.BatchDetail{ID: batchID}, nil)

	handler := &handlers.BatchesHandler{Service: serviceMock}
	// use new ctx to force error
	err := handler.Get(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: batchID}})
	assert.Contains(t, err.Error(), respondErrMsg)
}
