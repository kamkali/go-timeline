package codec

import (
	"github.com/kamkali/go-timeline/internal/server/schema"
	"github.com/kamkali/go-timeline/internal/timeline"
)

func HTTPToDomainUser(u *schema.User) (*timeline.User, error) {
	domainUser := &timeline.User{
		Email:    u.Username,
		Password: u.Password,
	}

	return domainUser, nil
}
