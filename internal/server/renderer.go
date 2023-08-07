package server

import (
	"github.com/kamkali/go-timeline/internal/server/schema"
	"net/http"
)

func (s *Server) renderTimeline() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		events, err := s.eventService.ListEvents(ctx)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
			return
		}

		site, err := s.renderer.RenderSite(events)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
			return
		}
		w.Write(site)
	}
}
