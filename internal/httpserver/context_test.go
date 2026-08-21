package httpserver

import (
	"context"
	"errors"
	"github.com/example/federated-learning-coordinator/internal/config"
	app "github.com/example/federated-learning-coordinator/internal/round"
	"github.com/example/federated-learning-coordinator/internal/storage"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRequestDeadlineStopsModelRead(t *testing.T) {
	h := (&Server{log: slog.Default()}).middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			w.WriteHeader(http.StatusRequestTimeout)
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	cancel()
	r := httptest.NewRequest(http.MethodGet, "/v1/models/x/download", nil).WithContext(ctx)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, r)
	if rw.Code != http.StatusRequestTimeout && rw.Code != http.StatusInternalServerError {
		t.Fatalf("status %d", rw.Code)
	}
	t.Setenv("GRPC_ADDR", ":29000")
	t.Setenv("NOISE_SCALE", "0.7")
	c := config.Load()
	if c.GRPCAddr != ":29000" || c.NoiseScale != 0.7 {
		t.Fatalf("config lost: %#v", c)
	}
	ctx2, cancel2 := context.WithCancel(context.Background())
	cancel2()
	svc := app.NewService(storage.NewMemory(), nil, nil)
	if _, err := svc.GetModel(ctx2, "x"); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancel lost: %v", err)
	}
	if _, err := svc.Publish(ctx2, "x"); !errors.Is(err, context.Canceled) {
		t.Fatalf("publish cancel lost: %v", err)
	}
	if _, err := svc.Aggregate(ctx2, "x"); !errors.Is(err, context.Canceled) {
		t.Fatalf("aggregate cancel lost: %v", err)
	}
}
