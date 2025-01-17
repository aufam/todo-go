package models

import (
	"time"
	"todo-go/core"
)

// database model
type User struct {
	ID             string    `json:"id,omitempty" bson:"_id,omitempty"`
	Username       string    `json:"username"`
	HashedPassword string    `json:"hashedPassword"`
	CreatedAt      time.Time `json:"createdAt"`
}

// response model
type UserResponse struct {
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}

// request model
type UserForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) AsResponse() UserResponse {
	return UserResponse{
		Username:  u.Username,
		CreatedAt: u.CreatedAt,
	}
}

func (u *UserForm) CreateModel() (User, error) {
	hashedPassword, err := core.HashPassword(u.Password)
	if err != nil {
		return User{}, err
	}

	return User{
		Username:       u.Username,
		HashedPassword: string(hashedPassword),
		CreatedAt:      time.Now(),
	}, nil
}
