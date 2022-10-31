package schema

import (
	"github.com/kamkali/go-timeline/internal/domain"
)

const (
	ErrInternal   = "Internal server error"
	ErrBadRequest = "Bad request"
	ErrNotFound   = "Not Found"
	ErrTimedOut   = "Timed out"
)

type ServerError struct {
	Description string `json:"description"`
}

type EventCreatedResponse struct {
	EventID uint `json:"event_id,omitempty"`
}

type EventResponse struct {
	Event domain.Event `json:"event"`
}

type EventsResponse struct {
	Events []domain.Event `json:"events"`
}

type TypeResponse struct {
	Type domain.Type `json:"type"`
}

type TypesResponse struct {
	Types []domain.Type `json:"types"`
}

type TypeCreatedResponse struct {
	TypeID uint `json:"type_id,omitempty"`
}