package domain

import (
	"golang.org/x/net/context"
)

type User struct {
	ID       uint   `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"-,omitempty"`
}

type UserService interface {
	LoginUser(ctx context.Context, user *User) (User, error)
	CreateUser(ctx context.Context, user User) error
	ChangePassword(ctx context.Context, email, password string) error
}

//go:generate mockery --name=UserService

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (User, error)
	CreateUser(ctx context.Context, user User) error
	ChangePassword(ctx context.Context, email, password string) error
}

//go:generate mockery --name=UserRepository
