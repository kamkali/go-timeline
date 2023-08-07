package service

import (
	"github.com/kamkali/go-timeline/internal/timeline"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type TypeService struct {
	log *zap.Logger

	repo timeline.TypeRepository
}

func (t TypeService) GetType(ctx context.Context, id uint) (timeline.Type, error) {
	return t.repo.GetType(ctx, id)
}

func (t TypeService) UpdateType(ctx context.Context, id uint, dt *timeline.Type) error {
	return t.repo.UpdateType(ctx, id, dt)
}

func (t TypeService) DeleteType(ctx context.Context, id uint) error {
	return t.repo.DeleteType(ctx, id)
}

func (t TypeService) CreateType(ctx context.Context, dt *timeline.Type) (uint, error) {
	return t.repo.CreateType(ctx, dt)
}

func (t TypeService) ListTypes(ctx context.Context) ([]timeline.Type, error) {
	return t.repo.ListTypes(ctx)
}

func NewTypeService(log *zap.Logger, repo timeline.TypeRepository) *TypeService {
	return &TypeService{log: log, repo: repo}
}
