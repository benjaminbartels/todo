package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/benjaminbartels/todo/internal"
	"github.com/benjaminbartels/todo/internal/lambda/handlers"
	"github.com/pkg/errors"
)

const testUUID = "a8a43435-20d8-4af2-8f94-f504aff2c6f3"

var newToDo = internal.ToDo{
	Title: "Some ToDo",
}

var savedToDo = internal.ToDo{
	ID:    testUUID,
	Title: "Some ToDo",
}

func TestToDoRepo(t *testing.T) {
	t.Run("GetToDoOK", testGetToDoOK)
	t.Run("GetToDoNotFound", testGetToDoNotFound)
	t.Run("GetToDoInternalError", testGetToDoInternalError)
	t.Run("GetAllToDoOK", testGetAllToDoOK)
	t.Run("GetAllToDoInternalError", testGetAllToDoInternalError)
	t.Run("CreateToDoOK", testCreateToDoOK)
	t.Run("CreateToDoBadRequest", testCreateToDoBadRequest)
	t.Run("CreateToDoInternalErrorOnParse", testCreateToDoInternalErrorOnParse)
	t.Run("CreateToDoInternalErrorOnSave", testCreateToDoInternalErrorOnSave)
	t.Run("UpdateToDoOK", testUpdateToDoOK)
	t.Run("UpdateToDoBadRequestMissingID", testUpdateToDoBadRequestMissingID)
	t.Run("UpdateToDoBadRequestNoMatch", testUpdateToDoBadRequestNoMatch)
	t.Run("UpdateToDoNotFound", testUpdateToDoNotFound)
	t.Run("UpdateToDoInternalErrorOnParse", testUpdateToDoInternalErrorOnParse)
	t.Run("UpdateToDoInternalErrorOnGet", testUpdateToDoInternalErrorOnGet)
	t.Run("UpdateToDoInternalErrorOnSave", testUpdateToDoInternalErrorOnSave)
	t.Run("DeleteToDoOK", testDeleteToDoOK)
	t.Run("DeleteToDoBadRequestMissingID", testDeleteToDoBadRequestMissingID)
	t.Run("DeleteToDoNotFound", testDeleteToDoNotFound)
	t.Run("DeleteToDoInternalErrorOnGet", testDeleteToDoInternalErrorOnGet)
	t.Run("DeleteToDoInternalErrorOnDelete", testDeleteToDoInternalErrorOnDelete)
	t.Run("MethodNotAllowed", testMethodNotAllowed)
}

func testGetToDoOK(t *testing.T) {

	m := &RepoMock{
		GetFn: func(string) (*internal.ToDo, error) {
			return &savedToDo, nil
		},
	}

	req := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": testUUID},
		HTTPMethod:     http.MethodGet,
	}

	resp, err := handlers.NewToDoHandler(m).Handle(req)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(resp.Body, testUUID) {
		t.Fatalf("Expected body to contain '%s'", testUUID)
	}

	if !m.GetInvoked {
		t.Fatal("Get not invoked")
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected %d http response code, got %d", http.StatusOK, resp.StatusCode)
	}

}

func testGetToDoNotFound(t *testing.T) {

	m := &RepoMock{
		GetFn: func(string) (*internal.ToDo, error) {
			return nil, nil
		},
	}

	req := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": testUUID},
		HTTPMethod:     http.MethodGet,
	}

	resp, err := handlers.NewToDoHandler(m).Handle(req)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(resp.Body, handlers.ErrNotFound.Error()) {
		t.Fatalf("Expected body to contain '%s'", handlers.ErrNotFound.Error())
	}

	if !m.GetInvoked {
		t.Fatal("Get not invoked")
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected %d http response code, got %d", http.StatusNotFound, resp.StatusCode)
	}

}

