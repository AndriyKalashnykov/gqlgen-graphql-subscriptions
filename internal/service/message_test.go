package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/internal/constants"
	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/internal/testutil"
)

func TestNewMessageService(t *testing.T) {
	mock := &testutil.MockRedisClient{}
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
	mock := &testutil.MockRedisClient{
		XAddFunc: func(ctx context.Context, args *redis.XAddArgs) *redis.StringCmd {
			if args.Stream != constants.RedisStreamRoom {
				t.Errorf("expected stream %s, got %s", constants.RedisStreamRoom, args.Stream)
			}
			if args.MaxLen != constants.RedisStreamMaxLen {
				t.Errorf("expected maxlen %d, got %d", constants.RedisStreamMaxLen, args.MaxLen)
			}
			cmd := redis.NewStringCmd(ctx)
			cmd.SetVal("1234567890-0")
			return cmd
		},
	}

	svc := NewMessageService(mock)
	msg, err := svc.PublishMessage(ctx, "hello")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if msg == nil {
		t.Fatal("expected message, got nil")
	}

	if msg.Message != "hello" {
		t.Errorf("expected message 'hello', got %s", msg.Message)
	}

	if msg.ID != "1234567890-0" {
		t.Errorf("expected ID '1234567890-0', got %s", msg.ID)
	}
}

func TestPublishMessage_EmptyMessage(t *testing.T) {
	ctx := context.Background()
	mock := &testutil.MockRedisClient{}
	svc := NewMessageService(mock)

	_, err := svc.PublishMessage(ctx, "")

	if err == nil {
		t.Fatal("expected error for empty message, got nil")
	}
}

func TestPublishMessage_TooLong(t *testing.T) {
	ctx := context.Background()
	mock := &testutil.MockRedisClient{}
	svc := NewMessageService(mock)

	longMsg := make([]byte, constants.MaxMessageLength+1)
	for i := range longMsg {
		longMsg[i] = 'a'
	}

	_, err := svc.PublishMessage(ctx, string(longMsg))

	if err == nil {
		t.Fatal("expected error for message exceeding max length, got nil")
	}
}

