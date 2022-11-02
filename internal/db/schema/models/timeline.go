package models

import (
	"gorm.io/gorm"
	"time"
)

type Event struct {
	gorm.Model

	Name                string
	EventTime           time.Time
	ShortDescription    string
	DetailedDescription string
	Graphic             []byte
	TypeID              uint
}

type Process struct {
	gorm.Model

	Name                string
	StartTime           time.Time
	EndTime             time.Time
	ShortDescription    string
	DetailedDescription string
	Graphic             []byte
	TypeID              uint
}

type Type struct {
	gorm.Model

	Name      string
	Color     string
	Events    []Event
	Processes []Process
}
