package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/platform/app"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// FermentationResource is a client side rest resource used to manage Fermentations
type FermentationResource struct {
	url    *url.URL
	logger *logrus.Logger
}

func newFermentationResource(base string, version string) (*FermentationResource, error) {

	u, err := url.Parse(base + "/" + version + "/fermentations/")
	if err != nil {
		return nil, err
	}

	return &FermentationResource{url: u}, nil
}

// Get returns a fermentation by id
func (r *FermentationResource) Get(id uint64) (*internal.Fermentation, error) {

	resp, err := http.Get(r.url.String() + url.QueryEscape(strconv.FormatUint(id, 10)))
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GET Fermentation %d", id)
	}

	defer safeClose(resp.Body, &err)

	if resp.StatusCode == http.StatusNotFound {
		return nil, app.ErrNotFound
	}

	var fermentation *internal.Fermentation
	if err = json.NewDecoder(resp.Body).Decode(&fermentation); err != nil {
		return nil, errors.Wrapf(err, "Could not decode Beer %d", id)
	}

	return fermentation, nil
}

// Save creates or updates the stored fermentation with the given Fermentation
func (r *FermentationResource) Save(f *internal.Fermentation) error {

	reqBody, err := json.Marshal(f)
	if err != nil {
		return errors.Wrapf(err, "Could not marshal Fermentation %d", f.ID)
	}

	resp, err := http.Post(r.url.String(), "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return errors.Wrapf(err, "Could not POST Fermentation %d", f.ID)
	}

	defer safeClose(resp.Body, &err)

	if err := json.NewDecoder(resp.Body).Decode(&f); err != nil {
		return errors.Wrapf(err, "Could not decode Fermentation %s", f.ID)
	}

	return nil
}
