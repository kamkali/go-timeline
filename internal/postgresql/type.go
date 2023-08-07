package postgresql

import (
	"errors"
	"fmt"
	timeline2 "github.com/kamkali/go-timeline/internal/timeline"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type TypeRepository struct {
	log *zap.Logger

	db *gorm.DB
}

func NewTypeRepository(log *zap.Logger, db *gorm.DB) *TypeRepository {
	return &TypeRepository{log: log, db: db}
}

func toDBType(dt *timeline2.Type) (*eventType, error) {
	return &eventType{
		Name:  dt.Name,
		Color: dt.Color,
	}, nil
}

func (tr TypeRepository) GetType(ctx context.Context, id uint) (timeline2.Type, error) {
	var t eventType
	if err := tr.db.WithContext(ctx).First(&t, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return timeline2.Type{}, timeline2.ErrNotFound
		}
		return timeline2.Type{}, fmt.Errorf("db error on select query: %w", err)
	}
	domainType, err := toDomainType(t)
	if err != nil {
		return timeline2.Type{}, fmt.Errorf("cannot translate db model to domain")
	}
	return domainType, nil
}

func (tr TypeRepository) UpdateType(ctx context.Context, id uint, dt *timeline2.Type) error {
	var t eventType
	r := tr.db.WithContext(ctx).Find(&t, id)
	if r.Error != nil {
		return fmt.Errorf("db error on select query: %w", r.Error)
	}

	t.Name = dt.Name
	t.Color = dt.Color

	if err := tr.db.WithContext(ctx).Save(&t).Error; err != nil {
		return fmt.Errorf("db error on update query: %w", r.Error)
	}

	return nil
}

func (tr TypeRepository) DeleteType(ctx context.Context, id uint) error {
	if err := tr.db.WithContext(ctx).Delete(&eventType{}, id).Error; err != nil {
		return fmt.Errorf("error while deleting: %w", err)
	}
	return nil
}

func (tr TypeRepository) CreateType(ctx context.Context, dt *timeline2.Type) (uint, error) {
	dbType, err := toDBType(dt)
	if err != nil {
		return 0, err
	}

	if err := tr.db.WithContext(ctx).Create(dbType).Error; err != nil {
		return 0, fmt.Errorf("cannot create eventType: %w", err)
	}
	return dbType.ID, nil
}

func (tr TypeRepository) ListTypes(ctx context.Context) ([]timeline2.Type, error) {
	var types []eventType
	r := tr.db.WithContext(ctx).Find(&types)
	if r.Error != nil {
		return nil, fmt.Errorf("db error on select query: %w", r.Error)
	}

	domainTypes := []timeline2.Type{}
	for _, e := range types {
		domainType, err := toDomainType(e)
		if err != nil {
			return nil, fmt.Errorf("cannot translate db model to domain")
		}
		domainTypes = append(domainTypes, domainType)
	}

	return domainTypes, nil
}

func toDomainType(mt eventType) (timeline2.Type, error) {
	domainType := timeline2.Type{
		ID:    mt.ID,
		Name:  mt.Name,
		Color: mt.Color,
	}
	return domainType, nil
}
