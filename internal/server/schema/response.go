package schema

import (
    "github.com/kamkali/go-timeline/internal/domain"
)

const (
    ErrInternal   = "Internal server error"
    ErrBadRequest = "Bad request"
    ErrTimedOut   = "Timed out"
)

type ServerError struct {
    Description string `json:"description"`
}

type EventResponse struct {
    EventID uint `json:"event_id,omitempty"`
}

type EventsResponse struct {
    Events []domain.Event `json:"events"`
}