func TestPublishMessage_RedisError(t *testing.T) {
	ctx := context.Background()
	redisErr := errors.New("redis connection error")
	mock := &testutil.MockRedisClient{
		XAddFunc: func(ctx context.Context, args *redis.XAddArgs) *redis.StringCmd {
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
	mock := &testutil.MockRedisClient{
		XReadFunc: func(ctx context.Context, args *redis.XReadArgs) *redis.XStreamSliceCmd {
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

	if err != nil {
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
	mock := &testutil.MockRedisClient{
		XReadFunc: func(ctx context.Context, args *redis.XReadArgs) *redis.XStreamSliceCmd {
			cmd := redis.NewXStreamSliceCmd(ctx)
			cmd.SetVal([]redis.XStream{})
			return cmd
		},
	}

	svc := NewMessageService(mock)
	messages, err := svc.ReadMessages(ctx)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(messages) != 0 {
		t.Errorf("expected empty messages, got %d", len(messages))
	}
}

func TestReadMessages_RedisNil(t *testing.T) {
	ctx := context.Background()
	mock := &testutil.MockRedisClient{
		XReadFunc: func(ctx context.Context, args *redis.XReadArgs) *redis.XStreamSliceCmd {
			cmd := redis.NewXStreamSliceCmd(ctx)
			cmd.SetErr(redis.Nil)
			return cmd
		},
	}

	svc := NewMessageService(mock)
	messages, err := svc.ReadMessages(ctx)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(messages) != 0 {
		t.Errorf("expected empty messages for redis.Nil, got %d", len(messages))
	}
}

func TestReadMessages_RedisError(t *testing.T) {
	ctx := context.Background()
	redisErr := errors.New("connection reset")
	mock := &testutil.MockRedisClient{
		XReadFunc: func(ctx context.Context, args *redis.XReadArgs) *redis.XStreamSliceCmd {
			cmd := redis.NewXStreamSliceCmd(ctx)
			cmd.SetErr(redisErr)
			return cmd
		},
	}

	svc := NewMessageService(mock)
	_, err := svc.ReadMessages(ctx)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, redisErr) {
		t.Errorf("expected error to wrap redis error")
	}
}

func TestStreamMessages_ReceivesMessage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	callCount := 0
	mock := &testutil.MockRedisClient{
		XReadFunc: func(ctx context.Context, args *redis.XReadArgs) *redis.XStreamSliceCmd {
			callCount++
			cmd := redis.NewXStreamSliceCmd(ctx)
			if callCount == 1 {
				cmd.SetVal([]redis.XStream{
					{
						Stream: constants.RedisStreamRoom,
						Messages: []redis.XMessage{
							{
								ID:     "100-0",
								Values: map[string]interface{}{constants.RedisMessageField: "streamed-msg"},
							},
						},
					},
				})
			} else {
				// Cancel after first message to stop the loop
				cancel()
				cmd.SetErr(context.Canceled)
			}
			return cmd
		},
	}

	svc := NewMessageService(mock)
	msgChan, errChan := svc.StreamMessages(ctx)

	select {
	case msg := <-msgChan:
		if msg.ID != "100-0" {
			t.Errorf("expected ID '100-0', got %s", msg.ID)
		}
		if msg.Message != "streamed-msg" {
			t.Errorf("expected message 'streamed-msg', got %s", msg.Message)
		}
	case err := <-errChan:
		t.Fatalf("unexpected error: %v", err)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for streamed message")
	}
}

func TestStreamMessages_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	mock := &testutil.MockRedisClient{
		XReadFunc: func(ctx context.Context, args *redis.XReadArgs) *redis.XStreamSliceCmd {
			cmd := redis.NewXStreamSliceCmd(ctx)
			cmd.SetErr(context.Canceled)
			return cmd
		},
	}

	svc := NewMessageService(mock)
	msgChan, _ := svc.StreamMessages(ctx)

	cancel()

	// Channel should close without sending messages
	select {
	case _, ok := <-msgChan:
		if ok {
			t.Error("expected channel to be closed")
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for channel close")
	}
}

func TestStreamMessages_RedisError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mock := &testutil.MockRedisClient{
		XReadFunc: func(ctx context.Context, args *redis.XReadArgs) *redis.XStreamSliceCmd {
			cmd := redis.NewXStreamSliceCmd(ctx)
			cmd.SetErr(errors.New("redis down"))
			return cmd
		},
	}

	svc := NewMessageService(mock)
	_, errChan := svc.StreamMessages(ctx)

	select {
	case err := <-errChan:
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for error")
	}
}

func TestStreamMessages_InvalidFormat(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mock := &testutil.MockRedisClient{
		XReadFunc: func(ctx context.Context, args *redis.XReadArgs) *redis.XStreamSliceCmd {
			cmd := redis.NewXStreamSliceCmd(ctx)
			cmd.SetVal([]redis.XStream{
				{
					Stream: constants.RedisStreamRoom,
					Messages: []redis.XMessage{
						{
							ID:     "1-0",
							Values: map[string]interface{}{constants.RedisMessageField: 12345},
						},
					},
				},
			})
			return cmd
		},
	}

	svc := NewMessageService(mock)
	_, errChan := svc.StreamMessages(ctx)

	select {
	case err := <-errChan:
		if err == nil {
			t.Fatal("expected error for invalid format, got nil")
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for error")
	}
}

func TestReadMessages_InvalidFormat(t *testing.T) {
	ctx := context.Background()
	mock := &testutil.MockRedisClient{
		XReadFunc: func(ctx context.Context, args *redis.XReadArgs) *redis.XStreamSliceCmd {
			cmd := redis.NewXStreamSliceCmd(ctx)
			cmd.SetVal([]redis.XStream{
				{
					Stream: constants.RedisStreamRoom,
					Messages: []redis.XMessage{
						{
							ID:     "1-0",
							Values: map[string]interface{}{constants.RedisMessageField: 12345},
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
