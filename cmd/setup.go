package main

import (
	"context"
	"fmt"
	"go-distributed-worker/internal/config"
	"go-distributed-worker/internal/queue"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Dependencies struct { //needed for service
}

func setUpMux() *http.ServeMux {
	mux := http.NewServeMux()
	return mux
}

func registerRoutes(mux *http.ServeMux, rdb *queue.Client) {
	mux.HandleFunc("/enqueue", queue.EnqueueHandler(rdb))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.Handle("/metrics", promhttp.Handler()) // server interface satisfied by promhttp handler
}

func setUpRedis(cfg *config.Config) (*queue.Client, error) {
	rdb, err := queue.NewRedisClient(cfg)
	if err != nil {
		return nil, err
	}
	return rdb, nil
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
