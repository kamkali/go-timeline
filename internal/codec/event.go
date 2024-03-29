package codec

import (
	"fmt"
	"github.com/kamkali/go-timeline/internal/server/schema"
	"github.com/kamkali/go-timeline/internal/timeline"
	"time"
)

func HTTPToDomainEvent(e *schema.Event) (*timeline.Event, error) {
	parsedTime, err := time.Parse(time.RFC3339, e.EventTime)
	if err != nil {
		return nil, fmt.Errorf("invalid time")
	}
	domainEvent := &timeline.Event{
		Name:                e.Name,
		EventTime:           parsedTime,
		ShortDescription:    e.ShortDescription,
		DetailedDescription: e.DetailedDescription,
		Graphic:             e.Graphic,
		TypeID:              e.TypeID,
	}

	return domainEvent, nil
}

func HTTPFromDomainEvent(e *timeline.Event) (*schema.Event, error) {
	httpEvent := &schema.Event{
		ID:                  e.ID,
		Name:                e.Name,
		EventTime:           e.EventTime.Format(time.RFC3339),
		ShortDescription:    e.ShortDescription,
		DetailedDescription: e.DetailedDescription,
		Graphic:             e.Graphic,
		TypeID:              e.TypeID,
	}

	return httpEvent, nil
}
