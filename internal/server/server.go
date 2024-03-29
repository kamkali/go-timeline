package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kamkali/go-timeline/internal/auth"
	"github.com/kamkali/go-timeline/internal/config"
	"github.com/kamkali/go-timeline/internal/generator"
	"github.com/kamkali/go-timeline/internal/server/schema"
	"github.com/kamkali/go-timeline/internal/timeline"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

//go:embed static
var staticFS embed.FS

type Server struct {
	router       *mux.Router
	config       *config.Config
	httpServer   *http.Server
	staticServer http.Handler

	log        *zap.Logger
	jwtManager *auth.JWTManager

	eventService timeline.EventService
	typeService  timeline.TypeService
	userService  timeline.UserService
	renderer     *generator.Renderer
}

func New(
	cfg *config.Config,
	log *zap.Logger,
	manager *auth.JWTManager,
	eventService timeline.EventService,
	typesService timeline.TypeService,
	userService timeline.UserService,
) (*Server, error) {
	r := mux.NewRouter()
	siteRenderer, err := generator.NewRenderer()
	if err != nil {
		return nil, fmt.Errorf("cannot instantiate site renderer: %w", err)
	}
	handler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000", "https://apollo11timeline.herokuapp.com"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "DELETE", "PUT", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Origin", "Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)(r)
	s := &Server{
		router: r,
		config: cfg,
		httpServer: &http.Server{
			Addr:    net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
			Handler: handler,
		},
		log:          log,
		jwtManager:   manager,
		eventService: eventService,
		typeService:  typesService,
		userService:  userService,
		renderer:     siteRenderer,
	}

	fSys, err := fs.Sub(staticFS, "static")
	if err != nil {
		return nil, err
	}
	s.staticServer = http.FileServer(http.FS(fSys))

	s.registerRoutes()

	return s, nil
}

func (s *Server) registerRoutes() {
	{ // public routes
		s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", s.staticServer))
		s.router.HandleFunc("/", s.renderTimeline()).Methods("GET")
	}

	{ // Events routes
		s.router.HandleFunc("/api/events",
			s.withTimeout(s.config.Server.TimeoutSeconds, s.listEvents()),
		).Methods("GET")

		s.router.HandleFunc("/api/events/{id}",
			s.withTimeout(s.config.Server.TimeoutSeconds, s.getEvent()),
		).Methods("GET")

		s.router.HandleFunc("/api/events/{id}",
			s.withAuth(s.withTimeout(s.config.Server.TimeoutSeconds, s.updateEvent())),
		).Methods("PUT")

		s.router.HandleFunc("/api/events/{id}",
			s.withAuth(s.withTimeout(s.config.Server.TimeoutSeconds, s.deleteEvent())),
		).Methods("DELETE")

		s.router.HandleFunc("/api/events",
			s.withAuth(s.withTimeout(s.config.Server.TimeoutSeconds, s.createEvent())),
		).Methods("POST")
	}

	{ // Types routes
		s.router.HandleFunc("/api/types",
			s.withTimeout(s.config.Server.TimeoutSeconds, s.listTypes()),
		).Methods("GET")

		s.router.HandleFunc("/api/types/{id}",
			s.withTimeout(s.config.Server.TimeoutSeconds, s.getType()),
		).Methods("GET")

		s.router.HandleFunc("/api/types/{id}",
			s.withAuth(s.withTimeout(s.config.Server.TimeoutSeconds, s.updateType())),
		).Methods("PUT")

		s.router.HandleFunc("/api/types/{id}",
			s.withAuth(s.withTimeout(s.config.Server.TimeoutSeconds, s.deleteType())),
		).Methods("DELETE")

		s.router.HandleFunc("/api/types",
			s.withAuth(s.withTimeout(s.config.Server.TimeoutSeconds, s.createType())),
		).Methods("POST")
	}

	{ // User routes
		s.router.HandleFunc("/api/login",
			s.withTimeout(s.config.Server.TimeoutSeconds, s.login()),
		).Methods("POST")

		s.router.HandleFunc("/api/change_password",
			s.withAuth(s.withTimeout(s.config.Server.TimeoutSeconds, s.changePassword())),
		).Methods("POST")

		s.router.HandleFunc("/api/check",
			s.withAuth(s.withTimeout(s.config.Server.TimeoutSeconds, s.check())),
		).Methods("GET")
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
	s.log.Info(fmt.Sprintf("Server Started on host=%s:%s", s.config.Server.Host, s.config.Server.Port))

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
	s.log.Info(fmt.Errorf("error response: %w", err).Error())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonErr, err := json.Marshal(schema.ServerError{Description: desc})
	if err != nil {
		return
	}
	if _, err := w.Write(jsonErr); err != nil {
		s.log.Error(fmt.Errorf("cannot write error response: %w", err).Error())
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

func (s *Server) check() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}
}
