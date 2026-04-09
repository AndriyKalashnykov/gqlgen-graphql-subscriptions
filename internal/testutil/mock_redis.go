package testutil

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// MockRedisClient is a test double for datastore.RedisClient
type MockRedisClient struct {
	XAddFunc  func(ctx context.Context, args *redis.XAddArgs) *redis.StringCmd
	XReadFunc func(ctx context.Context, args *redis.XReadArgs) *redis.XStreamSliceCmd
}

func (m *MockRedisClient) XAdd(ctx context.Context, args *redis.XAddArgs) *redis.StringCmd {
	if m.XAddFunc != nil {
		return m.XAddFunc(ctx, args)
	}
	return redis.NewStringCmd(ctx)
}

func (m *MockRedisClient) XRead(ctx context.Context, args *redis.XReadArgs) *redis.XStreamSliceCmd {
	if m.XReadFunc != nil {
		return m.XReadFunc(ctx, args)
	}
	return redis.NewXStreamSliceCmd(ctx)
}

func (m *MockRedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	return redis.NewStatusCmd(ctx)
}

func (m *MockRedisClient) Close() error {
	return nil
}
