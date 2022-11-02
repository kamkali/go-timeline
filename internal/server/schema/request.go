package schema

type Event struct {
	Name                string `json:"name,omitempty"`
	EventTime           string `json:"event_time"`
	ShortDescription    string `json:"short_description,omitempty"`
	DetailedDescription string `json:"detailed_description,omitempty"`
	Graphic             []byte `json:"graphic,omitempty"`
	Type                string `json:"type,omitempty"`
}

type Type struct {
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
}

type Process struct {
	Name                string `json:"name,omitempty"`
	StartTime           string `json:"start_time"`
	EndTime             string `json:"end_time"`
	ShortDescription    string `json:"short_description,omitempty"`
	DetailedDescription string `json:"detailed_description,omitempty"`
	Graphic             []byte `json:"graphic,omitempty"`
	Type                string `json:"type,omitempty"`
}

type User struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type PasswordChange struct {
	NewPassword string `json:"new_password,omitempty"`
}
