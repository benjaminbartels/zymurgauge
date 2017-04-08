package http

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/orangesword/zymurgauge"
	"github.com/orangesword/zymurgauge/gpio"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var headerData = []byte("data:")

// check implementation at compile time
var _ zymurgauge.ChamberService = &ChamberService{}

// ChamberService is the HTTP implementation of zymurgauge.ChamberService
type ChamberService struct {
	url    url.URL
	logger *logrus.Logger
	stream *stream
}

// Get returns a controller by MAC address
func (s *ChamberService) Get(mac string) (*zymurgauge.Chamber, error) {

	u := s.url
	u.Path = "/api/chambers/" + url.QueryEscape(mac)

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GET Chamber %s", mac)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			s.logger.Error(err)
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return nil, zymurgauge.ErrNotFound
	}

	var respBody getChamberResponse

	respBody.Chamber = &zymurgauge.Chamber{Controller: &gpio.Thermostat{}}

	// buf := bytes.NewBuffer(make([]byte, 0, resp.ContentLength))
	// buf.ReadFrom(resp.Body)
	// fmt.Printf("Response: %s\n", string(buf.Bytes()))

	if err = json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, errors.Wrapf(err, "Could not decode Chamber %s", mac)
	} else if respBody.Err != "" {

		return nil, zymurgauge.Error(respBody.Err)
	}

	return respBody.Chamber, nil
}

// Save creates or updates the stored controller with the given Chamber
func (s *ChamberService) Save(f *zymurgauge.Chamber) error {

	if f == nil {
		return zymurgauge.ErrChamberRequired
	} else if f.MacAddress == "" {
		return zymurgauge.ErrMacAddressRequired
	}

	u := s.url
	u.Path = "/api/chambers"

	reqBody, err := json.Marshal(postChamberRequest{Chamber: f})
	if err != nil {
		return errors.Wrapf(err, "Could not marshal Chamber %s", f.MacAddress)
	}

	resp, err := http.Post(u.String(), "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return errors.Wrapf(err, "Could not POST Chamber %s", f.MacAddress)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			s.logger.Error(err)
		}
	}()

	var respBody postChamberResponse
	respBody.Chamber = &zymurgauge.Chamber{Controller: &gpio.Thermostat{}}

	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return errors.Wrapf(err, "Could not decode Chamber %s", f.MacAddress)
	} else if respBody.Err != "" {
		return zymurgauge.Error(respBody.Err)
	}

	*f = *respBody.Chamber

	return nil
}

// Subscribe registers the caller to receives updates to the given controller on the given channel
func (s *ChamberService) Subscribe(mac string, ch chan zymurgauge.Chamber) error {

	if ch == nil {
		return zymurgauge.ErrChamberRequired // ToDo: Add new error type
	}

	if s.stream != nil {
		return zymurgauge.ErrChamberRequired // ToDo: Add new error type
	}

	u := s.url
	u.Path = "/api/chambers/" + url.QueryEscape(mac) + "/events"

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return errors.Wrapf(err, "Could not GET Chamber events for %s", mac)
	}

	s.stream = &stream{
		ch:     ch,
		logger: s.logger,
	}

	return s.stream.open(req)
}

// Unsubscribe unregisters the caller to receives updates to the given controller
func (s *ChamberService) Unsubscribe(mac string) {
	s.stream.close()
	s.stream = nil
}

// stream represents a http event stream
type stream struct {
	ch     chan zymurgauge.Chamber
	client *http.Client // ToDo: make all methods use this client
	resp   *http.Response
	logger *logrus.Logger
}

// open opens the http event stream using the provied http request
func (s *stream) open(req *http.Request) error {
	req.Header.Set("Accept", "text/event-stream") // ToDo: Revisit maybe

	s.client = &http.Client{}

	go func() {

		var err error
		s.resp, err = s.client.Do(req)
		if err != nil {
			s.logger.Error(err)
		}

		scanner := bufio.NewScanner(s.resp.Body)

		defer func() {
			err = s.resp.Body.Close()
			if err != nil {
				s.logger.Error(err)
			}
		}()

		for {
			if !scanner.Scan() {
				break
			}

			msg := scanner.Bytes()
			if err != nil {
				s.logger.Error(err)
				continue
			}

			s.logger.Debugf("Received: [%s]\n", string(msg))

			if bytes.Contains(msg, headerData) {

				data := trimHeader(headerData, msg)

				var c = &zymurgauge.Chamber{}
				c.Controller = &gpio.Thermostat{}

				err = json.Unmarshal(data, c)
				if err != nil {
					s.logger.Error(err)
					continue
				}

				s.ch <- *c

			} else {
				s.logger.Infof("Unrecognized Message: %s", msg)
			}
		}
	}()

	return nil
}

// close closes the http event stream
func (s *stream) close() {
	close(s.ch)
	err := s.resp.Body.Close()
	if err != nil {
		s.logger.Error(err)
	}
}

// trimHeader remove the header label from the provided byte array
func trimHeader(h []byte, data []byte) []byte {
	data = data[len(h):]
	// Remove optional leading whitespace
	if data[0] == 32 {
		data = data[1:]
	}
	// Remove trailing new line
	if data[len(data)-1] == 10 {
		data = data[:len(data)-1]
	}
	return data
}
