package worker

import (
	"context"
	"go-distributed-worker/internal/job"
	"go-distributed-worker/internal/queue"
	"log"

	"golang.org/x/sync/errgroup"
)

// consume from redis queue, with brpop and processes

func worker(ctx context.Context, id int, rdb *queue.Client) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d: shutting down\n", id)
			return
		default:
			job, err := rdb.PopJob("job_queue")
			if err != nil {
				log.Printf("Worker %d error: %v\n", id, err)
				continue
			}
			process(job)
		}
	}

}

func runWorkerPool(ctx context.Context, rdb *queue.Client, numWorkers int) error {
	g, ctx := errgroup.WithContext(ctx)
	for i := 1; i <= numWorkers; i++ {
		workerID := i
		g.Go(func() error {
			worker(ctx, workerID, rdb)
			return nil
		})
	}
	return g.Wait()
}

func process(job *job.Job) {
	log.Printf("Processing job %s of type %s", job.ID, job.Type)
	// Add actual processing logic here
}
