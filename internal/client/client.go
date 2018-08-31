package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"github.com/benjaminbartels/zymurgauge/internal/platform/safeclose"
	"github.com/pkg/errors"
)

const (
	authURL   = "https://zymurgauge.auth0.com/oauth/token"
	audience  = "zymurgauge.com/api"
	grantType = "client_credentials"
)

// Client provides resources used to manage entities via REST.
type Client struct {
	ChamberResource      *ChamberResource
	BeerResource         *BeerResource
	FermentationResource *FermentationResource
}

// NewClient creates a new instance of the HTTP client
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
		ChamberResource:      chamberResource,
		BeerResource:         beerResource,
		FermentationResource: fermentationResource,
	}

	return c, nil
}

func authenticate(clientID, clientSecret string) (string, error) {

	req := getTokenRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Audience:     audience,
		GrantType:    grantType,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", errors.Wrap(err, "Could not marshal getTokenRequest")
	}

	resp, err := http.Post(authURL, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return "", errors.Wrap(err, "Could not POST getTokenRequest")
	}

	defer safeclose.Close(resp.Body, &err)

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
