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

type TimelineService interface {
    ListEvents(ctx context.Context) ([]Event, error)
    CreateEvent(ctx context.Context, event *Event) (uint, error)
}

//go:generate mockery --name=TimelineService

type TimelineRepository interface {
    ListEvents(ctx context.Context) ([]Event, error)
    CreateEvent(ctx context.Context, event *Event) (uint, error)
}

//go:generate mockery --name=TimelineRepository
