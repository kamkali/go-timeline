package codec

import (
	"github.com/kamkali/go-timeline/internal/domain"
	"github.com/kamkali/go-timeline/internal/server/schema"
)

func HTTPToDomainUser(u *schema.User) (*domain.User, error) {
	domainUser := &domain.User{
		Email:    u.Username,
		Password: u.Password,
	}

	return domainUser, nil
}
