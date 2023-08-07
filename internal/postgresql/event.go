package postgresql

import (
	"errors"
	"fmt"
	timeline2 "github.com/kamkali/go-timeline/internal/timeline"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type EventRepository struct {
	log *zap.Logger

	db *gorm.DB
}

func NewEventRepository(log *zap.Logger, db *gorm.DB) *EventRepository {
	return &EventRepository{log: log, db: db}
}

func toDBEvent(de *timeline2.Event) (*event, error) {
	return &event{
		Name:                de.Name,
		EventTime:           de.EventTime,
		ShortDescription:    de.ShortDescription,
		DetailedDescription: de.DetailedDescription,
		Graphic:             de.Graphic,
	}, nil
}

func (t EventRepository) GetEvent(ctx context.Context, id uint) (timeline2.Event, error) {
	var event event
	if err := t.db.WithContext(ctx).First(&event, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return timeline2.Event{}, timeline2.ErrNotFound
		}
		return timeline2.Event{}, fmt.Errorf("db error on select query: %w", err)
	}
	domainEvent, err := toDomainEvent(event)
	if err != nil {
		return timeline2.Event{}, fmt.Errorf("cannot translate db model to domain")
	}
	return domainEvent, nil
}

func (t EventRepository) UpdateEvent(ctx context.Context, id uint, domainEvent *timeline2.Event) error {
	var e event
	r := t.db.WithContext(ctx).Find(&e, id)
	if r.Error != nil {
		return fmt.Errorf("db error on select query: %w", r.Error)
	}

	e.Name = domainEvent.Name
	e.EventTime = domainEvent.EventTime
	e.ShortDescription = domainEvent.ShortDescription
	e.DetailedDescription = domainEvent.DetailedDescription
	e.Graphic = domainEvent.Graphic
	e.TypeID = domainEvent.TypeID

	if err := t.db.WithContext(ctx).Save(&e).Error; err != nil {
		return fmt.Errorf("db error on update query: %w", r.Error)
	}

	return nil
}

func (t EventRepository) DeleteEvent(ctx context.Context, id uint) error {
	if err := t.db.WithContext(ctx).Delete(&event{}, id).Error; err != nil {
		return fmt.Errorf("error while deleting: %w", err)
	}
	return nil
}

func (t EventRepository) CreateEvent(ctx context.Context, event *timeline2.Event) (uint, error) {
	dbEvent, err := toDBEvent(event)
	if err != nil {
		return 0, err
	}

	typ := eventType{}
	result := t.db.Model(eventType{}).Find(&typ, event.TypeID)
	if result.Error != nil || result.RowsAffected == 0 {
		return 0, fmt.Errorf("cannot find type %d", event.TypeID)
	}
	dbEvent.TypeID = typ.ID

	result = t.db.WithContext(ctx).Create(dbEvent)
	if result.Error != nil {
		return 0, fmt.Errorf("cannot create event: %w", result.Error)
	}
	return dbEvent.ID, nil
}

func (t EventRepository) ListEvents(ctx context.Context) ([]timeline2.Event, error) {
	var events []event
	r := t.db.WithContext(ctx).Find(&events)
	if r.Error != nil {
		return nil, fmt.Errorf("db error on select query: %w", r.Error)
	}

	domainEvents := []timeline2.Event{}
	for _, e := range events {
		domainEvent, err := toDomainEvent(e)
		if err != nil {
			return nil, fmt.Errorf("cannot translate db model to domain")
		}
		domainEvents = append(domainEvents, domainEvent)
	}

	return domainEvents, nil
}

func toDomainEvent(e event) (timeline2.Event, error) {
	domainEvent := timeline2.Event{
		ID:                  e.ID,
		Name:                e.Name,
		EventTime:           e.EventTime,
		ShortDescription:    e.ShortDescription,
		DetailedDescription: e.DetailedDescription,
		Graphic:             e.Graphic,
		TypeID:              e.TypeID,
	}
	return domainEvent, nil
}
