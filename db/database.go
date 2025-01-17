package db

import (
	"fmt"
	"net/url"
	"todo-go/core"
	"todo-go/models"
)

type Database interface {
	GetUserID(user models.UserForm) (string, error)
	GetUser(userID string) (models.UserResponse, error)
	GetUsers() ([]models.UserResponse, error)
	AddUser(user models.UserForm) (userID string, err error)
	DelUser(userID string) error

	GetTodo(userID string, todoID string) (models.TodoResponse, error)
	GetTodos(userID string) ([]models.TodoResponse, error)
	AddTodo(userID string, todo models.TodoForm) (todoID string, err error)
	ModTodo(userID string, todoID string, todo models.TodoForm) error
	DelTodo(userID string, todoID string) error

	Close() error
}

func OpenDefault() (Database, error) {
	uri, err := core.LoadEnv("DATABASE_URL")
	if err != nil {
		return nil, err
	}

	parsed, err := url.Parse(uri)
	if parsed.Scheme == "mongodb" {
		return MongoDBOpen(uri)
	}

	return nil, fmt.Errorf("Unknown database URL `%s`", uri)
}
