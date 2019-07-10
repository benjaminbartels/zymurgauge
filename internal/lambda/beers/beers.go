package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/database/dynamodb"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type handler struct {
	repo *dynamodb.BeerRepo
}

func (h *handler) handle(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return h.get(req)
	case "POST":
		return h.post(req)
	default:
		return createErrorResponse(ErrMethodNotAllowed)
	}
}

func (h *handler) get(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	beers, err := h.repo.GetAll()
	if err != nil {
		return createErrorResponse(err)
	}
	return createResponse(beers, http.StatusOK)
}

func (h *handler) post(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	beer, err := parseBeer(req.Body)
	if err != nil {
		return createErrorResponse(err)
	}

	if beer.ID == "" {
		beer.ID = uuid.NewV4().String()
	}

	err = h.repo.Save(&beer)
	if err != nil {
		return createErrorResponse(err)
	}
	return createResponse(beer, http.StatusOK)
}

func (h *handler) put(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	beer, err := parseBeer(req.Body)
	if err != nil {
		return createErrorResponse(err)
	}
	
	err = h.repo.Save(&beer)
	if err != nil {
		return createErrorResponse(err)
	}
	return createResponse(beer, http.StatusOK)
}


// func create(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
//     if req.Headers["Content-Type"] != "application/json" {
//         return clientError(http.StatusNotAcceptable)
//     }

//     bk := new(beer)
//     err := json.Unmarshal([]byte(req.Body), bk)
//     if err != nil {
//         return clientError(http.StatusUnprocessableEntity)
//     }

//     if !isbnRegexp.MatchString(bk.ISBN) {
//         return clientError(http.StatusBadRequest)
//     }
//     if bk.Title == "" || bk.Author == "" {
//         return clientError(http.StatusBadRequest)
//     }

//     err = putItem(bk)
//     if err != nil {
//         return serverError(err)
//     }

//     return events.APIGatewayProxyResponse{
//         StatusCode: 201,
//         Headers:    map[string]string{"Location": fmt.Sprintf("/beers?isbn=%s", bk.ISBN)},
//     }, nil
// }

func parseBeer(body string) (internal.Beer, error) {
	var b internal.Beer
    err := json.Unmarshal([]byte(body), b)
	return b, err
}

func createResponse(data interface{}, code int) (events.APIGatewayProxyResponse, error) {

	r := events.APIGatewayProxyResponse{
		StatusCode: code,
	}

	// No Content
	if data == nil {
		r.StatusCode = http.StatusNoContent
		data = errorResponse{Err: http.StatusText(http.StatusNoContent)}
	}

	// Marshal into a JSON
	js, err := json.Marshal(data)
	if err != nil {
		r.StatusCode = http.StatusInternalServerError
		js, err = json.Marshal(errorResponse{Err: err.Error()})
		if err != nil {
			return r, err
		}
	}

	r.Body = string(js)

	return r, err
}

var (
	// ErrNotFound is returned when an entity is not found
	ErrNotFound = errors.New("not found")
	// ErrInternal is returned when an internal error has occurred
	ErrInternal = errors.New("internal error")
	// ErrBadRequest is returned when the request is invalid
	ErrBadRequest = errors.New("bad request")
	// ErrMethodNotAllowed is returned when the request method (GET, POST, etc.) is not allowed
	ErrMethodNotAllowed = errors.New("method not allowed")
	// ErrUnauthorized is returned when the request is not authorized
	ErrUnauthorized = errors.New("unauthorized")
)

func createErrorResponse(err error) (events.APIGatewayProxyResponse, error) {

	var code int

	switch errors.Cause(err) {
	case ErrNotFound:
		code = http.StatusNotFound
	case ErrBadRequest: //ToDO: what was bad?
		code = http.StatusBadRequest
	case ErrMethodNotAllowed:
		code = http.StatusMethodNotAllowed
	case ErrUnauthorized:
		code = http.StatusUnauthorized
	default:
		code = http.StatusInternalServerError
	}

	return createResponse(errorResponse{Err: err.Error()}, code)
}

// errorResponse is the response sent to the client in the event of a error
type errorResponse struct {
	Err string `json:"error,omitempty"`
}

func main() {

	s, _ := session.NewSession(aws.NewConfig().WithRegion("us-west-2"))

	db := awsdynamodb.New(s)
	repo := dynamodb.NewBeerRepo(db)

	h := handler{repo: repo}

	lambda.Start(h.handle)
}
