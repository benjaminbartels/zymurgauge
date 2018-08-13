package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"github.com/benjaminbartels/zymurgauge/internal/platform/safeclose"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/pkg/errors"
)

var headerData = []byte("data:")

// ChamberResource is a client side rest resource used to manage Chambers
type ChamberResource struct {
	url    *url.URL
	stream *stream
	logger log.Logger
}

func newChamberResource(base string, version string, logger log.Logger) (*ChamberResource, error) {

	u, err := url.Parse(base + "/" + version + "/chambers/")
	if err != nil {
		return nil, errors.Wrapf(err, "Could not create new ChamberResource")
	}

	return &ChamberResource{url: u, logger: logger}, nil
}

// Get returns a controller by MAC address
func (r ChamberResource) Get(mac string) (*internal.Chamber, error) {

	resp, err := http.Get(r.url.String() + url.QueryEscape(mac))
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GET Chamber %s", mac)
	}

	defer safeclose.Close(resp.Body, &err)

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.Wrapf(web.ErrNotFound, "Could not GET Chamber %s", mac)
	}

	var chamber *internal.Chamber
	if err = json.NewDecoder(resp.Body).Decode(&chamber); err != nil {
		return nil, errors.Wrapf(err, "Could not decode Chamber %s", mac)
	}

	return chamber, nil
}

// Save creates or updates the stored controller with the given Chamber
func (r ChamberResource) Save(c *internal.Chamber) error {

	reqBody, err := json.Marshal(c)
	if err != nil {
		return errors.Wrapf(err, "Could not marshal Chamber %s", c.MacAddress)
	}

	resp, err := http.Post(r.url.String(), "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return errors.Wrapf(err, "Could not POST Chamber %s", c.MacAddress)
	}

	defer safeclose.Close(resp.Body, &err)

	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return errors.Wrapf(err, "Could not decode Chamber %s", c.MacAddress)
	}

	return nil
}

// Subscribe registers the caller to receives updates to the given controller on the given channel
func (r ChamberResource) Subscribe(mac string, ch chan internal.Chamber) error {

	r.url.Path = r.url.Path + url.QueryEscape(mac) + "/events"

	req, err := http.NewRequest("GET", r.url.String(), nil)
	if err != nil {
		return errors.Wrapf(err, "Could not GET Chamber events for %s", mac)
	}

	r.stream = &stream{
		ch:     ch,
		logger: r.logger,
	}

	err = r.stream.open(req)
	if err != nil {
		return errors.Wrapf(err, "Could not GET Chamber events for %s", mac)
	}
	return nil
}

// Unsubscribe unregisters the caller to receives updates to the given controller
func (r ChamberResource) Unsubscribe(mac string) { // ToDo: Does Unsubscribe need to return an error?
	r.stream.close()
}

// stream represents a http event stream
type stream struct {
	ch     chan internal.Chamber
	client *http.Client // ToDo: make all methods use this client
	resp   *http.Response
	logger log.Logger
}

// open opens the http event stream using the provied http request
func (s *stream) open(req *http.Request) error {

	// To Do: Does error handling need to be improved?
	req.Header.Set("Accept", "text/event-stream")

	s.client = &http.Client{}

	go func() {

		var err error
		s.resp, err = s.client.Do(req)
		if err != nil {
			s.logger.Println(err)
		}

		scanner := bufio.NewScanner(s.resp.Body)

		defer safeclose.Close(s.resp.Body, &err)

		for {
			if !scanner.Scan() {
				break
			}

			msg := scanner.Bytes()
			if err != nil {
				s.logger.Print(err)
				continue
			}

			s.logger.Printf("Received: [%s]\n", string(msg))

			if bytes.Contains(msg, headerData) {

				data := trimHeader(headerData, msg)

				var c = &internal.Chamber{}

				err = json.Unmarshal(data, c)
				if err != nil {
					s.logger.Print(err)
					continue
				}

				s.ch <- *c

			} else {
				s.logger.Printf("Unrecognized Message: %s\n", msg)
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
		s.logger.Print(err)
	} // To Do: Does error need to be returned?
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
