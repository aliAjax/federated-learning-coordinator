package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/example/federated-learning-coordinator/internal/model"
	app "github.com/example/federated-learning-coordinator/internal/round"
	"github.com/example/federated-learning-coordinator/internal/storage"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"
)

type Server struct {
	app      *app.Service
	log      *slog.Logger
	requests atomic.Uint64
	failures atomic.Uint64
}

func New(a *app.Service, l *slog.Logger) http.Handler {
	s := &Server{app: a, log: l}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.health)
	mux.HandleFunc("GET /readyz", s.ready)
	mux.HandleFunc("GET /metrics", s.metrics)
	mux.HandleFunc("POST /v1/cohorts", s.createCohort)
	mux.HandleFunc("POST /v1/participants", s.createParticipant)
	mux.HandleFunc("POST /v1/rounds", s.createRound)
	mux.HandleFunc("GET /v1/rounds/{id}", s.getRound)
	mux.HandleFunc("POST /v1/rounds/{id}/abort", s.abortRound)
	mux.HandleFunc("POST /v1/rounds/{id}/updates", s.submitUpdate)
	mux.HandleFunc("POST /v1/rounds/{id}/aggregate", s.aggregate)
	mux.HandleFunc("POST /v1/models/{id}/publish", s.publish)
	mux.HandleFunc("GET /v1/models/{id}/download", s.download)
	mux.HandleFunc("GET /v1/privacy-budgets", s.budgets)
	return s.middleware(mux)
}
func (s *Server) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.requests.Add(1)
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = app.ID("req")
		}
		w.Header().Set("X-Request-ID", id)
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Type", "application/json")
		ctx, cancel := requestContext(r)
		defer cancel()
		defer func() {
			if v := recover(); v != nil {
				s.failures.Add(1)
				s.log.Error("panic recovered", "request_id", id, "panic", v, "stack", string(debug.Stack()))
				writeError(w, http.StatusInternalServerError, "internal_error", "unexpected error")
			}
		}()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func requestContext(r *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), 15*time.Second)
}
func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 8<<20)
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(v); err != nil {
		writeError(w, 400, "invalid_request", err.Error())
		return false
	}
	return true
}
func write(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeError(w http.ResponseWriter, status int, code, msg string) {
	write(w, status, map[string]any{"error": map[string]string{"code": code, "message": msg}})
}
func (s *Server) respond(w http.ResponseWriter, v any, err error, created bool) {
	if err == nil {
		if created {
			write(w, 201, v)
		} else {
			write(w, 200, v)
		}
		return
	}
	s.failures.Add(1)
	status := 400
	code := "invalid_request"
	if errors.Is(err, storage.ErrNotFound) {
		status = 404
		code = "not_found"
	}
	if errors.Is(err, storage.ErrConflict) {
		status = 409
		code = "conflict"
	}
	if errors.Is(err, app.ErrInvalidState) {
		status = 409
		code = "invalid_state"
	}
	writeError(w, status, code, err.Error())
}
func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	write(w, 200, map[string]string{"status": "ok"})
}
func (s *Server) ready(w http.ResponseWriter, _ *http.Request) {
	write(w, 200, map[string]string{"status": "ready"})
}
func (s *Server) metrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	_, _ = w.Write([]byte("federated_http_requests_total " + itoa(s.requests.Load()) + "\nfederated_http_failures_total " + itoa(s.failures.Load()) + "\n"))
}
func (s *Server) createCohort(w http.ResponseWriter, r *http.Request) {
	var q model.CohortRequest
	if !decode(w, r, &q) {
		return
	}
	v, e := s.app.CreateCohort(r.Context(), q)
	s.respond(w, v, e, true)
}
func (s *Server) createParticipant(w http.ResponseWriter, r *http.Request) {
	var q model.ParticipantRequest
	if !decode(w, r, &q) {
		return
	}
	v, e := s.app.Register(r.Context(), q)
	s.respond(w, v, e, true)
}
func (s *Server) createRound(w http.ResponseWriter, r *http.Request) {
	var q model.RoundRequest
	if !decode(w, r, &q) {
		return
	}
	v, e := s.app.CreateRound(r.Context(), q)
	s.respond(w, v, e, true)
}
func (s *Server) getRound(w http.ResponseWriter, r *http.Request) {
	v, e := s.app.GetRound(r.Context(), r.PathValue("id"))
	s.respond(w, v, e, false)
}
func (s *Server) abortRound(w http.ResponseWriter, r *http.Request) {
	v, e := s.app.Abort(r.Context(), r.PathValue("id"))
	s.respond(w, v, e, false)
}
func (s *Server) submitUpdate(w http.ResponseWriter, r *http.Request) {
	var q model.UpdateRequest
	if !decode(w, r, &q) {
		return
	}
	v, e := s.app.Submit(r.Context(), r.PathValue("id"), q)
	s.respond(w, v, e, true)
}
func (s *Server) aggregate(w http.ResponseWriter, r *http.Request) {
	v, e := s.app.Aggregate(r.Context(), r.PathValue("id"))
	s.respond(w, v, e, false)
}
func (s *Server) publish(w http.ResponseWriter, r *http.Request) {
	v, e := s.app.Publish(r.Context(), r.PathValue("id"))
	s.respond(w, v, e, false)
}
func (s *Server) download(w http.ResponseWriter, r *http.Request) {
	v, e := s.app.GetModel(r.Context(), r.PathValue("id"))
	s.respond(w, v, e, false)
}
func (s *Server) budgets(w http.ResponseWriter, r *http.Request) {
	write(w, 200, map[string]any{"items": s.app.Budgets(r.Context())})
}
func itoa(v uint64) string {
	if v == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for v > 0 {
		i--
		b[i] = byte('0' + v%10)
		v /= 10
	}
	return string(b[i:])
}

var _ = strings.TrimSpace
