package service

import (
	"context"
	"fmt"
	"go-distributed-worker/internal/config"
	"go-distributed-worker/internal/job"
	"log"
)

type Service struct {
	cfg *config.Config
	log *log.Logger
	ws  WorkerService
	QS  QueueService
}

type WorkerService interface {
	RunWorkerPool(ctx context.Context, qs QueueService, numWorkers int) error
}

type QueueService interface {
	PushJob(ctx context.Context, queueName string, job *job.Job) error
	PopJob(queueName string) (*job.Job, error)
	PeekQueue(queueName string) ([]string, error)
}

func NewService(cfg *config.Config, log *log.Logger, qs QueueService, ws WorkerService) *Service {
	return &Service{
		cfg: cfg,
		log: log,
		ws:  ws,
		QS:  qs,
	}
}
func (s *Service) Execute(ctx context.Context) error {
	if err := s.runService(ctx); err != nil {
		s.log.Printf("service execution failed: %v", err)
		return err
	}
	return nil
}

func (s *Service) runService(ctx context.Context) error {

	if err := s.viewQueue(ctx, s.cfg.JobQueueName); err != nil {
		return fmt.Errorf("failed to view job queue: %w", err)
	}

	if err := s.ws.RunWorkerPool(ctx, s.QS, s.cfg.NumWorkers); err != nil {
		return fmt.Errorf("worker pool failed: %w", err)
	}

	return nil
}

func (s *Service) viewQueue(ctx context.Context, queueName string) error {

	vals, err := s.QS.PeekQueue(queueName)
	if err != nil {
		return err
	}

	for i, v := range vals {
		fmt.Printf("%d: %s\n", i, v)
	}
	return nil
}
