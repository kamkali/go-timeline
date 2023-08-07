package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kamkali/go-timeline/internal/codec"
	schema2 "github.com/kamkali/go-timeline/internal/server/schema"
	timeline2 "github.com/kamkali/go-timeline/internal/timeline"
	"golang.org/x/net/context"
	"io"
	"net/http"
)

func (s *Server) getEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := s.getIDFromRequest(r)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusBadRequest, schema2.ErrBadRequest)
			return
		}

		event, err := s.eventService.GetEvent(ctx, id)
		if err != nil {
			if errors.Is(err, timeline2.ErrNotFound) {
				s.writeErrResponse(w, err, http.StatusNotFound, schema2.ErrNotFound)
				return
			}
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema2.ErrInternal)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		httpEvent, err := codec.HTTPFromDomainEvent(&event)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema2.ErrInternal)
			return
		}
		eventResponse, err := json.Marshal(schema2.EventResponse{Event: httpEvent})
		if err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema2.ErrInternal)
			return
		}
		if _, err := w.Write(eventResponse); err != nil {
			s.log.Error("cannot write response")
			return
		}
	}
}

func (s *Server) updateEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		domainEvent, err := s.getEventPayload(r)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusBadRequest, schema2.ErrBadRequest)
			return
		}

		id, err := s.getIDFromRequest(r)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusBadRequest, schema2.ErrBadRequest)
			return
		}

		if err := s.eventService.UpdateEvent(ctx, id, domainEvent); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				s.writeErrResponse(w, err, http.StatusRequestTimeout, schema2.ErrTimedOut)
				return
			}
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema2.ErrInternal)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) deleteEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := s.getIDFromRequest(r)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusBadRequest, schema2.ErrBadRequest)
			return
		}

		if err := s.eventService.DeleteEvent(ctx, id); err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema2.ErrInternal)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) listEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		events, err := s.eventService.ListEvents(ctx)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				s.writeErrResponse(w, err, http.StatusRequestTimeout, schema2.ErrTimedOut)
				return
			}
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema2.ErrInternal)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		var httpEvents []*schema2.Event
		for i := range events {
			httpEvent, err := codec.HTTPFromDomainEvent(&events[i])
			if err != nil {
				s.writeErrResponse(w, err, http.StatusInternalServerError, schema2.ErrInternal)
				return
			}
			httpEvents = append(httpEvents, httpEvent)
		}

		eventsResponse, err := json.Marshal(schema2.EventsResponse{Events: httpEvents})
		if err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema2.ErrInternal)
			return
		}
		if _, err := w.Write(eventsResponse); err != nil {
			s.log.Error("cannot write response")
			return
		}
	}
}

func (s *Server) createEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		domainEvent, err := s.getEventPayload(r)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusBadRequest, schema2.ErrBadRequest)
			return
		}

		created, err := s.eventService.CreateEvent(ctx, domainEvent)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				s.writeErrResponse(w, err, http.StatusRequestTimeout, schema2.ErrTimedOut)
				return
			}
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema2.ErrInternal)
			return
		}

		er := schema2.EventCreatedResponse{EventID: created}
		resp, err := json.Marshal(er)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema2.ErrInternal)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write(resp); err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema2.ErrInternal)
			return
		}
	}
}

func (s *Server) getEventPayload(r *http.Request) (*timeline2.Event, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read body")
	}
	var event schema2.Event
	if err := json.Unmarshal(body, &event); err != nil {
		return nil, fmt.Errorf("cannot unmarshal body")
	}
	domainEvent, err := codec.HTTPToDomainEvent(&event)
	if err != nil {
		return nil, fmt.Errorf("cannot codec to domain entity")
	}

	return domainEvent, nil
}
