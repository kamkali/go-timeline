package server

import (
    "encoding/json"
    "errors"
    "github.com/gorilla/mux"
    "github.com/kamkali/go-timeline/internal/codec"
    "github.com/kamkali/go-timeline/internal/config"
    "github.com/kamkali/go-timeline/internal/domain"
    "github.com/kamkali/go-timeline/internal/logger"
    "github.com/kamkali/go-timeline/internal/server/schema"
    "golang.org/x/net/context"
    "io"
    "log"
    "net"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

type Server struct {
    log logger.Logger

    router          *mux.Router
    config          *config.Config
    timelineService domain.TimelineService
    httpServer      *http.Server
}

func New(cfg *config.Config, timelineService domain.TimelineService, log logger.Logger) *Server {
    r := mux.NewRouter()

    s := &Server{
        log:             log,
        config:          cfg,
        timelineService: timelineService,
        router:          r,
        httpServer: &http.Server{
            Addr:    net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
            Handler: r,
        },
    }

    s.registerRoutes()

    return s
}

func (s *Server) registerRoutes() {
    s.router.HandleFunc("/events", s.withTimeout(s.config.Server.TimeoutSeconds, s.listEvents())).Methods("GET")
    s.router.HandleFunc("/events", s.withTimeout(s.config.Server.TimeoutSeconds, s.createEvent())).Methods("POST")
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
    log.Print("Server Exited Properly")
}

func (s *Server) listEvents() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        events, err := s.timelineService.ListEvents(ctx)
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

        body, err := io.ReadAll(r.Body)
        if err != nil {
            s.writeErrResponse(w, err, http.StatusBadRequest, schema.ErrBadRequest)
            return
        }
        var event schema.Event
        if err := json.Unmarshal(body, &event); err != nil {
            s.writeErrResponse(w, err, http.StatusBadRequest, schema.ErrBadRequest)
            return
        }
        domainEvent, err := codec.HTTPToDomainEvent(&event)
        if err != nil {
            s.writeErrResponse(w, err, http.StatusBadRequest, schema.ErrBadRequest)
            return
        }

        created, err := s.timelineService.CreateEvent(ctx, domainEvent)
        if err != nil {
            if errors.Is(err, context.DeadlineExceeded) {
                s.writeErrResponse(w, err, http.StatusRequestTimeout, schema.ErrTimedOut)
                return
            }
            s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
            return
        }

        er := schema.EventResponse{EventID: created}
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

func (s *Server) withTimeout(timeout uint, next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeout)*time.Second)
        defer cancel()
        r = r.WithContext(ctx)
        next.ServeHTTP(w, r)
    }
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
