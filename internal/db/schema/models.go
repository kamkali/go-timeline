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
}
