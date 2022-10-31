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

type TypeRepository struct {
	log logger.Logger

	db *gorm.DB
}

func NewTypeRepository(log logger.Logger, db *gorm.DB) *TypeRepository {
	return &TypeRepository{log: log, db: db}
}

func toDBType(dt *domain.Type) (*models.Type, error) {
	return &models.Type{
		Name:  dt.Name,
		Color: dt.Color,
	}, nil
}

func (tr TypeRepository) GetType(ctx context.Context, id uint) (domain.Type, error) {
	var t models.Type
	if err := tr.db.WithContext(ctx).First(&t, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Type{}, domain.ErrNotFound
		}
		return domain.Type{}, fmt.Errorf("db error on select query: %w", err)
	}
	domainType, err := toDomainType(t)
	if err != nil {
		return domain.Type{}, fmt.Errorf("cannot translate db model to domain")
	}
	return domainType, nil
}

func (tr TypeRepository) UpdateType(ctx context.Context, id uint, dt *domain.Type) error {
	var t models.Type
	r := tr.db.WithContext(ctx).Find(&t, id)
	if r.Error != nil {
		return fmt.Errorf("db error on select query: %w", r.Error)
	}

	t.Name = dt.Name
	t.Color = dt.Color

	if err := tr.db.Save(&t).Error; err != nil {
		return fmt.Errorf("db error on update query: %w", r.Error)
	}

	return nil
}

func (tr TypeRepository) DeleteType(ctx context.Context, id uint) error {
	if err := tr.db.WithContext(ctx).Delete(&models.Type{}, id).Error; err != nil {
		return fmt.Errorf("error while deleting: %w", err)
	}
	return nil
}

func (tr TypeRepository) CreateType(ctx context.Context, dt *domain.Type) (uint, error) {
	dbType, err := toDBType(dt)
	if err != nil {
		return 0, err
	}

	if err := tr.db.WithContext(ctx).Create(dbType).Error; err != nil {
		return 0, fmt.Errorf("cannot create Type: %w", err)
	}
	return dbType.ID, nil
}

func (tr TypeRepository) ListTypes(ctx context.Context) ([]domain.Type, error) {
	var Types []models.Type
	r := tr.db.WithContext(ctx).Find(&Types)
	if r.Error != nil {
		return nil, fmt.Errorf("db error on select query: %w", r.Error)
	}

	var domainTypes []domain.Type
	for _, e := range Types {
		domainType, err := toDomainType(e)
		if err != nil {
			return nil, fmt.Errorf("cannot translate db model to domain")
		}
		domainTypes = append(domainTypes, domainType)
	}

	return domainTypes, nil
}

func toDomainType(mt models.Type) (domain.Type, error) {
	domainType := domain.Type{
		ID:    mt.ID,
		Name:  mt.Name,
		Color: mt.Color,
	}
	return domainType, nil
}
