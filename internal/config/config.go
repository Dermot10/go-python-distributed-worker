package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	RedisAddr    string `env:"REDIS_ADDR,required"`
	ServicePort  string `env:"SERVICE_PORT" envDefault:"8080"`
	NumWorkers   int    `env:"WORKER_COUNT" envDefault:"5"`
	JobQueueName string `env:"JOB_QUEUE_NAME" envDefault:"job_queue"`
}

func LoadConfig() (*Config, error) {
	props := Config{}
	if err := env.Parse(&props); err != nil {
		return nil, fmt.Errorf("error loading properties: %w", err)
	}
	return &props, nil
}
