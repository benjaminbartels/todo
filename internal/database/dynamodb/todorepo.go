package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/benjaminbartels/todo/internal"
	"github.com/pkg/errors"
)

const todosTableName = "todos"

// ToDoRepo represents a boltdb repository for managing todos
type ToDoRepo struct {
	db *dynamodb.DynamoDB
}

// NewToDoRepo returns a new ToDo repository using the given bolt database. It also creates the ToDos
// bucket if it is not yet created on disk.
func NewToDoRepo(db *dynamodb.DynamoDB) *ToDoRepo {
	return &ToDoRepo{db}
}

// Get returns a ToDo by its ID
func (r *ToDoRepo) Get(id string) (*internal.ToDo, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(todosTableName),
		Key:       mapID(id),
	}

	result, err := r.db.GetItem(input)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get ToDo %s from database", id)
	}

	var todo *internal.ToDo

	err = dynamodbattribute.UnmarshalMap(result.Item, todo)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not unmarshal ToDo %s", id)
	}

	return todo, nil
}

// GetAll returns all ToDos
func (r *ToDoRepo) GetAll() ([]internal.ToDo, error) {

	input := &dynamodb.ScanInput{
		TableName: aws.String(todosTableName),
	}

	result, err := r.db.Scan(input)
	if err != nil {
		return nil, errors.Wrap(err, "Could not ToDos from database")
	}

	todos := []internal.ToDo{}

	// Unmarshal the Items field in the result value to the Item Go type.
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &todos)
	if err != nil {
		return nil, errors.Wrap(err, "Could not unmarshal ToDos")
	}

	return todos, nil
}

// Save creates or updates a ToDo
func (r *ToDoRepo) Save(todo *internal.ToDo) error {

	b, err := dynamodbattribute.MarshalMap(todo)
	if err != nil {
		return errors.Wrap(err, "Could not marshal ToDo")
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(todosTableName),
		Item:      b,
	}

	if _, err := r.db.PutItem(input); err != nil {
		return errors.Wrap(err, "Could not put ToDo")
	}

	return nil
}

// Delete permanently removes a ToDo
func (r *ToDoRepo) Delete(id string) error {

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(todosTableName),
		Key:       mapID(id),
	}

	if _, err := r.db.DeleteItem(input); err != nil {
		return errors.Wrapf(err, "Could not delete ToDo %s", id)
	}

	return nil
}
