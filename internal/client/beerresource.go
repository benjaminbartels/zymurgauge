package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/platform/app"
	"github.com/benjaminbartels/zymurgauge/internal/platform/safeclose"
	"github.com/pkg/errors"
)

// BeerResource is a client side rest resource used to manage Beers
type BeerResource struct {
	url *url.URL
}

func newBeerResource(base string, version string) (*BeerResource, error) {
	u, err := url.Parse(base + "/" + version + "/beers/")
	if err != nil {
		return nil, err
	}
	return &BeerResource{url: u}, nil
}

// Get returns a beer by id
func (r *BeerResource) Get(id uint64) (*internal.Beer, error) {

	resp, err := http.Get(r.url.String() + url.QueryEscape(strconv.FormatUint(id, 10)))
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GET Beer %d", id)
	}

	defer safeclose.Close(resp.Body, &err)

	if resp.StatusCode == http.StatusNotFound {
		return nil, app.ErrNotFound
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

	resp, err := http.Post(r.url.String(), "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return errors.Wrapf(err, "Could not POST Beer %d", b.ID)
	}

	defer safeclose.Close(resp.Body, &err)

	if err := json.NewDecoder(resp.Body).Decode(&b); err != nil {
		return errors.Wrapf(err, "Could not decode Beer %s", b.ID)
	}

	return nil
}
