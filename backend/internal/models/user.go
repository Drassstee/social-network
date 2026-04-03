package models

import (
	"context"
	"time"
)

type User struct {
	ID          int       `json:"id"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Password    string    `json:"-"`
	BirthDay    time.Time `json:"date_of_birth"`
	Avatar      string    `json:"avatar,omitempty"`
	Nickname    string    `json:"nickname,omitempty"`
	Bio         string    `json:"about_me,omitempty"`
	ProfileType string    `json:"profile_type"`
}

type UserService interface {
	GetByID(ctx context.Context, id int) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByIDs(ctx context.Context, ids []int) ([]User, error)
}

type UserRepo interface {
	GetByID(ctx context.Context, id int) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByIDs(ctx context.Context, ids []int) ([]User, error)
	CreateUser(user *User) error
}

type UserHandler interface {
}
