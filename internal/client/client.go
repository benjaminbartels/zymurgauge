package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"github.com/benjaminbartels/zymurgauge/internal/storage"
	"github.com/pkg/errors"
)

const (
	authURL   = "https://zymurgauge.auth0.com/oauth/token"
	audience  = "zymurgauge.com/api"
	grantType = "client_credentials"
)

// Client provides resources used to manage entities via REST.
type Client struct {
	ChamberProvider      ChamberProvider
	BeerResource         *BeerResource
	FermentationProvider FermentationProvider
}

// NewClient creates a new instance of the HTTP client.
func NewClient(url fmt.Stringer, version, clientID, clientSecret string, logger log.Logger) (*Client, error) {
	tok, err := authenticate(clientID, clientSecret)
	if err != nil {
		return nil, err
	}

	chamberResource, err := newChamberResource(url.String()+"/api", version, tok, logger)
	if err != nil {
		return nil, err
	}

	beerResource, err := newBeerResource(url.String()+"/api", version, tok)
	if err != nil {
		return nil, err
	}

	fermentationResource, err := newFermentationResource(url.String()+"/api", version, tok)
	if err != nil {
		return nil, err
	}

	c := &Client{
		ChamberProvider:      chamberResource,
		BeerResource:         beerResource,
		FermentationProvider: fermentationResource,
	}

	return c, nil
}

func authenticate(clientID, clientSecret string) (string, error) {
	tokReq := getTokenRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Audience:     audience,
		GrantType:    grantType,
	}

	reqBody, err := json.Marshal(tokReq)
	if err != nil {
		return "", errors.Wrap(err, "Could not marshal getTokenRequest")
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, authURL, bytes.NewReader(reqBody))
	if err != nil {
		return "", errors.Wrap(err, "Could not create POST request")
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "Could not POST getTokenRequest")
	}

	defer resp.Body.Close()

	var t getTokenResponse

	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return "", errors.Wrap(err, "Could not decode getTokenResponse")
	}

	return t.AccessToken, nil
}

type getTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
	GrantType    string `json:"grant_type"`
}

type getTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// ChamberProvider is a data provider for Chambers.
type ChamberProvider interface {
	Get(mac string) (*storage.Chamber, error)
	Save(c *storage.Chamber) error
}

// FermentationProvider is a data provider for Fermentations.
type FermentationProvider interface {
	Get(id uint64) (*storage.Fermentation, error)
	Save(f *storage.Fermentation) error
	SaveTemperatureChange(t *internal.TemperatureChange) error
}
