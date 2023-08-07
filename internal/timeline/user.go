package timeline

import (
	"golang.org/x/net/context"
)

type User struct {
	ID       uint
	Email    string
	Password string
}

type UserService interface {
	LoginUser(ctx context.Context, user *User) (User, error)
	CreateUser(ctx context.Context, user User) error
	ChangePassword(ctx context.Context, email, password string) error
}

//go:generate mockery --output=../mocks --name=UserService

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (User, error)
	CreateUser(ctx context.Context, user User) error
	ChangePassword(ctx context.Context, email, password string) error
}

//go:generate mockery --output=../mocks --name=UserRepository
