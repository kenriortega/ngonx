package redisdb

import (
	"context"
	"os"
	"time"

	"github.com/kenriortega/ngonx/pkg/logger"

	"github.com/go-redis/redis/v8"
)

func GetRedisDbClient() *redis.Client {

	clientInstance := redis.NewClient(&redis.Options{
		Addr:         os.Getenv("REDIS_URI"),  // use default Addr
		Password:     os.Getenv("REDIS_PASS"), // no password set
		DB:           0,
		DialTimeout:  60 * time.Second,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	})

	_, err := clientInstance.Ping(context.TODO()).Result()
	if err != nil {
		logger.LogError(err.Error())
	}
	return clientInstance
}
