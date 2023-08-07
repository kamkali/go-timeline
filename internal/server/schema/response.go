package schema

const (
	ErrInternal     = "Internal server error"
	ErrBadRequest   = "Bad request"
	ErrNotFound     = "Not Found"
	ErrTimedOut     = "Timed out"
	ErrUnauthorized = "Unauthorized"
)

type ServerError struct {
	Description string `json:"description"`
}

type TokenResponse struct {
	Token string `json:"token,omitempty"`
}

type (
	EventCreatedResponse struct {
		EventID uint `json:"event_id,omitempty"`
	}

	EventResponse struct {
		Event *Event `json:"event"`
	}

	EventsResponse struct {
		Events []*Event `json:"events"`
	}
)

type (
	TypeResponse struct {
		Type *Type `json:"type"`
	}

	TypesResponse struct {
		Types []*Type `json:"types"`
	}

	TypeCreatedResponse struct {
		TypeID uint `json:"type_id,omitempty"`
	}
)
