package codec

import (
	"github.com/kamkali/go-timeline/internal/server/schema"
	"github.com/kamkali/go-timeline/internal/timeline"
)

func HTTPToDomainType(e *schema.Type) (*timeline.Type, error) {
	domainType := &timeline.Type{
		Name:  e.Name,
		Color: e.Color,
	}

	return domainType, nil
}

func HTTPFromDomainType(t *timeline.Type) (*schema.Type, error) {
	return &schema.Type{
		ID:    t.ID,
		Name:  t.Name,
		Color: t.Color,
	}, nil
}
