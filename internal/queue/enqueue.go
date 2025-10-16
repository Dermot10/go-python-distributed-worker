package queue

import (
	"encoding/json"
	"fmt"
	"go-distributed-worker/internal/job"
	"log"
	"net/http"
)

// handler pushing tasks to redis, event producer from incoming requests

func EnqueueHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse, validate, push, repsonse

		var job job.Job
		if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		dataValidator(job, w)

		data, _ := json.Marshal(job)
		if err := c.redis.LPush(r.Context(), "job_queue", data).Err(); err != nil {
			http.Error(w, "failed to enqueue job", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, "job %s enqueued successfully", job.ID)
		log.Printf("Job enqueued: %s", job.ID)

	}
}

func dataValidator(job job.Job, w http.ResponseWriter) {
	if job.ID == "" {
		http.Error(w, "missing job ID", http.StatusBadRequest)
		return
	}

	if job.Type == "" {
		http.Error(w, "missing job type", http.StatusBadRequest)
		return
	}

	if len(job.Payload) == 0 {
		http.Error(w, "missing job type", http.StatusBadRequest)
		return
	}

	if job.CreatedAt.IsZero() {
		http.Error(w, "timestamp is missing", http.StatusBadRequest)
		return
	}
}
