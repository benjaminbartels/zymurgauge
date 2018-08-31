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
func (r *FermentationResource) Get(id uint64) (*internal.Fermentation, error) {

	req, err := http.NewRequest(http.MethodGet, r.url.String()+url.QueryEscape(strconv.FormatUint(id, 10)), nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not create GET request for Fermentation %d", id)
	}

	req.Header.Add("authorization", "Bearer "+r.token)

	resp, err := http.DefaultClient.Do(req) // ToDo: Dont use default client...
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GET Fermentation %d", id)
	}

	defer safeclose.Close(resp.Body, &err)

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.Wrapf(web.ErrNotFound, "Fermentation %d does not exist", id)
	}

	var fermentation *internal.Fermentation
	if err = json.NewDecoder(resp.Body).Decode(&fermentation); err != nil {
		return nil, errors.Wrapf(err, "Could not decode Fermentation %d", id)
	}

	return fermentation, nil
}

// Save creates or updates the stored fermentation with the given Fermentation
func (r *FermentationResource) Save(f *internal.Fermentation) error {

	reqBody, err := json.Marshal(f)
	if err != nil {
		return errors.Wrapf(err, "Could not marshal Fermentation %d", f.ID)
	}

	req, err := http.NewRequest(http.MethodPost, r.url.String(), bytes.NewReader(reqBody))
	if err != nil {
		return errors.Wrapf(err, "Could not create POST request for Fermentation %d", f.ID)
	}

	req.Header.Add("Authorization", "Bearer "+r.token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req) // ToDo: Dont use default client...
	if err != nil {
		return errors.Wrapf(err, "Could not POST Fermentation %d", f.ID)
	}

	defer safeclose.Close(resp.Body, &err)

	if err := json.NewDecoder(resp.Body).Decode(&f); err != nil {
		return errors.Wrapf(err, "Could not decode Fermentation %d", f.ID)
	}

	return nil
}

// SaveTemperatureChange creates or updates the stored temperature change
func (r *FermentationResource) SaveTemperatureChange(t *internal.TemperatureChange) error {

	reqBody, err := json.Marshal(t)
	if err != nil {
		return errors.Wrapf(err, "Could not marshal TemperatureChange %d", t.ID)
	}

	req, err := http.NewRequest(http.MethodPost,
		r.url.String()+url.QueryEscape(strconv.FormatUint(t.FermentationID, 10))+"/temperaturechanges",
		bytes.NewReader(reqBody))
	if err != nil {
		return errors.Wrapf(err, "Could not POST TemperatureChange %d", t.ID)
	}

	req.Header.Add("Authorization", "Bearer "+r.token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req) // ToDo: Dont use default client...
	if err != nil {
		return errors.Wrapf(err, "Could not POST TemperatureChange %d", t.ID)
	}

	defer safeclose.Close(resp.Body, &err)

	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return errors.Wrapf(err, "Could not decode TemperatureChange %d", t.ID)
	}

	return nil
}