func testGetToDoInternalError(t *testing.T) {

	m := &RepoMock{
		GetFn: func(string) (*internal.ToDo, error) {
			return nil, errors.New("DB Error")
		},
	}

	req := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": testUUID},
		HTTPMethod:     http.MethodGet,
	}

	resp, err := handlers.NewToDoHandler(m).Handle(req)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(resp.Body, handlers.ErrInternal.Error()) {
		t.Fatalf("Expected body to contain '%s'", handlers.ErrInternal.Error())
	}

	if !m.GetInvoked {
		t.Fatal("Get not invoked")
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected %d http response code, got %d", http.StatusInternalServerError, resp.StatusCode)
	}

}

func testGetAllToDoOK(t *testing.T) {

	m := &RepoMock{
		GetAllFn: func() ([]internal.ToDo, error) {
			return []internal.ToDo{savedToDo}, nil
		},
	}

	req := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodGet,
	}

	resp, err := handlers.NewToDoHandler(m).Handle(req)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(resp.Body, testUUID) {
		t.Fatalf("Expected body to contain '%s'", testUUID)
	}

	if !m.GetAllInvoked {
		t.Fatal("GetAll not invoked")
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected %d http response code, got %d", http.StatusOK, resp.StatusCode)
	}

}

func testGetAllToDoInternalError(t *testing.T) {

	m := &RepoMock{
		GetAllFn: func() ([]internal.ToDo, error) {
			return []internal.ToDo{}, errors.New("DB Error")
		},
	}

	req := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodGet,
	}

	resp, err := handlers.NewToDoHandler(m).Handle(req)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(resp.Body, handlers.ErrInternal.Error()) {
		t.Fatalf("Expected body to contain '%s'", handlers.ErrInternal.Error())
	}

	if !m.GetAllInvoked {
		t.Fatal("GetAll not invoked")
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected %d http response code, got %d", http.StatusInternalServerError, resp.StatusCode)
	}

}

func testCreateToDoOK(t *testing.T) {

	m := &RepoMock{
		SaveFn: func(todo *internal.ToDo) error {
			todo.ID = testUUID
			return nil
		},
	}

	req := events.APIGatewayProxyRequest{
		Body:       toDoToString(&newToDo),
		HTTPMethod: http.MethodPost,
	}

	resp, err := handlers.NewToDoHandler(m).Handle(req)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(resp.Body, testUUID) {
		t.Fatalf("Expected body to contain '%s'", testUUID)
	}

	if !m.SaveInvoked {
		t.Fatal("Save not invoked")
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected %d http response code, got %d", http.StatusOK, resp.StatusCode)
	}

}

func testCreateToDoBadRequest(t *testing.T) {

	m := &RepoMock{
		SaveFn: func(*internal.ToDo) error {
			return nil
		},
	}

	req := events.APIGatewayProxyRequest{
		Body:       toDoToString(&savedToDo),
		HTTPMethod: http.MethodPost,
	}

	resp, err := handlers.NewToDoHandler(m).Handle(req)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(resp.Body, handlers.ErrBadRequest.Error()) {
		t.Fatalf("Expected body to contain '%s'", handlers.ErrBadRequest.Error())
	}

	if m.SaveInvoked {
		t.Fatal("Save invoked")
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected %d http response code, got %d", http.StatusBadRequest, resp.StatusCode)
	}

}

func testCreateToDoInternalErrorOnParse(t *testing.T) {

	m := &RepoMock{
		SaveFn: func(*internal.ToDo) error {
			return nil
		},
	}

	req := events.APIGatewayProxyRequest{
		Body:       "garbage",
		HTTPMethod: http.MethodPost,
	}

	resp, err := handlers.NewToDoHandler(m).Handle(req)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(resp.Body, handlers.ErrInternal.Error()) {
		t.Fatalf("Expected body to contain '%s'", handlers.ErrInternal.Error())
	}

	if m.SaveInvoked {
		t.Fatal("Save invoked")
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected %d http response code, got %d", http.StatusInternalServerError, resp.StatusCode)
	}

}

