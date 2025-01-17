package models

import (
	"fmt"
	"strings"
	"time"
)

// database model
type Todo struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    string    `json:"userId"`
	Task      string    `json:"task"`
	IsDone    bool      `json:"isDone"`
	CreatedAt time.Time `json:"createdAt"`
}

// response model
type TodoResponse struct {
	ID        string    `json:"id"`
	Task      string    `json:"task"`
	IsDone    bool      `json:"isDone"`
	CreatedAt time.Time `json:"createdAt"`
}

// request model
type TodoForm struct {
	Task   string `json:"task,omitempty"`
	IsDone bool   `json:"isDone,omitempty"`
}

func (t *Todo) AsResponse() TodoResponse {
	return TodoResponse{
		ID:        t.ID,
		Task:      t.Task,
		IsDone:    t.IsDone,
		CreatedAt: t.CreatedAt,
	}
}

func (t *TodoForm) CreateModel(userID string) (Todo, error) {
	if strings.TrimSpace(t.Task) == "" {
		return Todo{}, fmt.Errorf("Task cannot be empty")
	}

	return Todo{
		UserID:    userID,
		Task:      t.Task,
		IsDone:    t.IsDone,
		CreatedAt: time.Now(),
	}, nil
}
