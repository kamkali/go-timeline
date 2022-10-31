package codec

import (
	"github.com/kamkali/go-timeline/internal/domain"
	"github.com/kamkali/go-timeline/internal/server/schema"
)

func HTTPToDomainType(e *schema.Type) (*domain.Type, error) {
	domainType := &domain.Type{
		Name:  e.Name,
		Color: e.Color,
	}

	return domainType, nil
}
