package schema

type Event struct {
	ID                  uint   `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	EventTime           string `json:"event_time"`
	ShortDescription    string `json:"short_description,omitempty"`
	DetailedDescription string `json:"detailed_description,omitempty"`
	Graphic             string `json:"graphic,omitempty"`
	TypeID              uint   `json:"type_id,omitempty"`
}

type Type struct {
	ID    uint   `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
}

type User struct {
	ID       uint   `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type PasswordChange struct {
	NewPassword string `json:"new_password,omitempty"`
}
