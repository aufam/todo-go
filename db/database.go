package db

import (
	"context"
	"fmt"
	"net/url"
	"todo-go/core"
	"todo-go/models"
)

type Database interface {
	GetUserID(ctx context.Context, user models.UserForm) (string, error)
	GetUser(ctx context.Context, userID string) (models.UserResponse, error)
	GetUsers(ctx context.Context) ([]models.UserResponse, error)
	AddUser(ctx context.Context, user models.UserForm) (userID string, err error)
	DelUser(ctx context.Context, userID string) error

	GetTodo(ctx context.Context, userID string, todoID string) (models.TodoResponse, error)
	GetTodos(ctx context.Context, userID string) ([]models.TodoResponse, error)
	AddTodo(ctx context.Context, userID string, todo models.TodoForm) (todoID string, err error)
	ModTodo(ctx context.Context, userID string, todoID string, todo models.TodoForm) error
	DelTodo(ctx context.Context, userID string, todoID string) error

	Close(ctx context.Context) error
}

func OpenDefault(ctx context.Context) (Database, error) {
	uri, err := core.LoadEnv("DATABASE_URL")
	if err != nil {
		return nil, err
	}

	parsed, err := url.Parse(uri)
	if parsed.Scheme == "mongodb" {
		return MongoDBOpen(ctx, uri)
	}

	return nil, fmt.Errorf("Unknown database URL `%s`", uri)
}
