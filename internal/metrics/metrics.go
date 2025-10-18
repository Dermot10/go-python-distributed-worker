package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	prometheus.MustRegister(
		ProcessedRequests,
		FailedRequests,
	)
}

var (
	ProcessedRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "processed_request_total",
			Help: "Total number of processed requests from concurrent workers",
		},
	)
	FailedRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "failed_request_total",
			Help: "Total number of requests failed to be processed",
		},
	)
)
