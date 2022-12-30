package typeservice

import (
	"github.com/kamkali/go-timeline/internal/domain"
	"github.com/kamkali/go-timeline/internal/logger"
	"golang.org/x/net/context"
)

type TypeService struct {
	log logger.Logger

	repo domain.TypeRepository
}

func (t TypeService) GetType(ctx context.Context, id uint) (domain.Type, error) {
	return t.repo.GetType(ctx, id)
}

func (t TypeService) UpdateType(ctx context.Context, id uint, dt *domain.Type) error {
	return t.repo.UpdateType(ctx, id, dt)
}

func (t TypeService) DeleteType(ctx context.Context, id uint) error {
	return t.repo.DeleteType(ctx, id)
}

func (t TypeService) CreateType(ctx context.Context, dt *domain.Type) (uint, error) {
	return t.repo.CreateType(ctx, dt)
}

func (t TypeService) ListTypes(ctx context.Context) ([]domain.Type, error) {
	return t.repo.ListTypes(ctx)
}

func New(log logger.Logger, repo domain.TypeRepository) *TypeService {
	return &TypeService{log: log, repo: repo}
}
