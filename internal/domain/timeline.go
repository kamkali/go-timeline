package domain

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

type Process struct {
	ID                  uint
	Name                string
	StartTime           time.Time
	EndTime             time.Time
	ShortDescription    string
	DetailedDescription string
	Graphic             string
	TypeID              uint
}

type Type struct {
	ID      uint
	Name    string
	Color   string
	Events  []Event
	Process []Process
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
