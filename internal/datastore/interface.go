package datastore

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// RedisClient defines the interface for Redis operations used in this application
type RedisClient interface {
	XAdd(ctx context.Context, args *redis.XAddArgs) *redis.StringCmd
	XRead(ctx context.Context, args *redis.XReadArgs) *redis.XStreamSliceCmd
	Ping(ctx context.Context) *redis.StatusCmd
	Close() error
}
