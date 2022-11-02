package generator

import (
	"bytes"
	"fmt"
	"github.com/kamkali/go-timeline/internal/domain"
	"github.com/kamkali/go-timeline/internal/logger"
	"html/template"
	"os"
	"sort"
)

type Renderer struct {
	log logger.Logger

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
	tmplName := "./internal/generator/template.gohtml"
	htmlContent, err := os.ReadFile(tmplName)
	if err != nil {
		return err
	}
	r.template = template.Must(template.New("").Parse(string(htmlContent)))
	return nil
}

type data struct {
	Events []domain.Event
}

func (d *data) Sort() {
	sort.Sort(d)
}

func (d *data) Len() int      { return len(d.Events) }
func (d *data) Swap(i, j int) { d.Events[i], d.Events[j] = d.Events[j], d.Events[i] }
func (d *data) Less(i, j int) bool {
	return d.Events[i].EventTime.Year() > d.Events[j].EventTime.Year()
}

func (r *Renderer) RenderSite(events []domain.Event) ([]byte, error) {
	d := data{Events: events}
	d.Sort()

	var buf bytes.Buffer
	if err := r.template.Execute(&buf, d); err != nil {
		return nil, fmt.Errorf("execute tmpl: %w", err)
	}
	return buf.Bytes(), nil
}
