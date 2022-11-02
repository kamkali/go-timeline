package userservice

import (
	"fmt"
	"github.com/kamkali/go-timeline/internal/domain"
	"github.com/kamkali/go-timeline/internal/logger"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

type UserService struct {
	log logger.Logger

	repo domain.UserRepository
}

func (t UserService) ChangePassword(ctx context.Context, email, password string) error {
	if password == "" {
		return fmt.Errorf("empty password")
	}
	return t.repo.ChangePassword(ctx, email, password)
}

func (t UserService) LoginUser(ctx context.Context, loggingUser *domain.User) (domain.User, error) {
	if loggingUser.Email == "" || loggingUser.Password == "" {
		return domain.User{}, fmt.Errorf("empty email or password")
	}
	validUser, err := t.repo.GetUserByEmail(ctx, loggingUser.Email)
	if err != nil {
		return domain.User{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(validUser.Password), []byte(loggingUser.Password)); err != nil {
		return domain.User{}, domain.ErrUnauthorized
	}

	return validUser, nil
}

func (t UserService) CreateUser(ctx context.Context, user domain.User) error {
	if user.Email == "" || user.Password == "" {
		return fmt.Errorf("empty email or password")
	}
	return t.repo.CreateUser(ctx, user)
}

func New(log logger.Logger, repo domain.UserRepository) *UserService {
	return &UserService{log: log, repo: repo}
}
