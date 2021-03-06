package brewfather

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const (
	APIURL      = "https://api.brewfather.app/v1"
	recipesPath = "recipes"
)

var (
	ErrUserAccessDenied = errors.New("access deined")
	ErrNotFound         = errors.New("resource not found")
	ErrTooManyRequests  = errors.New("too many request")
)

type Service struct {
	client  *http.Client
	baseURL string
}

type transport struct {
	userID string
	apiKey string
}

func New(baseURL, userID, apiKey string) *Service {
	t := transport{
		userID: userID,
		apiKey: apiKey,
	}

	return &Service{
		client:  &http.Client{Transport: &t},
		baseURL: baseURL,
	}
}

func (s *Service) GetRecipes(ctx context.Context) ([]Recipe, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", s.baseURL, recipesPath), nil)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create GET request for Recipes")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Could not GET Recipes")
	}

	defer resp.Body.Close()

	if err := parseStatusCode(resp.StatusCode); err != nil {
		return nil, err
	}

	var recipes []Recipe
	if err = json.NewDecoder(resp.Body).Decode(&recipes); err != nil {
		return nil, errors.Wrap(err, "Could not decode Recipes")
	}

	return recipes, nil
}

func (s *Service) GetRecipe(ctx context.Context, id string) (*Recipe, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s/%s", s.baseURL, recipesPath, id), nil)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create GET request for Recipes")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Could not GET Recipe")
	}

	defer resp.Body.Close()

	if err := parseStatusCode(resp.StatusCode); err != nil {
		return nil, err
	}

	var recipe *Recipe
	if err = json.NewDecoder(resp.Body).Decode(&recipe); err != nil {
		return nil, errors.Wrap(err, "Could not decode Recipe")
	}

	return recipe, nil
}

func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	req := r.Clone(r.Context())
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Basic "+basicAuth(t.userID, t.apiKey))

	return http.DefaultTransport.RoundTrip(req)
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
