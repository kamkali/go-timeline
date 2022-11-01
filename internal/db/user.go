package db

import (
	"errors"
	"fmt"
	"github.com/kamkali/go-timeline/internal/db/schema/models"
	"github.com/kamkali/go-timeline/internal/domain"
	"github.com/kamkali/go-timeline/internal/logger"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type UserRepository struct {
	log logger.Logger

	db *gorm.DB
}

func NewUserRepository(log logger.Logger, db *gorm.DB) *UserRepository {
	return &UserRepository{log: log, db: db}
}

func toDomainUser(u models.User) (domain.User, error) {
	domainUser := domain.User{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
	}
	return domainUser, nil
}

func (ur UserRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	var t models.User
	if err := ur.db.WithContext(ctx).Where("email = ?", email).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("db error on select query: %w", err)
	}
	domainUser, err := toDomainUser(t)
	if err != nil {
		return domain.User{}, fmt.Errorf("cannot translate db model to domain")
	}
	return domainUser, nil
}

func (ur UserRepository) CreateUser(ctx context.Context, user domain.User) error {
	dbType, err := toDBUser(user)
	if err != nil {
		return err
	}

	if err := ur.db.WithContext(ctx).Create(dbType).Error; err != nil {
		return fmt.Errorf("cannot create User: %w", err)
	}
	return nil
}

func toDBUser(user domain.User) (*models.User, error) {
	return &models.User{
		Email:    user.Email,
		Password: user.Password,
	}, nil
}
