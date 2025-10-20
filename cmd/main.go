package main

import (
	"context"
	"go-distributed-worker/internal/config"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {

	log := log.New(os.Stdout, "distributed-worker ", log.LstdFlags)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// broader application ctx for goroutine lifecycle, both service and health endpoints
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	s, err := setUpServiceRunner(ctx, log, cfg)
	if err != nil {
		log.Fatalln("failed to set up distributed worker service", err)
	}
	// instantiate redis, pass to service
	mux := setUpMux()

	// Register all routes
	registerRoutes(mux, s.QS)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return s.Execute(ctx)
	})

	g.Go(func() error {
		return serveHealthAndMetrics(ctx, mux)
	})

	if err := g.Wait(); err != nil {
		// log service error
		log.Fatal("service encountered unexpected failure: %w", err)
	}
}
