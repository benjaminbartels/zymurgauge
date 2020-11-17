package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/benjaminbartels/zymurgauge/internal/storage"
	"github.com/pkg/errors"
)

// ChamberResource is a client side rest resource used to manage Chambers.
type ChamberResource struct {
	url    *url.URL
	token  string
	logger log.Logger
}

func newChamberResource(base, version, token string, logger log.Logger) (*ChamberResource, error) {
	u, err := url.Parse(base + "/" + version + "/chambers/")
	if err != nil {
		return nil, errors.Wrapf(err, "Could not create new ChamberResource")
	}

	return &ChamberResource{url: u, token: token, logger: logger}, nil
}

// Get returns a controller by id.
func (r ChamberResource) Get(mac string) (*storage.Chamber, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, r.url.String()+url.QueryEscape(mac), nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not create GET request for Chamber %s", mac)
	}

	req.Header.Add("authorization", "Bearer "+r.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GET Chamber %s", mac)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.Wrapf(web.ErrNotFound, "Chamber %s does not exist", mac)
	}

	var chamber *storage.Chamber
	if err = json.NewDecoder(resp.Body).Decode(&chamber); err != nil {
		return nil, errors.Wrapf(err, "Could not decode Chamber %s", mac)
	}

	return chamber, nil
}

// Save creates or updates the stored controller with the given Chamber.
func (r ChamberResource) Save(c *storage.Chamber) error {
	reqBody, err := json.Marshal(c)
	if err != nil {
		return errors.Wrapf(err, "Could not marshal Chamber %s", c.MacAddress)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, r.url.String(), bytes.NewReader(reqBody))
	if err != nil {
		return errors.Wrapf(err, "Could not create POST request for Chamber %s", c.MacAddress)
	}

	req.Header.Add("Authorization", "Bearer "+r.token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "Could not POST Chamber %s", c.MacAddress)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return errors.Wrapf(err, "Could not decode Chamber %s", c.MacAddress)
	}

	return nil
}
