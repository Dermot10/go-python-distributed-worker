package handler

import (
	"encoding/json"
	"fmt"
	"go-distributed-worker/internal/job"
	"go-distributed-worker/internal/service"
	"log"
	"net/http"
)

// handler pushing tasks to redis, event producer from incoming requests
func EnqueueHandler(qs service.QueueService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var j job.Job
		if err := json.NewDecoder(r.Body).Decode(&j); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if err := validateJob(j); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := qs.PushJob(r.Context(), "job_queue", &j); err != nil {
			http.Error(w, "failed to enqueue job", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, "job %s enqueued successfully", j.ID)
		log.Printf("Job enqueued: %s", j.ID)
	}
}

func validateJob(job job.Job) error {
	if job.ID == "" {
		return fmt.Errorf("missing job ID")
	}
	if job.Type == "" {
		return fmt.Errorf("missing job type")
	}
	if len(job.Payload) == 0 {
		return fmt.Errorf("missing job payload")
	}
	if job.CreatedAt.IsZero() {
		return fmt.Errorf("timestamp is missing")
	}
	return nil
}
