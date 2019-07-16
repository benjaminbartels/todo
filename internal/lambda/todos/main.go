package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	awslambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/benjaminbartels/todo/internal"
	"github.com/benjaminbartels/todo/internal/database/dynamodb"
	"github.com/benjaminbartels/todo/internal/lambda"
)

type handler struct {
	repo *dynamodb.ToDoRepo
}

func (h *handler) handle(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	switch req.HTTPMethod {
	case "GET":
		return h.get(req)
	case "POST":
		return h.post(req)
	case "PUT":
		return h.put(req)
	case "DELETE":
		return h.delete(req)
	default:
		return lambda.CreateErrorResponse(http.StatusMethodNotAllowed, fmt.Sprintf("%s not allowed", req.HTTPMethod))
	}
}

func (h *handler) get(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if id, ok := req.PathParameters["id"]; ok {

		todo, err := h.repo.Get(id)
		if err != nil {
			return lambda.CreateErrorResponse(http.StatusInternalServerError, err.Error())
		}

		if todo == nil {
			return lambda.CreateErrorResponse(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		}

		return lambda.CreateResponse(todo, http.StatusOK)

	}

	todos, err := h.repo.GetAll()
	if err != nil {
		return lambda.CreateErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return lambda.CreateResponse(todos, http.StatusOK)

}

func (h *handler) post(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	todo, err := parseToDo(req.Body)
	if err != nil {
		return lambda.CreateErrorResponse(http.StatusInternalServerError, err.Error())
	}

	if todo.ID != "" {
		return lambda.CreateErrorResponse(http.StatusBadRequest, "ID is required")
	}

	err = h.repo.Save(&todo)
	if err != nil {
		return lambda.CreateErrorResponse(http.StatusInternalServerError, err.Error())
	}
	return lambda.CreateResponse(todo, http.StatusOK)
}

func (h *handler) put(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	id, ok := req.PathParameters["id"]
	if !ok {
		return lambda.CreateErrorResponse(http.StatusBadRequest, "ID is required")
	}

	t, err := h.repo.Get(id)
	if err != nil {
		return lambda.CreateErrorResponse(http.StatusInternalServerError, err.Error())
	}

	if t == nil {
		return lambda.CreateErrorResponse(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	todo, err := parseToDo(req.Body)
	if err != nil {
		return lambda.CreateErrorResponse(http.StatusInternalServerError, err.Error())
	}

	if id != todo.ID {
		return lambda.CreateErrorResponse(http.StatusBadRequest, "ID in body does not match ID in path")
	}

	err = h.repo.Save(&todo)
	if err != nil {
		return lambda.CreateErrorResponse(http.StatusInternalServerError, err.Error())
	}
	return lambda.CreateResponse(todo, http.StatusOK)
}

func (h *handler) delete(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	id, ok := req.PathParameters["id"]

	if !ok {
		return lambda.CreateErrorResponse(http.StatusBadRequest, "ID is required")
	}

	t, err := h.repo.Get(id)
	if err != nil {
		return lambda.CreateErrorResponse(http.StatusInternalServerError, err.Error())
	}

	if t == nil {
		return lambda.CreateErrorResponse(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	if err := h.repo.Delete(id); err != nil {
		return lambda.CreateErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return lambda.CreateResponse(nil, http.StatusOK)

}

func parseToDo(body string) (internal.ToDo, error) {
	var t internal.ToDo
	err := json.Unmarshal([]byte(body), &t)
	return t, err
}

func main() {

	s, err := session.NewSession(aws.NewConfig().WithRegion("us-west-2"))
	if err != nil {
		panic(err) // ToDo: Is this ok?
	}

	db := awsdynamodb.New(s)
	repo := dynamodb.NewToDoRepo(db)

	h := handler{repo: repo}

	awslambda.Start(h.handle)
}
