# Distributed Worker Service

This project implements a Go-based distributed worker service with Redis as a queue backend and a Python client for testing. It simulates concurrent work processing in a scalable and observable manner.

# Overview

A Python script simulates network traffic by firing a configurable number of job requests to the Go service:

Go service: Accepts HTTP requests, validates and enqueues jobs, then processes them asynchronously using a worker pool.

Redis: Acts as a messaging queue between the HTTP client and worker goroutines.

Worker pool: Concurrent workers consume jobs from Redis, perform processing logic, and update Prometheus metrics.

Metrics: Exposed via /metrics endpoint (Prometheus counters for processed and failed jobs).

Health: /healthz endpoint reports service availability.

How It Works
The client sends a job via HTTP POST to /enqueue.

    The Go service:

        Parses and validates the job.

        Pushes it onto the Redis queue.

        Returns 202 Accepted to the client.

        Workers:

            Pop jobs from Redis.

            Process the job logic.

        Increment Prometheus counters (processed_request_total, failed_request_total).

# Design Considerations

Concurrent job processing with a configurable worker pool.

Prometheus metrics for observability.

Graceful shutdown via context

Dockerized setup with Redis + Go service.

Structured logging with INFO/ERROR levels.

# Docker Setup

docker-compose up --build

Go service listens on :8081 (can be configured).

Redis runs on port 6379.

Metrics available at http://localhost:8081/metrics.

# Testing

Use the Python client to enqueue jobs.

Monitor logs via:

docker-compose logs -f go-service

Unittests to be added
