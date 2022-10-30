package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kamkali/go-timeline/internal/codec"
	"github.com/kamkali/go-timeline/internal/domain"
	"github.com/kamkali/go-timeline/internal/server/schema"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"strconv"
)

func (s *Server) getEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := s.getIDFromRequest(r)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusBadRequest, schema.ErrBadRequest)
			return
		}

		event, err := s.eventService.GetEvent(ctx, id)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				s.writeErrResponse(w, err, http.StatusNotFound, schema.ErrNotFound)
				return
			}
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		eventResponse, err := json.Marshal(schema.EventResponse{Event: event})
		if err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
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
			s.writeErrResponse(w, err, http.StatusBadRequest, schema.ErrBadRequest)
			return
		}

		id, err := s.getIDFromRequest(r)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusBadRequest, schema.ErrBadRequest)
			return
		}

		if err := s.eventService.UpdateEvent(ctx, id, domainEvent); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				s.writeErrResponse(w, err, http.StatusRequestTimeout, schema.ErrTimedOut)
				return
			}
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
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
			s.writeErrResponse(w, err, http.StatusBadRequest, schema.ErrBadRequest)
			return
		}

		if err := s.eventService.DeleteEvent(ctx, id); err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
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
				s.writeErrResponse(w, err, http.StatusRequestTimeout, schema.ErrTimedOut)
				return
			}
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		eventsResponse, err := json.Marshal(schema.EventsResponse{Events: events})
		if err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
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
			s.writeErrResponse(w, err, http.StatusBadRequest, schema.ErrBadRequest)
			return
		}

		created, err := s.eventService.CreateEvent(ctx, domainEvent)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				s.writeErrResponse(w, err, http.StatusRequestTimeout, schema.ErrTimedOut)
				return
			}
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
			return
		}

		er := schema.EventCreatedResponse{EventID: created}
		resp, err := json.Marshal(er)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write(resp); err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
			return
		}
	}
}

func (s *Server) getIDFromRequest(r *http.Request) (uint, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return 0, fmt.Errorf("invalid id")
	}
	parseInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(parseInt), nil
}

func (s *Server) getEventPayload(r *http.Request) (*domain.Event, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read body")
	}
	var event schema.Event
	if err := json.Unmarshal(body, &event); err != nil {
		return nil, fmt.Errorf("cannot unmarshal body")
	}
	domainEvent, err := codec.HTTPToDomainEvent(&event)
	if err != nil {
		return nil, fmt.Errorf("cannot codec to domain entity")
	}

	return domainEvent, nil
}
