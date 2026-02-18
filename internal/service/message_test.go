package service

import (
	"context"
	"errors"
	"testing"

	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/internal/constants"
	"github.com/redis/go-redis/v9"
)

type mockRedisClient struct {
	xAddFunc  func(ctx context.Context, args *redis.XAddArgs) *redis.StringCmd
	xReadFunc func(ctx context.Context, args *redis.XReadArgs) *redis.XStreamSliceCmd
}

func (m *mockRedisClient) XAdd(ctx context.Context, args *redis.XAddArgs) *redis.StringCmd {
	if m.xAddFunc != nil {
		return m.xAddFunc(ctx, args)
	}
	return redis.NewStringCmd(ctx)
}

func (m *mockRedisClient) XRead(ctx context.Context, args *redis.XReadArgs) *redis.XStreamSliceCmd {
	if m.xReadFunc != nil {
		return m.xReadFunc(ctx, args)
	}
	return redis.NewXStreamSliceCmd(ctx)
}

func (m *mockRedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	return redis.NewStatusCmd(ctx)
}

func (m *mockRedisClient) Close() error {
	return nil
}

func TestNewMessageService(t *testing.T) {
	mock := &mockRedisClient{}
	svc := NewMessageService(mock)

	if svc == nil {
		t.Fatal("expected service to be created, got nil")
	}

	if svc.redis != mock {
		t.Error("expected redis client to be set correctly")
	}
}

func TestPublishMessage_Success(t *testing.T) {
	ctx := context.Background()
	mock := &mockRedisClient{
		xAddFunc: func(ctx context.Context, args *redis.XAddArgs) *redis.StringCmd {
			if args.Stream != constants.RedisStreamRoom {
				t.Errorf("expected stream %s, got %s", constants.RedisStreamRoom, args.Stream)
			}
			if args.MaxLen != constants.RedisStreamMaxLen {
				t.Errorf("expected maxlen %d, got %d", constants.RedisStreamMaxLen, args.MaxLen)
			}
			cmd := redis.NewStringCmd(ctx)
			cmd.SetVal("OK")
			return cmd
		},
	}

	svc := NewMessageService(mock)
	msg, err := svc.PublishMessage(ctx, "hello")

	if !errors.Is(err, nil) {
		t.Fatalf("unexpected error: %v", err)
	}

	if msg == nil {
		t.Fatal("expected message, got nil")
	}

	if msg.Message != "hello" {
		t.Errorf("expected message 'hello', got %s", msg.Message)
	}
}

func TestPublishMessage_EmptyMessage(t *testing.T) {
	ctx := context.Background()
	mock := &mockRedisClient{}
	svc := NewMessageService(mock)

	_, err := svc.PublishMessage(ctx, "")

	if err == nil {
		t.Fatal("expected error for empty message, got nil")
	}
}

func TestPublishMessage_RedisError(t *testing.T) {
	ctx := context.Background()
	redisErr := errors.New("redis connection error")
	mock := &mockRedisClient{
		xAddFunc: func(ctx context.Context, args *redis.XAddArgs) *redis.StringCmd {
			cmd := redis.NewStringCmd(ctx)
			cmd.SetErr(redisErr)
			return cmd
		},
	}

	svc := NewMessageService(mock)
	_, err := svc.PublishMessage(ctx, "hello")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, redisErr) {
		t.Errorf("expected error to wrap redis error")
	}
}

func TestReadMessages_Success(t *testing.T) {
	ctx := context.Background()
	mock := &mockRedisClient{
		xReadFunc: func(ctx context.Context, args *redis.XReadArgs) *redis.XStreamSliceCmd {
			cmd := redis.NewXStreamSliceCmd(ctx)
			cmd.SetVal([]redis.XStream{
				{
					Stream: constants.RedisStreamRoom,
					Messages: []redis.XMessage{
						{
							ID:     "1-0",
							Values: map[string]interface{}{constants.RedisMessageField: "message1"},
						},
						{
							ID:     "2-0",
							Values: map[string]interface{}{constants.RedisMessageField: "message2"},
						},
					},
				},
			})
			return cmd
		},
	}

	svc := NewMessageService(mock)
	messages, err := svc.ReadMessages(ctx)

	if !errors.Is(err, nil) {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(messages))
	}

	if messages[0].ID != "1-0" || messages[0].Message != "message1" {
		t.Errorf("unexpected first message: %+v", messages[0])
	}

	if messages[1].ID != "2-0" || messages[1].Message != "message2" {
		t.Errorf("unexpected second message: %+v", messages[1])
	}
}

func TestReadMessages_EmptyStream(t *testing.T) {
	ctx := context.Background()
	mock := &mockRedisClient{
		xReadFunc: func(ctx context.Context, args *redis.XReadArgs) *redis.XStreamSliceCmd {
			cmd := redis.NewXStreamSliceCmd(ctx)
			cmd.SetVal([]redis.XStream{})
			return cmd
		},
	}

	svc := NewMessageService(mock)
	messages, err := svc.ReadMessages(ctx)

	if !errors.Is(err, nil) {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(messages) != 0 {
		t.Errorf("expected empty messages, got %d", len(messages))
	}
}

func TestReadMessages_InvalidFormat(t *testing.T) {
	ctx := context.Background()
	mock := &mockRedisClient{
		xReadFunc: func(ctx context.Context, args *redis.XReadArgs) *redis.XStreamSliceCmd {
			cmd := redis.NewXStreamSliceCmd(ctx)
			cmd.SetVal([]redis.XStream{
				{
					Stream: constants.RedisStreamRoom,
					Messages: []redis.XMessage{
						{
							ID:     "1-0",
							Values: map[string]interface{}{constants.RedisMessageField: 12345}, // Invalid type
						},
					},
				},
			})
			return cmd
		},
	}

	svc := NewMessageService(mock)
	_, err := svc.ReadMessages(ctx)

	if err == nil {
		t.Fatal("expected error for invalid message format, got nil")
	}
}
