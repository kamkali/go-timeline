package models

import (
    "gorm.io/gorm"
    "time"
)

type EventType string

const (
    EventTypeNormal EventType = "Normal"
)

type Event struct {
    gorm.Model

    Name                string
    EventTime           time.Time
    ShortDescription    string
    DetailedDescription string
    Graphic             []byte
    Type                EventType
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
    Type                EventType
}

type Type struct {
    gorm.Model

    Name   string
    Color  string
    Events []Event
}
