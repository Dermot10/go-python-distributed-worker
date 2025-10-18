package queue

import (
	"context"
	"encoding/json"
	"go-distributed-worker/internal/config"
	"go-distributed-worker/internal/job"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisQueueClient struct {
	Rdb *redis.Client
}

func NewQueueService(cfg *config.Config) (*RedisQueueClient, error) {
	rdb := redis.NewClient(&redis.Options{
		// Addr: cfg.RedisAddr,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	log.Print("Successfully connected to redis")
	return &RedisQueueClient{Rdb: rdb}, nil
}

func (r *RedisQueueClient) PeekQueue(queueName string) ([]string, error) {
	return r.Rdb.LRange(context.Background(), queueName, 0, -1).Result()
}

func (r *RedisQueueClient) PopJob(queueName string) (*job.Job, error) {
	val, err := r.Rdb.BRPop(context.Background(), 0, queueName).Result()
	if err != nil {
		return nil, err
	}
	var job job.Job
	if err := json.Unmarshal([]byte(val[1]), &job); err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *RedisQueueClient) PushJob(ctx context.Context, queueName string, job job.Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return r.Rdb.LPush(ctx, queueName, data).Err()
}
