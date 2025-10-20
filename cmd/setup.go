package main

import (
	"context"
	"fmt"
	"go-distributed-worker/internal/config"
	"go-distributed-worker/internal/handler"
	"go-distributed-worker/internal/queue"
	"go-distributed-worker/internal/service"
	"go-distributed-worker/internal/worker"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Dependencies struct { //needed for service

}

func setUpServiceRunner(ctx context.Context, log *log.Logger, cfg *config.Config) (*service.Service, error) {

	qs, err := setUpQueueService(cfg)
	if err != nil {
		return nil, err
	}

	ws, err := setUpWorkerService(cfg)
	if err != nil {
		return nil, err
	}

	return service.NewService(cfg, log, qs, ws), nil
}

func setUpQueueService(cfg *config.Config) (*queue.RedisQueueClient, error) {
	qs, err := queue.NewQueueService(cfg)
	if err != nil {
		return nil, err
	}
	return qs, nil
}

func setUpWorkerService(cfg *config.Config) (*worker.Worker, error) {
	ws, err := worker.NewWorkerService(cfg)
	if err != nil {
		return nil, err
	}
	return ws, nil
}

func setUpMux() *http.ServeMux {
	mux := http.NewServeMux()
	return mux
}

func registerRoutes(mux *http.ServeMux, qs service.QueueService) {
	mux.HandleFunc("/enqueue", handler.EnqueueHandler(qs))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.Handle("/metrics", promhttp.Handler()) // server interface satisfied by promhttp handler
}

func serveHealthAndMetrics(ctx context.Context, mux *http.ServeMux) error {
	// lightweight server serves both the health and metrics endpoints

	server := &http.Server{
		Addr:              ":" + "8080", // cfg.UtilPortconfig, for the utility port
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() error {
		log.Println("Server running on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("server error: %w", err)
		}
		return nil
	}()

	// wait on ctx cancellation for a shutdown otherwise just run goroutine
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}
	return nil
}
