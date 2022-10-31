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
	ID                  uint      `json:"id,omitempty"`
	Name                string    `json:"name,omitempty"`
	EventTime           time.Time `json:"event_time"`
	ShortDescription    string    `json:"short_description,omitempty"`
	DetailedDescription string    `json:"detailed_description,omitempty"`
	Graphic             []byte    `json:"graphic,omitempty"`
	Type                EventType `json:"type,omitempty"`
}

type Process struct {
	ID                  uint      `json:"id,omitempty"`
	Name                string    `json:"name,omitempty"`
	StartTime           time.Time `json:"start_time"`
	EndTime             time.Time `json:"end_time"`
	ShortDescription    string    `json:"short_description,omitempty"`
	DetailedDescription string    `json:"detailed_description,omitempty"`
	Graphic             []byte    `json:"graphic,omitempty"`
	Type                EventType `json:"type,omitempty"`
}

type Type struct {
	ID     uint    `json:"id,omitempty"`
	Name   string  `json:"name,omitempty"`
	Color  string  `json:"color,omitempty"`
	Events []Event `json:"events,omitempty"`
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
	ListTypes(ctx context.Context) ([]Type, error)
	CreateType(ctx context.Context, t *Type) (uint, error)
	GetType(ctx context.Context, id uint) (Type, error)
	UpdateType(ctx context.Context, id uint, Type *Type) error
	DeleteType(ctx context.Context, id uint) error
}

//go:generate mockery --name=TypeService

type TypeRepository interface {
	ListTypes(ctx context.Context) ([]Type, error)
	CreateType(ctx context.Context, t *Type) (uint, error)
	GetType(ctx context.Context, id uint) (Type, error)
	UpdateType(ctx context.Context, id uint, Type *Type) error
	DeleteType(ctx context.Context, id uint) error
}

//go:generate mockery --name=TypeRepository