func testCreateToDoInternalErrorOnSave(t *testing.T) {

	m := &RepoMock{
		SaveFn: func(*internal.ToDo) error {
			return errors.New("DB Error")
		},
	}

	req := events.APIGatewayProxyRequest{
		Body:       toDoToString(&newToDo),
		HTTPMethod: http.MethodPost,
	}

	resp, err := handlers.NewToDoHandler(m).Handle(req)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(resp.Body, handlers.ErrInternal.Error()) {
		t.Fatalf("Expected body to contain '%s'", handlers.ErrInternal.Error())
	}

	if !m.SaveInvoked {
		t.Fatal("Save not invoked")
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected %d http response code, got %d", http.StatusInternalServerError, resp.StatusCode)
	}

}

func testUpdateToDoOK(t *testing.T) {

	m := &RepoMock{
		GetFn: func(string) (*internal.ToDo, error) {
			return &savedToDo, nil
		},
		SaveFn: func(*internal.ToDo) error {
			return nil
		},
	}

	req := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": testUUID},
		Body:           toDoToString(&savedToDo),
		HTTPMethod:     http.MethodPut,
	}

	resp, err := handlers.NewToDoHandler(m).Handle(req)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(resp.Body, testUUID) {
		t.Fatalf("Expected body to contain '%s'", testUUID)
	}

	if !m.GetInvoked {
		t.Fatal("Get not invoked")
	}

	if !m.SaveInvoked {
		t.Fatal("Save not invoked")
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected %d http response code, got %d", http.StatusOK, resp.StatusCode)
	}

}

func testUpdateToDoBadRequestMissingID(t *testing.T) {

	m := &RepoMock{
		GetFn: func(string) (*internal.ToDo, error) {
			return &savedToDo, nil
		},
		SaveFn: func(*internal.ToDo) error {
			return nil
		},
	}

	req := events.APIGatewayProxyRequest{
		Body:       toDoToString(&savedToDo),
		HTTPMethod: http.MethodPut,
	}

	resp, err := handlers.NewToDoHandler(m).Handle(req)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(resp.Body, "ID is required") {
		t.Fatalf("Expected body to contain '%s'", "ID is required")
	}

	if m.GetInvoked {
		t.Fatal("Get invoked")
	}

	if m.SaveInvoked {
		t.Fatal("Save invoked")
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected %d http response code, got %d", http.StatusBadRequest, resp.StatusCode)
	}

}

func testUpdateToDoBadRequestNoMatch(t *testing.T) {

	m := &RepoMock{
		GetFn: func(string) (*internal.ToDo, error) {
			return &savedToDo, nil
		},
		SaveFn: func(*internal.ToDo) error {
			return nil
		},
	}

	req := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": "garbage"},
		Body:           toDoToString(&savedToDo),
		HTTPMethod:     http.MethodPut,
	}

	resp, err := handlers.NewToDoHandler(m).Handle(req)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(resp.Body, "ID in body does not match ID in path") {
		t.Fatalf("Expected body to contain '%s'", "ID in body does not match ID in path")
	}

	if m.GetInvoked {
		t.Fatal("Get invoked")
	}

	if m.SaveInvoked {
		t.Fatal("Save invoked")
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected %d http response code, got %d", http.StatusBadRequest, resp.StatusCode)
	}

}

func testUpdateToDoNotFound(t *testing.T) {

	m := &RepoMock{
		GetFn: func(string) (*internal.ToDo, error) {
			return nil, nil
		},
		SaveFn: func(*internal.ToDo) error {
			return nil
		},
	}

	req := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": testUUID},
		Body:           toDoToString(&savedToDo),
		HTTPMethod:     http.MethodPut,
	}

	resp, err := handlers.NewToDoHandler(m).Handle(req)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(resp.Body, handlers.ErrNotFound.Error()) {
		t.Fatalf("Expected body to contain '%s'", handlers.ErrNotFound.Error())
	}

	if !m.GetInvoked {
		t.Fatal("Get not invoked")
	}

	if m.SaveInvoked {
		t.Fatal("Save invoked")
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected %d http response code, got %d", http.StatusNotFound, resp.StatusCode)
	}

}

