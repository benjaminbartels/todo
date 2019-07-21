package database

import (
	"github.com/benjaminbartels/todo/internal"
)

// ToDoRepo is an interface for database actions
type ToDoRepo interface {
	Get(id string) (*internal.ToDo, error)
	GetAll() ([]internal.ToDo, error)
	Save(todo *internal.ToDo) error
	Delete(id string) error
}
