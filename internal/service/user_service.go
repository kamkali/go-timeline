package service

import (
	"fmt"
	timeline2 "github.com/kamkali/go-timeline/internal/timeline"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

type UserService struct {
	log *zap.Logger

	repo timeline2.UserRepository
}

func (t UserService) ChangePassword(ctx context.Context, email, password string) error {
	if password == "" {
		return fmt.Errorf("empty password")
	}
	return t.repo.ChangePassword(ctx, email, password)
}

func (t UserService) LoginUser(ctx context.Context, loggingUser *timeline2.User) (timeline2.User, error) {
	if loggingUser.Email == "" || loggingUser.Password == "" {
		return timeline2.User{}, fmt.Errorf("empty email or password")
	}
	validUser, err := t.repo.GetUserByEmail(ctx, loggingUser.Email)
	if err != nil {
		return timeline2.User{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(validUser.Password), []byte(loggingUser.Password)); err != nil {
		return timeline2.User{}, timeline2.ErrUnauthorized
	}

	return validUser, nil
}

func (t UserService) CreateUser(ctx context.Context, user timeline2.User) error {
	if user.Email == "" || user.Password == "" {
		return fmt.Errorf("empty email or password")
	}
	return t.repo.CreateUser(ctx, user)
}

func NewUserService(log *zap.Logger, repo timeline2.UserRepository) *UserService {
	return &UserService{log: log, repo: repo}
}
