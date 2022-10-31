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

type EventRepository struct {
	log logger.Logger

	db *gorm.DB
}

func NewEventRepository(log logger.Logger, db *gorm.DB) *EventRepository {
	return &EventRepository{log: log, db: db}
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

func (t EventRepository) GetEvent(ctx context.Context, id uint) (domain.Event, error) {
	var event models.Event
	if err := t.db.WithContext(ctx).First(&event, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Event{}, domain.ErrNotFound
		}
		return domain.Event{}, fmt.Errorf("db error on select query: %w", err)
	}
	domainEvent, err := toDomainEvent(event)
	if err != nil {
		return domain.Event{}, fmt.Errorf("cannot translate db model to domain")
	}
	return domainEvent, nil
}

func (t EventRepository) UpdateEvent(ctx context.Context, id uint, event *domain.Event) error {
	var e models.Event
	r := t.db.WithContext(ctx).Find(&e, id)
	if r.Error != nil {
		return fmt.Errorf("db error on select query: %w", r.Error)
	}
	eventType, ok := toDBEventType[event.Type]
	if !ok {
		return fmt.Errorf("unknown event type")
	}

	e.Name = event.Name
	e.EventTime = event.EventTime
	e.ShortDescription = event.ShortDescription
	e.DetailedDescription = event.DetailedDescription
	e.Graphic = event.Graphic
	e.Type = eventType

	if err := t.db.Save(&e).Error; err != nil {
		return fmt.Errorf("db error on update query: %w", r.Error)
	}

	return nil
}

func (t EventRepository) DeleteEvent(ctx context.Context, id uint) error {
	if err := t.db.WithContext(ctx).Delete(&models.Event{}, id).Error; err != nil {
		return fmt.Errorf("error while deleting: %w", err)
	}
	return nil
}

func (t EventRepository) CreateEvent(ctx context.Context, event *domain.Event) (uint, error) {
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

func (t EventRepository) ListEvents(ctx context.Context) ([]domain.Event, error) {
	var events []models.Event
	r := t.db.WithContext(ctx).Find(&events)
	if r.Error != nil {
		return nil, fmt.Errorf("db error on select query: %w", r.Error)
	}

	domainEvents := []domain.Event{}
	for _, e := range events {
		domainEvent, err := toDomainEvent(e)
		if err != nil {
			return nil, fmt.Errorf("cannot translate db model to domain")
		}
		domainEvents = append(domainEvents, domainEvent)
	}

	return domainEvents, nil
}

func toDomainEvent(e models.Event) (domain.Event, error) {
	eventType, ok := toDomainEventType[e.Type]
	if !ok {
		return domain.Event{}, fmt.Errorf("unknown event type")
	}
	domainEvent := domain.Event{
		ID:                  e.ID,
		Name:                e.Name,
		EventTime:           e.EventTime,
		ShortDescription:    e.ShortDescription,
		DetailedDescription: e.DetailedDescription,
		Graphic:             e.Graphic,
		Type:                eventType,
	}
	return domainEvent, nil
}
