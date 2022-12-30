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

func (s *Server) getProcess() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := s.getIDFromRequest(r)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusBadRequest, schema.ErrBadRequest)
			return
		}

		process, err := s.processService.GetProcess(ctx, id)
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

		httpProcess, err := codec.HTTPFromDomainProcess(&process)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
			return
		}
		processResponse, err := json.Marshal(schema.ProcessResponse{Process: httpProcess})
		if err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
			return
		}
		if _, err := w.Write(processResponse); err != nil {
			s.log.Error("cannot write response")
			return
		}
	}
}

func (s *Server) updateProcess() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		domainProcess, err := s.getProcessPayload(r)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusBadRequest, schema.ErrBadRequest)
			return
		}

		id, err := s.getIDFromRequest(r)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusBadRequest, schema.ErrBadRequest)
			return
		}

		if err := s.processService.UpdateProcess(ctx, id, domainProcess); err != nil {
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

func (s *Server) deleteProcess() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := s.getIDFromRequest(r)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusBadRequest, schema.ErrBadRequest)
			return
		}

		if err := s.processService.DeleteProcess(ctx, id); err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) listProcesses() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		processes, err := s.processService.ListProcesses(ctx)
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

		var httpProcesses []*schema.Process
		for i := range processes {
			process, err := codec.HTTPFromDomainProcess(&processes[i])
			if err != nil {
				s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
				return
			}
			httpProcesses = append(httpProcesses, process)
		}
		processesResponse, err := json.Marshal(schema.ProcessesResponse{Processes: httpProcesses})
		if err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
			return
		}
		if _, err := w.Write(processesResponse); err != nil {
			s.log.Error("cannot write response")
			return
		}
	}
}

func (s *Server) createProcess() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		domainProcess, err := s.getProcessPayload(r)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusBadRequest, schema.ErrBadRequest)
			return
		}

		created, err := s.processService.CreateProcess(ctx, domainProcess)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				s.writeErrResponse(w, err, http.StatusRequestTimeout, schema.ErrTimedOut)
				return
			}
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
			return
		}

		er := schema.ProcessCreatedResponse{ProcessID: created}
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

func (s *Server) getProcessPayload(r *http.Request) (*domain.Process, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read body")
	}
	var process schema.Process
	if err := json.Unmarshal(body, &process); err != nil {
		return nil, fmt.Errorf("cannot unmarshal body")
	}
	domainProcess, err := codec.HTTPToDomainProcess(&process)
	if err != nil {
		return nil, fmt.Errorf("cannot codec to domain entity")
	}

	return domainProcess, nil
}
