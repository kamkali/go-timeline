package db

import (
    "fmt"
    "github.com/kamkali/go-timeline/internal/db/schema/models"
    "github.com/kamkali/go-timeline/internal/domain"
    "github.com/kamkali/go-timeline/internal/logger"
    "github.com/kamkali/go-timeline/internal/server/schema"
    "golang.org/x/net/context"
    "gorm.io/gorm"
    "time"
)

type TimelineRepository struct {
    log logger.Logger

    db *gorm.DB
}

func NewTimelineRepository(log logger.Logger, db *gorm.DB) *TimelineRepository {
    return &TimelineRepository{log: log, db: db}
}

var toDBEventType = map[domain.EventType]models.EventType{
    domain.EventTypeNormal: models.EventTypeNormal,
}

var toDomainEventType = map[models.EventType]domain.EventType{
    models.EventTypeNormal: domain.EventTypeNormal,
}

func toDBEvent(de *domain.Event) (*models.Event, error) {
    dbType, ok := toDBEventType[de.Type]
    if !ok {
        return nil, fmt.Errorf("unknown event type")
    }
    return &models.Event{
        Name:                de.Name,
        EventTime:           de.EventTime,
        ShortDescription:    de.ShortDescription,
        DetailedDescription: de.DetailedDescription,
        Graphic:             de.Graphic,
        Type:                dbType,
    }, nil
}

func (t TimelineRepository) CreateEvent(ctx context.Context, event *domain.Event) (uint, error) {
    dbEvent, err := toDBEvent(event)
    if err != nil {
        return 0, err
    }

    typ := models.Type{}
    result := t.db.Table("types").Find(&typ, "name = ?", event.Type)
    if result.Error != nil || result.RowsAffected == 0 {
        return 0, fmt.Errorf("cannot find type of name %s", event.Type)
    }
    dbEvent.TypeID = typ.ID

    result = t.db.WithContext(ctx).Create(dbEvent)
    if result.Error != nil {
        return 0, fmt.Errorf("cannot create event: %w", result.Error)
    }
    return dbEvent.ID, nil
}

func (t TimelineRepository) ListEvents(ctx context.Context) ([]domain.Event, error) {
    var events []schema.Event
    r := t.db.WithContext(ctx).Find(&events)
    if r.Error != nil {
        return nil, fmt.Errorf("db error on select query: %w", r.Error)
    }

    var domainEvents []domain.Event
    for _, e := range events {
        domainEvent, err := toDomainEvent(e)
        if err != nil {
            return nil, fmt.Errorf("cannot translate db model to domain")
        }
        domainEvents = append(domainEvents, domainEvent)
    }

    return domainEvents, nil
}

func toDomainEvent(e schema.Event) (domain.Event, error) {
    t, err := time.Parse(time.RFC3339, e.EventTime)
    if err != nil {
        return domain.Event{}, fmt.Errorf("cannot parse time: %w", err)
    }
    eventType, ok := toDomainEventType[models.EventType(e.Type)]
    if !ok {
        return domain.Event{}, fmt.Errorf("unknown event type")
    }
    domainEvent := domain.Event{
        Name:                e.Name,
        EventTime:           t,
        ShortDescription:    e.ShortDescription,
        DetailedDescription: e.DetailedDescription,
        Graphic:             e.Graphic,
        Type:                eventType,
    }
    return domainEvent, nil
}
