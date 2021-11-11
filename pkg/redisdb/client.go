package redisdb

import (
	"context"
	"sync"
	"time"

	"github.com/kenriortega/ngonx/pkg/errors"
	"github.com/kenriortega/ngonx/pkg/logger"

	"github.com/go-redis/redis/v8"
)

//Used to execute client creation procedure only once.
var redisOnce sync.Once

func GetRedisDbClient(redisUri, redisPass string) *redis.Client {
	var clientInstance *redis.Client
	redisOnce.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:         redisUri,  // use default Addr
			Password:     redisPass, // no password set
			DB:           0,
			DialTimeout:  60 * time.Second,
			ReadTimeout:  60 * time.Second,
			WriteTimeout: 60 * time.Second,
		})

		_, err := clientInstance.Ping(context.TODO()).Result()
		if err != nil {
			logger.LogError(errors.Errorf("redis: %v", err).Error())

		}
		clientInstance = client
	})
	return clientInstance
}
