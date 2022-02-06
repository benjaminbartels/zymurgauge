package brewfather

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/benjaminbartels/zymurgauge/internal/brewfather/model"
	"github.com/pkg/errors"
)

const (
	APIURL      = "https://api.brewfather.app/v1"
	batchesPath = "batches"
)

var (
	ErrUserAccessDenied = errors.New("access deined")
	ErrNotFound         = errors.New("resource not found")
	ErrTooManyRequests  = errors.New("too many request")
)

// _ model.Repo = (*Client)(nil).
var _ Service = (*ServiceClient)(nil)

type ServiceClient struct {
	client *http.Client
	logURL string
}

func New(userID, apiKey string, options ...OptionsFunc) *ServiceClient {
	t := &transport{
		userID: userID,
		apiKey: apiKey,
	}

	c := &ServiceClient{
		client: &http.Client{Transport: t},
	}

	for _, option := range options {
		option(c)
	}

	return c
}

type OptionsFunc func(*ServiceClient)

func SetTiltURL(url string) OptionsFunc {
	return func(c *ServiceClient) {
		c.logURL = url
	}
}

func (s *ServiceClient) GetAllSummaries(ctx context.Context) ([]BatchSummary, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", APIURL, batchesPath), nil)
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

	var batches []model.BatchSummary
	if err = json.NewDecoder(resp.Body).Decode(&batches); err != nil {
		return nil, errors.Wrap(err, "could not decode Batches")
	}

	return convertBatchSummaries(batches), nil
}

func (s *ServiceClient) GetDetail(ctx context.Context, id string) (*BatchDetail, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s/%s", APIURL, batchesPath, id), nil)
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

	var batch model.BatchDetail
	if err = json.NewDecoder(resp.Body).Decode(&batch); err != nil {
		return nil, errors.Wrap(err, "could not decode Batch")
	}

	b := convertBatchDetail(batch)

	return &b, nil
}

func (s *ServiceClient) Log(ctx context.Context, log LogEntry) error {
	data, err := json.Marshal(log)
	if err != nil {
		return errors.Wrap(err, "could not marshal TiltLogEntry")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.logURL, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrap(err, "could not create POST request for TiltLogEntry")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "could not POST TiltLogEntry")
	}

	defer resp.Body.Close()

	if err := parseStatusCode(resp.StatusCode); err != nil {
		return err
	}

	return nil
}

type transport struct {
	userID string
	apiKey string
}

func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	req := r.Clone(r.Context())
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Basic "+basicAuth(t.userID, t.apiKey))

	res, err := http.DefaultTransport.RoundTrip(req)
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

func convertBatchSummaries(batches []model.BatchSummary) []BatchSummary {
	s := []BatchSummary{}
	for i := 0; i < len(batches); i++ {
		s = append(s, convertBatchSummary(batches[i]))
	}

	return s
}

func convertBatchSummary(b model.BatchSummary) BatchSummary {
	return BatchSummary{
		ID:         b.ID,
		Name:       b.Name,
		Number:     b.BatchNo,
		RecipeName: b.Recipe.Name,
	}
}

func convertBatchDetail(b model.BatchDetail) BatchDetail {
	return BatchDetail{
		ID:           b.ID,
		Name:         b.Name,
		Number:       b.BatchNo,
		RecipeName:   b.Recipe.Name,
		Fermentation: convertFermentation(b.Recipe.Fermentation),
	}
}

func convertFermentation(fermentation model.Fermentation) Fermentation {
	return Fermentation{
		Name:  fermentation.Name,
		Steps: convertFermentationSteps(fermentation.Steps),
	}
}

func convertFermentationSteps(steps []model.FermentationStep) []FermentationStep {
	s := []FermentationStep{}
	for i := 0; i < len(steps); i++ {
		s = append(s, convertFermentationStep(steps[i]))
	}

	return s
}

func convertFermentationStep(step model.FermentationStep) FermentationStep {
	return FermentationStep{
		Type:            step.Type,
		ActualTime:      step.ActualTime,
		StepTemperature: step.StepTemp,
		StepTime:        step.StepTime,
	}
}
