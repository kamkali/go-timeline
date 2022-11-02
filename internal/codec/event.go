package codec

import (
	"fmt"
	"github.com/kamkali/go-timeline/internal/domain"
	"github.com/kamkali/go-timeline/internal/server/schema"
	"time"
)

func HTTPToDomainEvent(e *schema.Event) (*domain.Event, error) {
	parsedTime, err := time.Parse(time.RFC3339, e.EventTime)
	if err != nil {
		return nil, fmt.Errorf("invalid time")
	}
	domainEvent := &domain.Event{
		Name:                e.Name,
		EventTime:           parsedTime,
		ShortDescription:    e.ShortDescription,
		DetailedDescription: e.DetailedDescription,
		Graphic:             e.Graphic,
		TypeID:              e.TypeID,
	}

	return domainEvent, nil
}
