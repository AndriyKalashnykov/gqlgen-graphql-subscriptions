package datastore

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v9"
)

func NewRedisClient(ctx context.Context, url string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping(ctx).Result()
	if !errors.Is(err, nil) {
		return nil, err
	}

	return client, nil
}
