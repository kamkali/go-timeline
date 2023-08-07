package service

import (
	"github.com/kamkali/go-timeline/internal/timeline"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type EventService struct {
	log *zap.Logger

	repo timeline.EventRepository
}

func (t EventService) GetEvent(ctx context.Context, id uint) (timeline.Event, error) {
	return t.repo.GetEvent(ctx, id)
}

func (t EventService) UpdateEvent(ctx context.Context, id uint, event *timeline.Event) error {
	return t.repo.UpdateEvent(ctx, id, event)
}

func (t EventService) DeleteEvent(ctx context.Context, id uint) error {
	return t.repo.DeleteEvent(ctx, id)
}

func (t EventService) CreateEvent(ctx context.Context, event *timeline.Event) (uint, error) {
	return t.repo.CreateEvent(ctx, event)
}

func (t EventService) ListEvents(ctx context.Context) ([]timeline.Event, error) {
	return t.repo.ListEvents(ctx)
}

func NewEventService(log *zap.Logger, repo timeline.EventRepository) *EventService {
	return &EventService{log: log, repo: repo}
}
