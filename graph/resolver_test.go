package graph

import (
	"context"
	"testing"
	"time"

	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/graph/model"
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

func TestNewResolver(t *testing.T) {
	mock := &mockRedisClient{}
	resolver := NewResolver(mock)

	if resolver == nil {
		t.Fatal("expected resolver to be created, got nil")
	}

	if resolver.RedisClient != mock {
		t.Error("expected redis client to be set")
	}

	if resolver.messageService == nil {
		t.Error("expected message service to be initialized")
	}

	if resolver.messageChannels == nil {
		t.Error("expected message channels to be initialized")
	}
}

func TestMutationResolver_CreateMessage(t *testing.T) {
	ctx := context.Background()
	mock := &mockRedisClient{
		xAddFunc: func(ctx context.Context, args *redis.XAddArgs) *redis.StringCmd {
			cmd := redis.NewStringCmd(ctx)
			cmd.SetVal("OK")
			return cmd
		},
	}

	resolver := NewResolver(mock)
	mr := &mutationResolver{resolver}

	msg, err := mr.CreateMessage(ctx, "test message")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if msg == nil {
		t.Fatal("expected message, got nil")
	}

	if msg.Message != "test message" {
		t.Errorf("expected message 'test message', got %s", msg.Message)
	}
}

func TestMutationResolver_CreateMessage_Error(t *testing.T) {
	ctx := context.Background()
	mock := &mockRedisClient{}

	resolver := NewResolver(mock)
	mr := &mutationResolver{resolver}

	_, err := mr.CreateMessage(ctx, "")

	if err == nil {
		t.Fatal("expected error for empty message, got nil")
	}
}

func TestQueryResolver_Messages(t *testing.T) {
	ctx := context.Background()
	mock := &mockRedisClient{
		xReadFunc: func(ctx context.Context, args *redis.XReadArgs) *redis.XStreamSliceCmd {
			cmd := redis.NewXStreamSliceCmd(ctx)
			cmd.SetVal([]redis.XStream{
				{
					Stream: "room",
					Messages: []redis.XMessage{
						{
							ID:     "1-0",
							Values: map[string]interface{}{"message": "test1"},
						},
					},
				},
			})
			return cmd
		},
	}

	resolver := NewResolver(mock)
	qr := &queryResolver{resolver}

	messages, err := qr.Messages(ctx)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(messages))
	}

	if messages[0].Message != "test1" {
		t.Errorf("expected message 'test1', got %s", messages[0].Message)
	}
}

func TestSubscriptionResolver_MessageCreated(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mock := &mockRedisClient{}
	resolver := NewResolver(mock)
	sr := &subscriptionResolver{resolver}

	ch, err := sr.MessageCreated(ctx)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ch == nil {
		t.Fatal("expected channel, got nil")
	}

	// Verify channel is registered
	resolver.mutex.Lock()
	channelCount := len(resolver.messageChannels)
	resolver.mutex.Unlock()

	if channelCount != 1 {
		t.Errorf("expected 1 message channel, got %d", channelCount)
	}

	// Test cleanup on context cancellation
	cancel()
	time.Sleep(10 * time.Millisecond)

	resolver.mutex.Lock()
	channelCountAfter := len(resolver.messageChannels)
	resolver.mutex.Unlock()

	if channelCountAfter != 0 {
		t.Errorf("expected message channels to be cleaned up, got %d", channelCountAfter)
	}
}

func TestSubscriptionResolver_MessageDelivery(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mock := &mockRedisClient{}
	resolver := NewResolver(mock)
	sr := &subscriptionResolver{resolver}

	ch, err := sr.MessageCreated(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Simulate message delivery
	testMsg := &model.Message{ID: "1-0", Message: "test"}

	resolver.mutex.Lock()
	for _, msgCh := range resolver.messageChannels {
		msgCh <- testMsg
	}
	resolver.mutex.Unlock()

	// Verify message received
	select {
	case msg := <-ch:
		if msg.Message != "test" {
			t.Errorf("expected message 'test', got %s", msg.Message)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("timeout waiting for message")
	}
}
