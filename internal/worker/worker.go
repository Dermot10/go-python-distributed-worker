package worker

import (
	"context"
	"go-distributed-worker/internal/job"
	"go-distributed-worker/internal/metrics"
	"go-distributed-worker/internal/service"
	"log"

	"golang.org/x/sync/errgroup"
)

// consume from redis queue, with brpop and processes
// worker pool spawns workers up to limit
// each worker will continually look for work until stopped or no work is left
// process is dummy function currently only printing

type Worker struct {
}

func NewWorkerService() (*Worker, error) {
	return &Worker{}, nil
}

func (w *Worker) RunWorkerPool(ctx context.Context, qs service.QueueService, numWorkers int) error {
	g, ctx := errgroup.WithContext(ctx)
	for i := 1; i <= numWorkers; i++ {
		workerID := i
		g.Go(func() error {
			w.worker(ctx, workerID, qs)
			return nil
		})
	}
	return g.Wait()
}

func (w *Worker) worker(ctx context.Context, id int, qs service.QueueService) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d: shutting down\n", id)
			return
		default:
			job, err := qs.PopJob("job_queue")
			if err != nil {
				log.Printf("Worker %d error: %v\n", id, err)
				continue
			}
			w.process(job)
		}
	}

}

func (w *Worker) process(job *job.Job) {
	metrics.ProcessedRequests.Inc()
	log.Printf("job processed successfully: %s", job.ID)

	// simulate some kind of work here, failures and success can be logged here or consider service level again
}
