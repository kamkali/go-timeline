package timeline

import (
	"golang.org/x/net/context"
	"time"
)

type Event struct {
	ID                  uint
	Name                string
	EventTime           time.Time
	ShortDescription    string
	DetailedDescription string
	Graphic             string
	TypeID              uint
}

type EventService interface {
	ListEvents(ctx context.Context) ([]Event, error)
	CreateEvent(ctx context.Context, event *Event) (uint, error)
	GetEvent(ctx context.Context, id uint) (Event, error)
	UpdateEvent(ctx context.Context, id uint, event *Event) error
	DeleteEvent(ctx context.Context, id uint) error
}

//go:generate mockery --output=../mocks --name=EventService

type EventRepository interface {
	ListEvents(ctx context.Context) ([]Event, error)
	CreateEvent(ctx context.Context, event *Event) (uint, error)
	GetEvent(ctx context.Context, id uint) (Event, error)
	UpdateEvent(ctx context.Context, id uint, event *Event) error
	DeleteEvent(ctx context.Context, id uint) error
}

//go:generate mockery --output=../mocks --name=EventRepository

type Type struct {
	ID     uint
	Name   string
	Color  string
	Events []Event
}

type TypeService interface {
	ListTypes(ctx context.Context) ([]Type, error)
	CreateType(ctx context.Context, t *Type) (uint, error)
	GetType(ctx context.Context, id uint) (Type, error)
	UpdateType(ctx context.Context, id uint, Type *Type) error
	DeleteType(ctx context.Context, id uint) error
}

//go:generate mockery --output=../mocks --name=TypeService

type TypeRepository interface {
	ListTypes(ctx context.Context) ([]Type, error)
	CreateType(ctx context.Context, t *Type) (uint, error)
	GetType(ctx context.Context, id uint) (Type, error)
	UpdateType(ctx context.Context, id uint, Type *Type) error
	DeleteType(ctx context.Context, id uint) error
}

//go:generate mockery --output=../mocks --name=TypeRepository
