package codec

import (
	"fmt"
	"github.com/kamkali/go-timeline/internal/domain"
	"github.com/kamkali/go-timeline/internal/server/schema"
	"time"
)

func HTTPToDomainProcess(e *schema.Process) (*domain.Process, error) {
	eventType, ok := HTTPEventTypeToDomain[e.Type]
	if !ok {
		return nil, fmt.Errorf("unknown event type")
	}
	parsedStartTime, err := time.Parse(time.RFC3339, e.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid time")
	}
	parsedEndTime, err := time.Parse(time.RFC3339, e.EndTime)
	if err != nil {
		return nil, fmt.Errorf("invalid time")
	}
	domainProcess := &domain.Process{
		Name:                e.Name,
		StartTime:           parsedStartTime,
		EndTime:             parsedEndTime,
		ShortDescription:    e.ShortDescription,
		DetailedDescription: e.DetailedDescription,
		Graphic:             e.Graphic,
		Type:                eventType,
	}

	return domainProcess, nil
}
