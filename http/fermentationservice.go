package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/orangesword/zymurgauge"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// check implementation at compile time
var _ zymurgauge.FermentationService = &FermentationService{}

// FermentationService is the HTTP implementation of zymurgauge.FermentationService
type FermentationService struct {
	url    url.URL
	logger *logrus.Logger
}

// Get returns a fermentation by id
func (s *FermentationService) Get(id uint64) (*zymurgauge.Fermentation, error) {

	u := s.url
	u.Path = "/api/fermentations/" + url.QueryEscape(strconv.FormatUint(id, 10))

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GET Fermentation %d", id)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			s.logger.Error(err)
		}
	}()

	var respBody getFermentationResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, errors.Wrapf(err, "Could not decode Fermentation %d", id)
	} else if respBody.Err != "" {
		return nil, zymurgauge.Error(respBody.Err)
	}

	return respBody.Fermentation, nil
}

// Save creates or updates the stored fermentation with the given Fermentation
func (s *FermentationService) Save(f *zymurgauge.Fermentation) error {

	if f == nil {
		return zymurgauge.ErrFermentationRequired
	}

	u := s.url
	u.Path = "/api/fermentations"

	reqBody, err := json.Marshal(postFermentationRequest{Fermentation: f})
	if err != nil {
		return errors.Wrapf(err, "Could not marshal Fermentation %d", f.ID)
	}

	resp, err := http.Post(u.String(), "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return errors.Wrapf(err, "Could not POST Fermentation %d", f.ID)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			s.logger.Error(err)
		}
	}()

	var respBody postFermentationResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return errors.Wrapf(err, "Could not decode Fermentation %d", f.ID)
	} else if respBody.Err != "" {
		return zymurgauge.Error(respBody.Err)
	}

	*f = *respBody.Fermentation

	return nil
}

// ToDo: Move to its own service
// LogEvent logs the given event for the given Fermentation by it FermentationID
// func (s *FermentationService) LogEvent(fermentationID uint64, event zymurgauge.FermentationEvent) error {
// 	u := *s.URL
// 	u.Path = "/api/fermentations/" + url.QueryEscape(string(fermentationID))

// 	reqBody, err := json.Marshal(patchFermentationRequest{FermentationID: fermentationID, Event: event})
// 	if err != nil {
// 		return err
// 	}

// 	// Create request.
// 	req, err := http.NewRequest("PATCH", u.String(), bytes.NewReader(reqBody))
// 	if err != nil {
// 		return err
// 	}

// 	// Execute request.
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	// Decode response into JSON.
// 	var respBody patchFermentationResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
// 		return err
// 	} else if respBody.Err != "" {
// 		return zymurgauge.Error(respBody.Err)
// 	}

// 	return nil
// }
