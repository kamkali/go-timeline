package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kamkali/go-timeline/internal/codec"
	"github.com/kamkali/go-timeline/internal/domain"
	"github.com/kamkali/go-timeline/internal/server/schema"
	"io"
	"net/http"
)

func (s Server) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, err := s.getUserPayload(r)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusBadRequest, schema.ErrBadRequest)
			return
		}
		loggedUser, err := s.userService.LoginUser(ctx, user)
		if err != nil {
			if errors.Is(err, domain.ErrUnauthorized) || errors.Is(err, domain.ErrNotFound) {
				s.writeErrResponse(w, err, http.StatusUnauthorized, schema.ErrUnauthorized)
			}
			return
		}

		token, err := s.jwtManager.GenerateToken(loggedUser.Email)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Authorization", "Bearer "+token)
		w.WriteHeader(http.StatusOK)
		tokenResponse, err := json.Marshal(schema.TokenResponse{Token: token})
		if err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
			return
		}
		if _, err := w.Write(tokenResponse); err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
			return
		}
	}
}

func (s *Server) getUserPayload(r *http.Request) (*domain.User, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read body")
	}
	var user schema.User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("cannot unmarshal body")
	}
	domainUser, err := codec.HTTPToDomainUser(&user)
	if err != nil {
		return nil, fmt.Errorf("cannot codec to domain entity")
	}

	return domainUser, nil
}
