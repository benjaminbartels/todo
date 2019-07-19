package lambda

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// Create Response generates an APIGatewayProxyResponse using the provided data and http code
func CreateResponse(data interface{}, code int) (events.APIGatewayProxyResponse, error) {

	r := events.APIGatewayProxyResponse{
		StatusCode: code,
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

	r.Headers = make(map[string]string)
	r.Headers["Access-Control-Allow-Origin"] = "*"
	r.Headers["Access-Control-Allow-Credentials"] = "true"

	r.Body = string(js)

	return r, err
}

// Create Error Response generates an APIGatewayProxyResponse using the provided message and http code
func CreateErrorResponse(msg string, code int) (events.APIGatewayProxyResponse, error) {
	return CreateResponse(errorResponse{Err: msg}, code)
}

// errorResponse is the response sent to the client in the event of a error
type errorResponse struct {
	Err string `json:"error,omitempty"`
}
