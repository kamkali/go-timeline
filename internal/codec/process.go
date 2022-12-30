package codec

import (
	"fmt"
	"github.com/kamkali/go-timeline/internal/domain"
	"github.com/kamkali/go-timeline/internal/server/schema"
	"time"
)

func HTTPToDomainProcess(e *schema.Process) (*domain.Process, error) {
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
		TypeID:              e.TypeID,
	}

	return domainProcess, nil
}

func HTTPFromDomainProcess(p *domain.Process) (*schema.Process, error) {
	httpProcess := &schema.Process{
		ID:                  p.ID,
		Name:                p.Name,
		StartTime:           p.StartTime.Format(time.RFC3339),
		EndTime:             p.EndTime.Format(time.RFC3339),
		ShortDescription:    p.ShortDescription,
		DetailedDescription: p.DetailedDescription,
		Graphic:             p.Graphic,
		TypeID:              p.TypeID,
	}

	return httpProcess, nil
}
