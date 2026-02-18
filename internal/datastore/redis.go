package datastore

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, url string) (RedisClient, error) {
	if url == "" {
		return nil, fmt.Errorf("redis URL cannot be empty")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping(ctx).Result()
	if !errors.Is(err, nil) {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}
