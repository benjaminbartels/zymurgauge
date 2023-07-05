package brewfather

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const (
	apiURL              = "https://api.brewfather.app/v2"
	batchesPath         = "batches"
	logInterval         = 15 * time.Minute
	requestTimeout      = 10 * time.Second
	dialTimeout         = 10 * time.Second
	tlsHandshakeTimeout = 10 * time.Second
)

var (
	ErrUserAccessDenied = errors.New("access deined")
	ErrNotFound         = errors.New("resource not found")
	ErrTooManyRequests  = errors.New("too many request")
	ErrLogURLNotSet     = errors.New("log url is not set")
)

// _ model.Repo = (*Client)(nil).
var _ Service = (*ServiceClient)(nil)

type ServiceClient struct {
	client       *http.Client
	userID       string
	apiKey       string
	logURL       string
	nextSendTime time.Time
}

func New(userID, apiKey, logURL string) *ServiceClient {
	c := &ServiceClient{
		client:       createHTTPClient(userID, apiKey),
		logURL:       logURL,
		nextSendTime: time.Now(),
	}

	return c
}

func createHTTPClient(userID, apiKey string) *http.Client {
	t := &transport{
		userID: userID,
		apiKey: apiKey,
		roundTripper: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: dialTimeout,
			}).Dial,
			TLSHandshakeTimeout: tlsHandshakeTimeout,
		},
	}

	return &http.Client{
		Timeout:   requestTimeout,
		Transport: t,
	}
}

func (s *ServiceClient) GetAllBatchSummaries(ctx context.Context) ([]BatchSummary, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", apiURL, batchesPath), nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create GET request for Batches")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not GET Batches")
	}

	defer resp.Body.Close()

	if err := parseStatusCode(resp.StatusCode); err != nil {
		return nil, err
	}

	var batches []BatchSummary
	if err = json.NewDecoder(resp.Body).Decode(&batches); err != nil {
		return nil, errors.Wrap(err, "could not decode Batches")
	}

	return batches, nil
}

func (s *ServiceClient) GetBatchDetail(ctx context.Context, id string) (*BatchDetail, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s/%s", apiURL, batchesPath, id), nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create GET request for Batch")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not GET Batch")
	}

	defer resp.Body.Close()

	if err := parseStatusCode(resp.StatusCode); err != nil {
		return nil, err
	}

	var batch BatchDetail
	if err = json.NewDecoder(resp.Body).Decode(&batch); err != nil {
		return nil, errors.Wrap(err, "could not decode Batch")
	}

	return &batch, nil
}

func (s *ServiceClient) Log(ctx context.Context, log LogEntry) error {
	if s.logURL == "" {
		return ErrLogURLNotSet
	}

	if time.Now().After(s.nextSendTime) {
		data, err := json.Marshal(log)
		if err != nil {
			return errors.Wrap(err, "could not marshal Tilt log entry")
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.logURL, bytes.NewBuffer(data))
		if err != nil {
			return errors.Wrap(err, "could not create POST request for Tilt log entry")
		}

		resp, err := s.client.Do(req)
		if err != nil {
			return errors.Wrap(err, "could not POST Tilt log entry")
		}

		defer resp.Body.Close()

		if err := parseStatusCode(resp.StatusCode); err != nil {
			return err
		}

		s.nextSendTime = time.Now().Add(logInterval)
	}

	return nil
}

func (s *ServiceClient) UpdateSettings(userID, apiKey, tiltURL string) {
	if userID != s.userID || apiKey != s.apiKey {
		s.client = createHTTPClient(userID, apiKey)
	}

	s.logURL = tiltURL
}

type transport struct {
	userID       string
	apiKey       string
	roundTripper http.RoundTripper
}

func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	req := r.Clone(r.Context())
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Basic "+basicAuth(t.userID, t.apiKey))

	res, err := t.roundTripper.RoundTrip(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not perform RoundTrip")
	}

	return res, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password

	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func parseStatusCode(code int) error {
	switch code {
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusUnauthorized,
		http.StatusForbidden:
		return ErrUserAccessDenied
	case http.StatusTooManyRequests:
		return ErrTooManyRequests
	default:
		return nil
	}
}
