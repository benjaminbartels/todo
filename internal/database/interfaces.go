package database

import (
	"github.com/benjaminbartels/todo/internal"
)

type ToDoRepo interface {
	Get(id string) (*internal.ToDo, error)
	GetAll() ([]internal.ToDo, error)
	Save(todo *internal.ToDo) error
	Delete(id string) error
}
