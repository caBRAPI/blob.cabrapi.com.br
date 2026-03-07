package services

import (
	"blob/src/functions"
	"os"

	"github.com/hibiken/asynq"
)

var AsynqClient *asynq.Client
var AsynqServer *asynq.Server

func InitAsynq() {

	redisURL := os.Getenv("REDIS_URL")

	opt, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		functions.Error("Invalid REDIS_URL: %v", err)
		return
	}

	AsynqClient = asynq.NewClient(opt)

	AsynqServer = asynq.NewServer(
		opt,
		asynq.Config{
			Concurrency: 5,
			Queues: map[string]int{
				"default": 1,
			},
		},
	)

	functions.Info("[ASYNQ] Initialized successfully.")
}
