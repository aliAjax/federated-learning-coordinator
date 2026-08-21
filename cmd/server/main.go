package main

import (
	"context"
	grpcapi "github.com/example/federated-learning-coordinator/api/grpc"
	"github.com/example/federated-learning-coordinator/internal/aggregation"
	"github.com/example/federated-learning-coordinator/internal/anomaly"
	"github.com/example/federated-learning-coordinator/internal/config"
	"github.com/example/federated-learning-coordinator/internal/httpserver"
	"github.com/example/federated-learning-coordinator/internal/privacy"
	app "github.com/example/federated-learning-coordinator/internal/round"
	"github.com/example/federated-learning-coordinator/internal/storage"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	store := storage.NewMemory()
	agg := &aggregation.Mean{Noise: privacy.NewGaussian(time.Now().UnixNano()), NoiseScale: cfg.NoiseScale}
	svc := app.NewService(store, agg, anomaly.New(5))
	server := &http.Server{Addr: cfg.HTTPAddr, Handler: httpserver.New(svc, log), ReadHeaderTimeout: 5 * time.Second}
	grpcStop, err := grpcapi.Serve(context.Background(), cfg.GRPCAddr, &grpcapi.API{App: svc})
	if err != nil {
		log.Error("grpc server failed", "error", err)
		return
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		log.Info("http server started", "addr", cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server failed", "error", err)
			stop()
		}
	}()
	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdown); err != nil {
		log.Error("shutdown failed", "error", err)
	}
	grpcStop()
	log.Info("server stopped")
}
