package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

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
func (r *BeerResource) Get(id string) (*internal.Beer, error) {

	req, err := http.NewRequest(http.MethodGet, r.url.String()+url.QueryEscape(id), nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not create GET request for Beer %s", id)
	}

	req.Header.Add("authorization", "Bearer "+r.token)

	resp, err := http.DefaultClient.Do(req) // ToDo: Don't use default client...
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GET Beer %s", id)
	}

	defer safeclose.Close(resp.Body, &err)

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.Wrapf(web.ErrNotFound, "Beer %s does not exist", id)
	}

	var beer *internal.Beer
	if err = json.NewDecoder(resp.Body).Decode(&beer); err != nil {
		return nil, errors.Wrapf(err, "Could not decode Beer %s", id)
	}

	return beer, nil
}

// Save creates or updates the stored beer with the given Beer
func (r *BeerResource) Save(b *internal.Beer) error {

	reqBody, err := json.Marshal(b)
	if err != nil {
		return errors.Wrapf(err, "Could not marshal Beer %s", b.ID)
	}

	req, err := http.NewRequest(http.MethodPost, r.url.String(), bytes.NewReader(reqBody))
	if err != nil {
		return errors.Wrapf(err, "Could not create POST request for Beer %s", b.ID)
	}

	req.Header.Add("Authorization", "Bearer "+r.token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req) // ToDo: Don't use default client...
	if err != nil {
		return errors.Wrapf(err, "Could not POST Beer %s", b.ID)
	}

	defer safeclose.Close(resp.Body, &err)

	if err := json.NewDecoder(resp.Body).Decode(&b); err != nil {
		return errors.Wrapf(err, "Could not decode Beer %s", b.ID)
	}

	return nil
}
