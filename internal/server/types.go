package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kamkali/go-timeline/internal/codec"
	"github.com/kamkali/go-timeline/internal/domain"
	"github.com/kamkali/go-timeline/internal/server/schema"
	"golang.org/x/net/context"
	"io"
	"net/http"
)

func (s *Server) getType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := s.getIDFromRequest(r)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusBadRequest, schema.ErrBadRequest)
			return
		}

		dt, err := s.typeService.GetType(ctx, id)
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

		httpType, err := codec.HTTPFromDomainType(&dt)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
			return
		}
		typeResponse, err := json.Marshal(schema.TypeResponse{Type: httpType})
		if err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
			return
		}
		if _, err := w.Write(typeResponse); err != nil {
			s.log.Error("cannot write response")
			return
		}
	}
}

func (s *Server) updateType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		domainType, err := s.getTypePayload(r)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusBadRequest, schema.ErrBadRequest)
			return
		}

		id, err := s.getIDFromRequest(r)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusBadRequest, schema.ErrBadRequest)
			return
		}

		if err := s.typeService.UpdateType(ctx, id, domainType); err != nil {
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

func (s *Server) deleteType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := s.getIDFromRequest(r)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusBadRequest, schema.ErrBadRequest)
			return
		}

		if err := s.typeService.DeleteType(ctx, id); err != nil {
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

func (s *Server) listTypes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		types, err := s.typeService.ListTypes(ctx)
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
		var httpTypes []*schema.Type
		for i := range types {
			httpType, err := codec.HTTPFromDomainType(&types[i])
			if err != nil {
				s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
				return
			}
			httpTypes = append(httpTypes, httpType)
		}
		typesResponse, err := json.Marshal(schema.TypesResponse{Types: httpTypes})
		if err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
			return
		}
		if _, err := w.Write(typesResponse); err != nil {
			s.log.Error("cannot write response")
			return
		}
	}
}

func (s *Server) createType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		domainType, err := s.getTypePayload(r)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusBadRequest, schema.ErrBadRequest)
			return
		}

		created, err := s.typeService.CreateType(ctx, domainType)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				s.writeErrResponse(w, err, http.StatusRequestTimeout, schema.ErrTimedOut)
				return
			}
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
			return
		}

		er := schema.TypeCreatedResponse{TypeID: created}
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

func (s *Server) getTypePayload(r *http.Request) (*domain.Type, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read body")
	}
	var dt schema.Type
	if err := json.Unmarshal(body, &dt); err != nil {
		return nil, fmt.Errorf("cannot unmarshal body")
	}
	domainType, err := codec.HTTPToDomainType(&dt)
	if err != nil {
		return nil, fmt.Errorf("cannot codec to domain entity")
	}

	return domainType, nil
}