func testUpdateToDoInternalErrorOnParse(t *testing.T) {

	m := &RepoMock{
		GetFn: func(string) (*internal.ToDo, error) {
			return nil, nil
		},
		SaveFn: func(*internal.ToDo) error {
			return nil
		},
	}

	req := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": testUUID},
		Body:           "garbage",
		HTTPMethod:     http.MethodPut,
	}

	resp, err := handlers.NewToDoHandler(m).Handle(req)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(resp.Body, handlers.ErrInternal.Error()) {
		t.Fatalf("Expected body to contain '%s'", handlers.ErrInternal.Error())
	}

	if m.GetInvoked {
		t.Fatal("Get invoked")
	}

	if m.SaveInvoked {
		t.Fatal("Save invoked")
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected %d http response code, got %d", http.StatusInternalServerError, resp.StatusCode)
	}

}

func testUpdateToDoInternalErrorOnGet(t *testing.T) {

	m := &RepoMock{
		GetFn: func(string) (*internal.ToDo, error) {
			return nil, errors.New("DB Error")
		},
		SaveFn: func(*internal.ToDo) error {
			return nil
		},
	}

	req := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": testUUID},
		Body:           toDoToString(&savedToDo),
		HTTPMethod:     http.MethodPut,
	}

	resp, err := handlers.NewToDoHandler(m).Handle(req)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(resp.Body, handlers.ErrInternal.Error()) {
		t.Fatalf("Expected body to contain '%s'", handlers.ErrInternal.Error())
	}

	if !m.GetInvoked {
		t.Fatal("Get not invoked")
	}

	if m.SaveInvoked {
		t.Fatal("Save invoked")
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected %d http response code, got %d", http.StatusInternalServerError, resp.StatusCode)
	}

}

func testUpdateToDoInternalErrorOnSave(t *testing.T) {

	m := &RepoMock{
		GetFn: func(string) (*internal.ToDo, error) {
			return &savedToDo, nil
		},
		SaveFn: func(*internal.ToDo) error {
			return errors.New("DB Error")
		},
	}

	req := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": testUUID},
		Body:           toDoToString(&savedToDo),
		HTTPMethod:     http.MethodPut,
	}

	resp, err := handlers.NewToDoHandler(m).Handle(req)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(resp.Body, handlers.ErrInternal.Error()) {
		fmt.Println(resp.Body)
		t.Fatalf("Expected body to contain '%s'", handlers.ErrInternal.Error())
	}

	if !m.GetInvoked {
		t.Fatal("Get not invoked")
	}

	if !m.SaveInvoked {
		t.Fatal("Save not invoked")
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected %d http response code, got %d", http.StatusInternalServerError, resp.StatusCode)
	}

}

func testDeleteToDoOK(t *testing.T) {

	m := &RepoMock{
		GetFn: func(string) (*internal.ToDo, error) {
			return &savedToDo, nil
		},
		DeleteFn: func(string) error {
			return nil
		},
	}

	req := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": testUUID},
		HTTPMethod:     http.MethodDelete,
	}

	resp, err := handlers.NewToDoHandler(m).Handle(req)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(resp.Body, "") {
		fmt.Println(resp.Body)
		t.Fatalf("Expected body to contain '%s'", "")
	}

	if !m.GetInvoked {
		t.Fatal("Get not invoked")
	}

	if !m.DeleteInvoked {
		t.Fatal("Delete not invoked")
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected %d http response code, got %d", http.StatusOK, resp.StatusCode)
	}

}

