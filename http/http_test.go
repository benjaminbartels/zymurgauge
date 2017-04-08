package http_test

import (
	"fmt"
	"net/url"

	"github.com/orangesword/zymurgauge/http"
	"github.com/orangesword/zymurgauge/mock"
	"github.com/sirupsen/logrus"
)

// TestService is a wrapper around the http.Service. It also includes references to the mocked services to allow
// for easier access to the mocked functions in the tests
type TestService struct {
	*http.Server
	BeerServiceMock         *mock.BeerService
	FermentationServiceMock *mock.FermentationService
	ChamberServiceMock      *mock.ChamberService
}

// MustOpenServerAndClient returns a new TestService with an open underlying Service as well as a http.Client
func MustOpenServerAndClient() (*TestService, *http.Client) {

	l := logrus.New()

	s := &TestService{
		Server: http.NewServer(),
	}
	s.Addr = ":0"

	// create an http handler, this will also instantate the router
	s.Handler = &http.Handler{
		BeerHandler:         http.NewBeerHandler(l),
		FermentationHandler: http.NewFermentationHandler(l),
		ChamberHandler:      http.NewChamberHandler(l),
	}

	// instantate mocks and assign them to service interfaces
	s.BeerServiceMock = &mock.BeerService{}
	s.Handler.BeerHandler.BeerService = s.BeerServiceMock
	s.FermentationServiceMock = &mock.FermentationService{}
	s.Handler.FermentationHandler.FermentationService = s.FermentationServiceMock
	s.ChamberServiceMock = &mock.ChamberService{}
	s.Handler.ChamberHandler.ChamberService = s.ChamberServiceMock

	if err := s.Open(); err != nil {
		panic(err)
	}

	c := http.NewClient()
	c.URL = url.URL{Scheme: "http", Host: fmt.Sprintf("localhost:%d", s.Port())}

	return s, c
}
