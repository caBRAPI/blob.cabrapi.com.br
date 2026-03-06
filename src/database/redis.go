package database

import (
	"context"
	"os"
	"time"

	"blob/src/functions"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func Redis() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379/0"
	}

	opt, err := redis.ParseURL(redisURL)
	now := time.Now().Format("Mon Jan 2 15:04:05 2006")

	if err != nil {
		functions.Error("[REDIS ERROR] Invalid URL: %v (%s)", err, now)
		return
	}

	RedisClient = redis.NewClient(opt)

	if err := RedisClient.Ping(Ctx).Err(); err != nil {
		functions.Error("[REDIS ERROR] %v (%s)", err, now)
	} else {
		functions.Info("[REDIS] Connected successfully. (%s)", now)
	}
}
