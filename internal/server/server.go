package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kamkali/go-timeline/internal/config"
	"github.com/kamkali/go-timeline/internal/domain"
	"github.com/kamkali/go-timeline/internal/logger"
	"github.com/kamkali/go-timeline/internal/server/schema"
	"golang.org/x/net/context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type Server struct {
	router     *mux.Router
	config     *config.Config
	httpServer *http.Server

	log          logger.Logger
	eventService domain.EventService
	typeService  domain.TypeService
}

func New(cfg *config.Config, log logger.Logger, eventService domain.EventService, typesService domain.TypeService) *Server {
	r := mux.NewRouter()

	s := &Server{
		router: r,
		config: cfg,
		httpServer: &http.Server{
			Addr:    net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
			Handler: r,
		},
		log:          log,
		eventService: eventService,
		typeService:  typesService,
	}

	s.registerRoutes()

	return s
}

func (s *Server) registerRoutes() {
	{ // Events routes
		s.router.HandleFunc("/events",
			s.withTimeout(s.config.Server.TimeoutSeconds, s.listEvents()),
		).Methods("GET")

		s.router.HandleFunc("/events/{id}",
			s.withTimeout(s.config.Server.TimeoutSeconds, s.getEvent()),
		).Methods("GET")

		s.router.HandleFunc("/events/{id}",
			s.withTimeout(s.config.Server.TimeoutSeconds, s.updateEvent()),
		).Methods("PUT")

		s.router.HandleFunc("/events/{id}",
			s.withTimeout(s.config.Server.TimeoutSeconds, s.deleteEvent()),
		).Methods("DELETE")

		s.router.HandleFunc("/events",
			s.withTimeout(s.config.Server.TimeoutSeconds, s.createEvent()),
		).Methods("POST")
	}

	{ // Types routes
		s.router.HandleFunc("/types",
			s.withTimeout(s.config.Server.TimeoutSeconds, s.listTypes()),
		).Methods("GET")

		s.router.HandleFunc("/types/{id}",
			s.withTimeout(s.config.Server.TimeoutSeconds, s.getType()),
		).Methods("GET")

		s.router.HandleFunc("/types/{id}",
			s.withTimeout(s.config.Server.TimeoutSeconds, s.updateType()),
		).Methods("PUT")

		s.router.HandleFunc("/types/{id}",
			s.withTimeout(s.config.Server.TimeoutSeconds, s.deleteType()),
		).Methods("DELETE")

		s.router.HandleFunc("/types",
			s.withTimeout(s.config.Server.TimeoutSeconds, s.createType()),
		).Methods("POST")
	}
}

func (s *Server) Start() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	s.log.Info("Server Started")

	<-done
	s.log.Info("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	s.log.Info("Server Exited Properly")
}

func (s *Server) writeErrResponse(w http.ResponseWriter, err error, code int, desc string) {
	s.log.Errorf("error response: %s", err.Error())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonErr, err := json.Marshal(schema.ServerError{Description: desc})
	if err != nil {
		return
	}
	if _, err := w.Write(jsonErr); err != nil {
		s.log.Errorf("cannot write error response: %w", err.Error())
		return
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
