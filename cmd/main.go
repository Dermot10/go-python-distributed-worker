package main

import (
	"context"
	"fmt"
	"go-distributed-worker/internal/config"
	"log"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {

	// load from config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// logger

	// broader application ctx for goroutine lifecycle, both service and health endpoints
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// instantiate redis, pass to service
	mux := setUpMux()
	rdb, err := setUpRedis(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Register all routes
	registerRoutes(mux, rdb)

	vals, err := rdb.PeekQueue("job_queue")
	if err != nil {
		log.Fatal(err)
	}

	for i, v := range vals {
		fmt.Printf("%d: %s\n", i, v)
	}

	g, ctx := errgroup.WithContext(ctx)

	// set up service, pass to goroutine

	g.Go(func() error {

		return nil
	})

	g.Go(func() error {
		return serveHealthAndMetrics(ctx, mux)
	})

	if err := g.Wait(); err != nil {
		// log service error
		log.Fatal("service encountered unexpected failure: %w", err)
	}
}
