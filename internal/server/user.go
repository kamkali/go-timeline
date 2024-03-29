package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kamkali/go-timeline/internal/codec"
	schema2 "github.com/kamkali/go-timeline/internal/server/schema"
	timeline2 "github.com/kamkali/go-timeline/internal/timeline"
	"io"
	"net/http"
	"strings"
	"time"
)

func (s *Server) changePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			s.writeErrResponse(w, fmt.Errorf("missing Authorization Header"), http.StatusUnauthorized, schema2.ErrUnauthorized)
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		newPassword, err := s.getPasswordPayload(r)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusBadRequest, schema2.ErrBadRequest)
			return
		}
		claims, err := s.jwtManager.GetClaims(tokenString)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusUnauthorized, schema2.ErrUnauthorized)
			return
		}
		user, ok := claims["user"]
		if !ok {
			s.writeErrResponse(w, err, http.StatusBadRequest, schema2.ErrBadRequest)
			return
		}
		username, ok := user.(string)
		if !ok {
			s.writeErrResponse(w, err, http.StatusBadRequest, schema2.ErrBadRequest)
			return
		}

		if err := s.userService.ChangePassword(ctx, username, newPassword); err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema2.ErrInternal)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, err := s.getUserPayload(r)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusBadRequest, schema2.ErrBadRequest)
			return
		}
		loggedUser, err := s.userService.LoginUser(ctx, user)
		if err != nil {
			if errors.Is(err, timeline2.ErrUnauthorized) || errors.Is(err, timeline2.ErrNotFound) {
				s.writeErrResponse(w, err, http.StatusUnauthorized, schema2.ErrUnauthorized)
				return
			}
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema2.ErrInternal)
			return
		}

		token, err := s.jwtManager.GenerateToken(loggedUser.Email)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema2.ErrInternal)
			return
		}
		claims, err := s.jwtManager.GetClaims(token)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema2.ErrInternal)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		cookie := http.Cookie{
			Name:    "Authorization",
			Value:   "Bearer " + token,
			Expires: time.Unix(int64(claims["exp"].(float64)), 0),
		}
		if err := cookie.Valid(); err != nil {
			s.log.Error(err.Error())
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema2.ErrInternal)
			return
		}

		w.WriteHeader(http.StatusOK)
		tokenResponse, err := json.Marshal(schema2.TokenResponse{Token: token})
		if err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema2.ErrInternal)
			return
		}
		http.SetCookie(w, &cookie)
		if _, err := w.Write(tokenResponse); err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema2.ErrInternal)
			return
		}
	}
}

func (s *Server) getUserPayload(r *http.Request) (*timeline2.User, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read body")
	}
	var user schema2.User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("cannot unmarshal body")
	}
	domainUser, err := codec.HTTPToDomainUser(&user)
	if err != nil {
		return nil, fmt.Errorf("cannot codec to domain entity")
	}

	return domainUser, nil
}

func (s *Server) getPasswordPayload(r *http.Request) (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read body")
	}
	var pass schema2.PasswordChange
	if err := json.Unmarshal(body, &pass); err != nil {
		return "", fmt.Errorf("cannot unmarshal body")
	}

	return pass.NewPassword, nil
}
