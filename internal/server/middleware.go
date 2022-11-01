package server

import (
	"fmt"
	"github.com/kamkali/go-timeline/internal/server/schema"
	"golang.org/x/net/context"
	"net/http"
	"strings"
	"time"
)

func (s *Server) withTimeout(timeout uint, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeout)*time.Second)
		defer cancel()
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}

func (s *Server) withAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			s.writeErrResponse(w, fmt.Errorf("missing Authorization Header"), http.StatusUnauthorized, schema.ErrUnauthorized)
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		verified, err := s.jwtManager.VerifyToken(tokenString)
		if err != nil {
			s.writeErrResponse(w, fmt.Errorf("token veryfication: %w", err), http.StatusUnauthorized, schema.ErrUnauthorized)
			return
		}

		if !verified {
			s.writeErrResponse(w, fmt.Errorf("token not veryfied"), http.StatusUnauthorized, schema.ErrUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
