package lambda

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
)

func CreateResponse(data interface{}, code int) (events.APIGatewayProxyResponse, error) {

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

func CreateErrorResponse(err error) (events.APIGatewayProxyResponse, error) {

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

	return CreateResponse(errorResponse{Err: err.Error()}, code)
}

// errorResponse is the response sent to the client in the event of a error
type errorResponse struct {
	Err string `json:"error,omitempty"`
}
