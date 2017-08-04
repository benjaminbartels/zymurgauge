package client

import (
	"net/url"
)

// Client provides services used to communicate with the HTTP server.
type Client struct {
	url                  url.URL
	chamberResource      *ChamberResource
	beerResource         *BeerResource
	fermentationResource *FermentationResource
	version              string
}

// NewClient creates a new instance of the HTTP client
func NewClient(url url.URL, version string) (*Client, error) { // ToDo: Why is url no set

	chamberResource, err := newChamberResource(url.String(), version)
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
		chamberResource:      chamberResource,
		beerResource:         beerResource,
		fermentationResource: fermentationResource,
	}

	return c, nil
}

func (c Client) ChamberResource() *ChamberResource {
	return c.chamberResource
}

func (c Client) BeerResource() *BeerResource {
	return c.beerResource
}

func (c Client) FermentationResource() *FermentationResource {
	return c.fermentationResource
}
