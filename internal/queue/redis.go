package queue

import (
	"context"
	"encoding/json"
	"go-distributed-worker/internal/config"
	"go-distributed-worker/internal/job"
	"log"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	cfg   *config.Config
	redis *redis.Client
}

func NewRedisClient(cfg *config.Config) (*Client, error) {
	rds := redis.NewClient(&redis.Options{
		Addr: "redis:6379", // docker compose container of redis instance, move to config
	})
	log.Println("Connected to Redis successfully")
	return &Client{
		cfg:   cfg,
		redis: rds,
	}, nil
}

func (c *Client) PeekQueue(queueName string) ([]string, error) {
	return c.redis.LRange(context.Background(), queueName, 0, -1).Result()
}

func (c *Client) PopJob(queueName string) (*job.Job, error) {
	val, err := c.redis.BRPop(context.Background(), 0, queueName).Result()
	if err != nil {
		return nil, err
	}
	var job job.Job
	if err := json.Unmarshal([]byte(val[1]), &job); err != nil {
		return nil, err
	}
	return &job, nil
}
