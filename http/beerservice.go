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
var _ zymurgauge.BeerService = &BeerService{}

// BeerService is the HTTP implementation of zymurgauge.BeerService
type BeerService struct {
	url    url.URL
	logger *logrus.Logger
}

// Get returns a beer by id
func (s *BeerService) Get(id uint64) (*zymurgauge.Beer, error) {

	u := s.url
	u.Path = "/api/beers/" + url.QueryEscape(strconv.FormatUint(id, 10))

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GET Beer %d", id)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			s.logger.Error(err)
		}
	}()

	var respBody getBeerResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, errors.Wrapf(err, "Could not decode Beer %d", id)
	} else if respBody.Err != "" {
		return nil, zymurgauge.Error(respBody.Err)
	}

	return respBody.Beer, nil
}

// Save creates or updates the stored beer with the given Beer
func (s *BeerService) Save(b *zymurgauge.Beer) error {

	if b == nil {
		return zymurgauge.ErrBeerRequired
	}

	u := s.url
	u.Path = "/api/beers"

	reqBody, err := json.Marshal(postBeerRequest{Beer: b})
	if err != nil {
		return errors.Wrapf(err, "Could not marshal Beer %d", b.ID)
	}

	resp, err := http.Post(u.String(), "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return errors.Wrapf(err, "Could not POST Beer %d", b.ID)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			s.logger.Error(err)
		}
	}()

	var respBody postBeerResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return errors.Wrapf(err, "Could not decode Beer %d", b.ID)
	} else if respBody.Err != "" {
		return zymurgauge.Error(respBody.Err)
	}

	*b = *respBody.Beer

	return nil
}
