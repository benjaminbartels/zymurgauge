package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/platform/safeclose"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/pkg/errors"
)

// BeerResource is a client side rest resource used to manage Beers
type BeerResource struct {
	url   *url.URL
	token string
}

func newBeerResource(base, version, token string) (*BeerResource, error) {
	u, err := url.Parse(base + "/" + version + "/beers/")
	if err != nil {
		return nil, errors.Wrapf(err, "Could not create new BeerResource")
	}
	return &BeerResource{url: u, token: token}, nil
}

// Get returns a beer by id
func (r *BeerResource) Get(id uint64) (*internal.Beer, error) {

	req, err := http.NewRequest(http.MethodGet, r.url.String()+url.QueryEscape(strconv.FormatUint(id, 10)), nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not create GET request for Beer %d", id)
	}

	req.Header.Add("authorization", "Bearer "+r.token)

	resp, err := http.DefaultClient.Do(req) // ToDo: Dont use default client...
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GET Beer %d", id)
	}

	defer safeclose.Close(resp.Body, &err)

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.Wrapf(web.ErrNotFound, "Could not GET Beer %d", id)
	}

	var beer *internal.Beer
	if err = json.NewDecoder(resp.Body).Decode(&beer); err != nil {
		return nil, errors.Wrapf(err, "Could not decode Beer %d", id)
	}

	return beer, nil
}

// Save creates or updates the stored beer with the given Beer
func (r *BeerResource) Save(b *internal.Beer) error {

	reqBody, err := json.Marshal(b)
	if err != nil {
		return errors.Wrapf(err, "Could not marshal Beer %d", b.ID)
	}

	req, err := http.NewRequest(http.MethodPost, r.url.String(), bytes.NewReader(reqBody))
	if err != nil {
		return errors.Wrapf(err, "Could not create POST request for Beer %d", b.ID)
	}

	req.Header.Add("Authorization", "Bearer "+r.token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req) // ToDo: Dont use default client...
	if err != nil {
		return errors.Wrapf(err, "Could not POST Beer %d", b.ID)
	}

	defer safeclose.Close(resp.Body, &err)

	if err := json.NewDecoder(resp.Body).Decode(&b); err != nil {
		return errors.Wrapf(err, "Could not decode Beer %d", b.ID)
	}

	return nil
}
