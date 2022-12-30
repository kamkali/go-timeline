package domain

import (
	"golang.org/x/net/context"
	"time"
)

// TODO: Get rid of JSON annotations

type Event struct {
	ID                  uint      `json:"id,omitempty"`
	Name                string    `json:"name,omitempty"`
	EventTime           time.Time `json:"event_time"`
	ShortDescription    string    `json:"short_description,omitempty"`
	DetailedDescription string    `json:"detailed_description,omitempty"`
	Graphic             string    `json:"graphic,omitempty"`
	TypeID              uint      `json:"type_id,omitempty"`
}

type Process struct {
	ID                  uint      `json:"id,omitempty"`
	Name                string    `json:"name,omitempty"`
	StartTime           time.Time `json:"start_time"`
	EndTime             time.Time `json:"end_time"`
	ShortDescription    string    `json:"short_description,omitempty"`
	DetailedDescription string    `json:"detailed_description,omitempty"`
	Graphic             string    `json:"graphic,omitempty"`
	TypeID              uint      `json:"type_id,omitempty"`
}

type Type struct {
	ID      uint      `json:"id,omitempty"`
	Name    string    `json:"name,omitempty"`
	Color   string    `json:"color,omitempty"`
	Events  []Event   `json:"-"`
	Process []Process `json:"-"`
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

type ProcessService interface {
	ListProcesses(ctx context.Context) ([]Process, error)
	CreateProcess(ctx context.Context, t *Process) (uint, error)
	GetProcess(ctx context.Context, id uint) (Process, error)
	UpdateProcess(ctx context.Context, id uint, process *Process) error
	DeleteProcess(ctx context.Context, id uint) error
}

//go:generate mockery --name=ProcessService

type ProcessRepository interface {
	ListProcesses(ctx context.Context) ([]Process, error)
	CreateProcess(ctx context.Context, t *Process) (uint, error)
	GetProcess(ctx context.Context, id uint) (Process, error)
	UpdateProcess(ctx context.Context, id uint, process *Process) error
	DeleteProcess(ctx context.Context, id uint) error
}

//go:generate mockery --name=ProcessRepository

type Renderer interface {
	RenderSite(events []Event) ([]byte, error)
}
