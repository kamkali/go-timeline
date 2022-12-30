package eventservice

import (
	"github.com/kamkali/go-timeline/internal/domain"
	"github.com/kamkali/go-timeline/internal/logger"
	"golang.org/x/net/context"
)

type EventService struct {
	log logger.Logger

	repo domain.EventRepository
}

func (t EventService) GetEvent(ctx context.Context, id uint) (domain.Event, error) {
	return t.repo.GetEvent(ctx, id)
}

func (t EventService) UpdateEvent(ctx context.Context, id uint, event *domain.Event) error {
	return t.repo.UpdateEvent(ctx, id, event)
}

func (t EventService) DeleteEvent(ctx context.Context, id uint) error {
	return t.repo.DeleteEvent(ctx, id)
}

func (t EventService) CreateEvent(ctx context.Context, event *domain.Event) (uint, error) {
	return t.repo.CreateEvent(ctx, event)
}

func (t EventService) ListEvents(ctx context.Context) ([]domain.Event, error) {
	return t.repo.ListEvents(ctx)
}

func New(log logger.Logger, repo domain.EventRepository) *EventService {
	return &EventService{log: log, repo: repo}
}
