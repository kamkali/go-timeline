package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/kamkali/go-timeline/internal/timeline"
	"html/template"
	"sort"
	"time"
)

//go:embed template.gohtml
var htmlContent string

type Renderer struct {
	template *template.Template
}

func NewRenderer() (*Renderer, error) {
	r := &Renderer{}
	if err := r.loadTemplate(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Renderer) loadTemplate() error {
	r.template = template.Must(template.New("").Parse(htmlContent))
	return nil
}

type Event struct {
	ID                  uint
	Name                string
	EventTime           time.Time
	ShortDescription    string
	DetailedDescription string
	Graphic             template.URL
	TypeID              uint
}

type data struct {
	Events []Event
}

func (d *data) Sort() {
	sort.Sort(d)
}

func (d *data) Len() int      { return len(d.Events) }
func (d *data) Swap(i, j int) { d.Events[i], d.Events[j] = d.Events[j], d.Events[i] }
func (d *data) Less(i, j int) bool {
	return d.Events[i].EventTime.Year() > d.Events[j].EventTime.Year()
}

func (r *Renderer) RenderSite(events []timeline.Event) ([]byte, error) {
	var d data
	for _, e := range events {
		d.Events = append(d.Events, Event{
			ID:                  e.ID,
			Name:                e.Name,
			EventTime:           e.EventTime,
			ShortDescription:    e.ShortDescription,
			DetailedDescription: e.DetailedDescription,
			Graphic:             template.URL(e.Graphic),
			TypeID:              e.TypeID,
		})
	}

	d.Sort()

	var buf bytes.Buffer
	if err := r.template.Execute(&buf, d); err != nil {
		return nil, fmt.Errorf("execute tmpl: %w", err)
	}
	return buf.Bytes(), nil
}
