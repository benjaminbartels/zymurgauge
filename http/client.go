package http

import (
	"net/url"

	"github.com/orangesword/zymurgauge"
	"github.com/sirupsen/logrus"
)

// Client provides services used to communicate with the HTTP server.
type Client struct {
	url                 url.URL
	logger              *logrus.Logger
	chamberService      ChamberService
	fermentationService FermentationService
	beerService         BeerService
}

// NewClient creates a new instance of the HTTP client
func NewClient(url url.URL, logger *logrus.Logger) *Client { // ToDo: Why is url no set
	c := &Client{
		url:    url,
		logger: logger,
	}
	c.chamberService.url = c.url
	c.fermentationService.url = c.url
	c.beerService.url = c.url
	c.chamberService.logger = c.logger
	c.fermentationService.logger = c.logger
	c.beerService.logger = c.logger

	return c
}

// ChamberService returns the service used to manage Chambers
func (c *Client) ChamberService() zymurgauge.ChamberService {
	return &c.chamberService
}

// FermentationService returns the service used to manage Fermentations
func (c *Client) FermentationService() zymurgauge.FermentationService {
	return &c.fermentationService
}

// BeerService returns the service used to manage Beers
func (c *Client) BeerService() zymurgauge.BeerService {
	return &c.beerService
}
