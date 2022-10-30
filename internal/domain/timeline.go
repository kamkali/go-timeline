package domain

import (
	"golang.org/x/net/context"
	"time"
)

type EventType string

const (
	EventTypeNormal EventType = "normal"
)

type Event struct {
	ID                  uint
	Name                string
	EventTime           time.Time
	ShortDescription    string
	DetailedDescription string
	Graphic             []byte
	Type                EventType
}

type Process struct {
	Name                string
	StartTime           time.Time
	EndTime             time.Time
	ShortDescription    string
	DetailedDescription string
	Graphic             []byte
	Type                EventType
}

type Type struct {
	Name   string
	Color  string
	Events []Event
}

type EventService interface {
	ListEvents(ctx context.Context) ([]Event, error)
	CreateEvent(ctx context.Context, event *Event) (uint, error)
	GetEvent(ctx context.Context, id uint) (Event, error)
	UpdateEvent(ctx context.Context, id uint, event *Event) error
	DeleteEvent(ctx context.Context, id uint) error
}

//go:generate mockery --name=EventService

type EventRepository interface {
	ListEvents(ctx context.Context) ([]Event, error)
	CreateEvent(ctx context.Context, event *Event) (uint, error)
	GetEvent(ctx context.Context, id uint) (Event, error)
	UpdateEvent(ctx context.Context, id uint, event *Event) error
	DeleteEvent(ctx context.Context, id uint) error
}

//go:generate mockery --name=EventRepository

type TypeService interface {
	CreateType(ctx context.Context, t *Type) (Type, error)
}

//go:generate mockery --name=TypeService

type TypeRepository interface {
	CreateType(ctx context.Context, t *Type) (Type, error)
}

//go:generate mockery --name=TypeRepository
