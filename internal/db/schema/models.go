package schema

import (
	"gorm.io/gorm"
	"time"
)

type EventType string

const (
	EventTypeNormal = "Normal"
)

type Event struct {
	gorm.Model

	Name                string
	EventTime           time.Time
	ShortDescription    string
	DetailedDescription string
	Graphic             []byte
	Type                EventType
	TypeID              string
}

type Process struct {
	gorm.Model

	Name                string
	StartTime           time.Time
	EndTime             time.Time
	ShortDescription    string
	DetailedDescription string
	Graphic             []byte
	Type                EventType
}

type Type struct {
	gorm.Model

	Name   string
	Color  string
	Events []Event
}
