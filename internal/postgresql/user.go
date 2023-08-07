package postgresql

import (
	"errors"
	"fmt"
	timeline2 "github.com/kamkali/go-timeline/internal/timeline"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type UserRepository struct {
	log *zap.Logger

	db *gorm.DB
}

func (ur UserRepository) ChangePassword(ctx context.Context, email, password string) error {
	user, err := ur.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	dbUser, err := toDBUser(user)
	if err != nil {
		return err
	}

	dbUser.ID = user.ID
	dbUser.Password = password
	if err := ur.db.WithContext(ctx).Save(&dbUser).Error; err != nil {
		return fmt.Errorf("db error on update query: %w", err)
	}

	return nil
}

func NewUserRepository(log *zap.Logger, db *gorm.DB) *UserRepository {
	return &UserRepository{log: log, db: db}
}

func toDomainUser(u user) (timeline2.User, error) {
	domainUser := timeline2.User{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
	}
	return domainUser, nil
}

func (ur UserRepository) GetUserByEmail(ctx context.Context, email string) (timeline2.User, error) {
	var t user
	if err := ur.db.WithContext(ctx).Where("email = ?", email).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return timeline2.User{}, timeline2.ErrNotFound
		}
		return timeline2.User{}, fmt.Errorf("db error on select query: %w", err)
	}
	domainUser, err := toDomainUser(t)
	if err != nil {
		return timeline2.User{}, fmt.Errorf("cannot translate db model to domain")
	}
	return domainUser, nil
}

func (ur UserRepository) CreateUser(ctx context.Context, user timeline2.User) error {
	dbType, err := toDBUser(user)
	if err != nil {
		return err
	}

	if err := ur.db.WithContext(ctx).Create(dbType).Error; err != nil {
		return fmt.Errorf("cannot create user: %w", err)
	}
	return nil
}

func toDBUser(u timeline2.User) (*user, error) {
	return &user{
		Email:    u.Email,
		Password: u.Password,
	}, nil
}
