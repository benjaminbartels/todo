package handlers_test

import (
	"github.com/benjaminbartels/todo/internal"
)

// ClientMock is used to mock a client that uses makes call to DynamoDBAPI
type RepoMock struct {
	GetFn         func(string) (*internal.ToDo, error)
	GetAllFn      func() ([]internal.ToDo, error)
	SaveFn        func(todo *internal.ToDo) error
	DeleteFn      func(string) error
	GetInvoked    bool
	GetAllInvoked bool
	SaveInvoked   bool
	DeleteInvoked bool
}

// Get returns a ToDo by its ID
func (m *RepoMock) Get(id string) (*internal.ToDo, error) {
	m.GetInvoked = true
	return m.GetFn(id)
}

// GetAll returns all ToDos
func (m *RepoMock) GetAll() ([]internal.ToDo, error) {
	m.GetAllInvoked = true
	return m.GetAllFn()
}

// Save creates or updates a ToDo
func (m *RepoMock) Save(todo *internal.ToDo) error {
	m.SaveInvoked = true
	return m.SaveFn(todo)
}

// Delete permanently removes a ToDo
func (m *RepoMock) Delete(id string) error {
	m.DeleteInvoked = true
	return m.DeleteFn(id)
}
