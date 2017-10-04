package client

import (
	"io"
	"net/url"

	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
)

// Client provides resources used to manage entities via REST.
type Client struct {
	ChamberResource      *ChamberResource
	BeerResource         *BeerResource
	FermentationResource *FermentationResource
}

// NewClient creates a new instance of the HTTP client
func NewClient(url url.URL, version string, logger log.Logger) (*Client, error) { // ToDo: Why is url no set

	chamberResource, err := newChamberResource(url.String(), version, logger)
	if err != nil {
		return nil, err
	}

	beerResource, err := newBeerResource(url.String(), version)
	if err != nil {
		return nil, err
	}

	fermentationResource, err := newFermentationResource(url.String(), version)
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

func safeClose(c io.Closer, err *error) {
	if cerr := c.Close(); cerr != nil && *err == nil {
		*err = cerr
	}
}