func testDeleteToDoBadRequestMissingID(t *testing.T) {

	m := &RepoMock{
		GetFn: func(string) (*internal.ToDo, error) {
			return &savedToDo, nil
		},
		DeleteFn: func(string) error {
			return nil
		},
	}

	req := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodDelete,
	}

	resp, err := handlers.NewToDoHandler(m).Handle(req)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(resp.Body, "ID is required") {
		fmt.Println(resp.Body)
		t.Fatalf("Expected body to contain '%s'", "ID is required")
	}

	if m.GetInvoked {
		t.Fatal("Get invoked")
	}

	if m.DeleteInvoked {
		t.Fatal("Delete invoked")
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected %d http response code, got %d", http.StatusBadRequest, resp.StatusCode)
	}

}

func testDeleteToDoNotFound(t *testing.T) {

	m := &RepoMock{
		GetFn: func(string) (*internal.ToDo, error) {
			return nil, nil
		},
		DeleteFn: func(string) error {
			return nil
		},
	}

	req := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": testUUID},
		HTTPMethod:     http.MethodDelete,
	}

	resp, err := handlers.NewToDoHandler(m).Handle(req)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(resp.Body, handlers.ErrNotFound.Error()) {
		fmt.Println(resp.Body)
		t.Fatalf("Expected body to contain '%s'", handlers.ErrNotFound.Error())
	}

	if !m.GetInvoked {
		t.Fatal("Get not invoked")
	}

	if m.DeleteInvoked {
		t.Fatal("Delete invoked")
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected %d http response code, got %d", http.StatusNotFound, resp.StatusCode)
	}

}

func testDeleteToDoInternalErrorOnGet(t *testing.T) {

	m := &RepoMock{
		GetFn: func(string) (*internal.ToDo, error) {
			return nil, errors.New("DB Error")
		},
		DeleteFn: func(string) error {
			return nil
		},
	}

	req := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": testUUID},
		HTTPMethod:     http.MethodDelete,
	}

	resp, err := handlers.NewToDoHandler(m).Handle(req)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(resp.Body, handlers.ErrInternal.Error()) {
		fmt.Println(resp.Body)
		t.Fatalf("Expected body to contain '%s'", handlers.ErrInternal.Error())
	}

	if !m.GetInvoked {
		t.Fatal("Get not invoked")
	}

	if m.DeleteInvoked {
		t.Fatal("Delete invoked")
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected %d http response code, got %d", http.StatusInternalServerError, resp.StatusCode)
	}

}

func testDeleteToDoInternalErrorOnDelete(t *testing.T) {

	m := &RepoMock{
		GetFn: func(string) (*internal.ToDo, error) {
			return &savedToDo, nil
		},
		DeleteFn: func(string) error {
			return errors.New("DB Error")
		},
	}

	req := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": testUUID},
		HTTPMethod:     http.MethodDelete,
	}

	resp, err := handlers.NewToDoHandler(m).Handle(req)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(resp.Body, handlers.ErrInternal.Error()) {
		fmt.Println(resp.Body)
		t.Fatalf("Expected body to contain '%s'", handlers.ErrInternal.Error())
	}

	if !m.GetInvoked {
		t.Fatal("Get not invoked")
	}

	if !m.DeleteInvoked {
		t.Fatal("Delete not invoked")
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected %d http response code, got %d", http.StatusInternalServerError, resp.StatusCode)
	}

}

func testMethodNotAllowed(t *testing.T) {

	m := &RepoMock{}

	req := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"id": testUUID},
		HTTPMethod:     http.MethodPatch,
	}

	resp, err := handlers.NewToDoHandler(m).Handle(req)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(resp.Body, handlers.ErrMethodNotAllowed.Error()) {
		fmt.Println(resp.Body)
		t.Fatalf("Expected body to contain '%s'", handlers.ErrMethodNotAllowed.Error())
	}

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("Expected %d http response code, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}

}

func toDoToString(todo *internal.ToDo) string {
	b, _ := json.Marshal(todo)
	return string(b)
}
