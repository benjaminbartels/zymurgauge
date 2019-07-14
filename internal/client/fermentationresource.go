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

// FermentationResource is a client side rest resource used to manage Fermentations
type FermentationResource struct {
	url   *url.URL
	token string
}

func newFermentationResource(base, version, token string) (*FermentationResource, error) {

	u, err := url.Parse(base + "/" + version + "/fermentations/")
	if err != nil {
		return nil, errors.Wrapf(err, "Could not create new FermentationResource")
	}

	return &FermentationResource{url: u, token: token}, nil
}

// Get returns a fermentation by id
func (r *FermentationResource) Get(id string) (*internal.Fermentation, error) {

	req, err := http.NewRequest(http.MethodGet, r.url.String()+url.QueryEscape(id), nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not create GET request for Fermentation %s", id)
	}
	req.Header.Add("authorization", "Bearer "+r.token)

	resp, err := http.DefaultClient.Do(req) // ToDo: Don't use default client...
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GET Fermentation %s", id)
	}

	defer safeclose.Close(resp.Body, &err)

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.Wrapf(web.ErrNotFound, "Fermentation %s does not exist", id)
	}

	var fermentation *internal.Fermentation
	if err = json.NewDecoder(resp.Body).Decode(&fermentation); err != nil {
		return nil, errors.Wrapf(err, "Could not decode Fermentation %s", id)
	}

	return fermentation, nil
}

// Save creates or updates the stored fermentation with the given Fermentation
func (r *FermentationResource) Save(f *internal.Fermentation) error {

	reqBody, err := json.Marshal(f)
	if err != nil {
		return errors.Wrapf(err, "Could not marshal Fermentation %s", f.ID)
	}

	req, err := http.NewRequest(http.MethodPost, r.url.String(), bytes.NewReader(reqBody))
	if err != nil {
		return errors.Wrapf(err, "Could not create POST request for Fermentation %s", f.ID)
	}

	req.Header.Add("Authorization", "Bearer "+r.token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req) // ToDo: Don't use default client...
	if err != nil {
		return errors.Wrapf(err, "Could not POST Fermentation %s", f.ID)
	}

	defer safeclose.Close(resp.Body, &err)

	if err := json.NewDecoder(resp.Body).Decode(&f); err != nil {
		return errors.Wrapf(err, "Could not decode Fermentation %s", f.ID)
	}

	return nil
}

// SaveTemperatureChange creates or updates the stored temperature change
func (r *FermentationResource) SaveTemperatureChange(t *internal.TemperatureChange) error {

	reqBody, err := json.Marshal(t)
	if err != nil {
		return errors.Wrapf(err, "Could not marshal TemperatureChange %s", t.ID)
	}

	req, err := http.NewRequest(http.MethodPost,
		r.url.String()+url.QueryEscape(t.FermentationID)+"/temperaturechanges",
		bytes.NewReader(reqBody))
	if err != nil {
		return errors.Wrapf(err, "Could not POST TemperatureChange %s", t.ID)
	}

	req.Header.Add("Authorization", "Bearer "+r.token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req) // ToDo: Don't use default client...
	if err != nil {
		return errors.Wrapf(err, "Could not POST TemperatureChange %s", t.ID)
	}

	defer safeclose.Close(resp.Body, &err)

	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return errors.Wrapf(err, "Could not decode TemperatureChange %s", t.ID)
	}

	return nil
}
