package timelineservice

import (
    "github.com/kamkali/go-timeline/internal/domain"
    "github.com/kamkali/go-timeline/internal/logger"
    "golang.org/x/net/context"
)

type TimelineService struct {
    log logger.Logger

    repo domain.TimelineRepository
}

func (t TimelineService) CreateEvent(ctx context.Context, event *domain.Event) (uint, error) {
    return t.repo.CreateEvent(ctx, event)
}

func (t TimelineService) ListEvents(ctx context.Context) ([]domain.Event, error) {
    return t.repo.ListEvents(ctx)
}

func New(log logger.Logger, repo domain.TimelineRepository) *TimelineService {
    return &TimelineService{log: log, repo: repo}
}
