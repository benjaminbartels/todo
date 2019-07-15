package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	awslambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/benjaminbartels/todo/internal"
	"github.com/benjaminbartels/todo/internal/database/dynamodb"
	"github.com/benjaminbartels/todo/internal/lambda"
	uuid "github.com/satori/go.uuid"
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
	// case "PUT":
	// 	return h.put(req)
	// case "DELETE":
	// 	return h.delete(req)
	default:
		return lambda.CreateErrorResponse(lambda.ErrMethodNotAllowed)
	}
}

func (h *handler) get(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	todos, err := h.repo.GetAll()
	if err != nil {
		return lambda.CreateErrorResponse(err)
	}
	return lambda.CreateResponse(todos, http.StatusOK)
}

func (h *handler) post(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	todo, err := parseToDo(req.Body)
	if err != nil {
		return lambda.CreateErrorResponse(err)
	}

	if todo.ID == "" {
		todo.ID = uuid.NewV4().String()
	}

	err = h.repo.Save(&todo)
	if err != nil {
		return lambda.CreateErrorResponse(err)
	}
	return lambda.CreateResponse(todo, http.StatusOK)
}

func (h *handler) put(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	todo, err := parseToDo(req.Body)
	if err != nil {
		return lambda.CreateErrorResponse(err)
	}

	err = h.repo.Save(&todo)
	if err != nil {
		return lambda.CreateErrorResponse(err)
	}
	return lambda.CreateResponse(todo, http.StatusOK)
}

func parseToDo(body string) (internal.ToDo, error) {
	var b internal.ToDo
	err := json.Unmarshal([]byte(body), b)
	return b, err
}

func main() {

	s, _ := session.NewSession(aws.NewConfig().WithRegion("us-west-2"))

	db := awsdynamodb.New(s)
	repo := dynamodb.NewToDoRepo(db)

	h := handler{repo: repo}

	awslambda.Start(h.handle)
}
