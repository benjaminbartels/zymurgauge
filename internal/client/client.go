package client

import (
	"fmt"

	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
)

// Client provides resources used to manage entities via REST.
type Client struct {
	ChamberResource      *ChamberResource
	BeerResource         *BeerResource
	FermentationResource *FermentationResource
}

// NewClient creates a new instance of the HTTP client
func NewClient(url fmt.Stringer, version string, logger log.Logger) (*Client, error) {

	chamberResource, err := newChamberResource(url.String()+"/api", version, logger)
	if err != nil {
		return nil, err
	}

	beerResource, err := newBeerResource(url.String()+"/api", version)
	if err != nil {
		return nil, err
	}

	fermentationResource, err := newFermentationResource(url.String()+"/api", version)
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
