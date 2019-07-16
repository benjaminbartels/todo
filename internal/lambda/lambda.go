package lambda

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
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

func CreateErrorResponse(code int, msg string) (events.APIGatewayProxyResponse, error) {
	return CreateResponse(errorResponse{Err: msg}, code)
}

// errorResponse is the response sent to the client in the event of a error
type errorResponse struct {
	Err string `json:"error,omitempty"`
}
