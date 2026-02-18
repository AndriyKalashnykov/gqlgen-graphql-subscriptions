package datastore

import (
	"context"
	"errors"
	"testing"
)

func TestNewRedisClient_EmptyURL(t *testing.T) {
	ctx := context.Background()
	_, err := NewRedisClient(ctx, "")

	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

func TestNewRedisClient_InvalidURL(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	_, err := NewRedisClient(ctx, "invalid:99999")

	if err == nil {
		t.Fatal("expected error for invalid URL, got nil")
	}
}

// TestNewRedisClient_Success is an integration test that requires Redis
func TestNewRedisClient_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	client, err := NewRedisClient(ctx, "localhost:6379")

	if !errors.Is(err, nil) {
		t.Skipf("Redis not available: %v", err)
	}

	if client == nil {
		t.Fatal("expected client, got nil")
	}

	defer client.Close()

	// Verify connection
	err = client.Ping(ctx).Err()
	if !errors.Is(err, nil) {
		t.Errorf("failed to ping Redis: %v", err)
	}
}
