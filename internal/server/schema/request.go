package schema

type Event struct {
    Name                string `json:"name,omitempty"`
    EventTime           string `json:"event_time"`
    ShortDescription    string `json:"short_description,omitempty"`
    DetailedDescription string `json:"detailed_description,omitempty"`
    Graphic             []byte `json:"graphic,omitempty"`
    Type                string `json:"type,omitempty"`
}
