package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/graph/model"
	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/internal/constants"
	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/internal/datastore"
)

// MessageService handles message publishing and retrieval via Redis
type MessageService struct {
	redis datastore.RedisClient
}

// NewMessageService creates a new MessageService
func NewMessageService(redis datastore.RedisClient) *MessageService {
	return &MessageService{
		redis: redis,
	}
}

// PublishMessage publishes a message to Redis stream
func (s *MessageService) PublishMessage(ctx context.Context, message string) (*model.Message, error) {
	if message == "" {
		return nil, fmt.Errorf("message cannot be empty")
	}
	if len(message) > constants.MaxMessageLength {
		return nil, fmt.Errorf("message exceeds maximum length of %d", constants.MaxMessageLength)
	}

	id, err := s.redis.XAdd(ctx, &redis.XAddArgs{
		Stream: constants.RedisStreamRoom,
		ID:     "*",
		MaxLen: constants.RedisStreamMaxLen,
		Values: map[string]interface{}{
			constants.RedisMessageField: message,
		},
	}).Result()

	if err != nil {
		return nil, fmt.Errorf("failed to publish message: %w", err)
	}

	return &model.Message{
		ID:      id,
		Message: message,
	}, nil
}

// ReadMessages reads messages from Redis stream
func (s *MessageService) ReadMessages(ctx context.Context) ([]*model.Message, error) {
	streams, err := s.redis.XRead(ctx, &redis.XReadArgs{
		Streams: []string{constants.RedisStreamRoom, "0"},
		Count:   constants.RedisStreamCount,
	}).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return []*model.Message{}, nil
		}
		return nil, fmt.Errorf("failed to read messages: %w", err)
	}

	if len(streams) == 0 {
		return []*model.Message{}, nil
	}

	stream := streams[0]
	messages := make([]*model.Message, len(stream.Messages))

	for i, v := range stream.Messages {
		msgValue, ok := v.Values[constants.RedisMessageField].(string)
		if !ok {
			return nil, fmt.Errorf("invalid message format at index %d", i)
		}

		messages[i] = &model.Message{
			ID:      v.ID,
			Message: msgValue,
		}
	}

	return messages, nil
}

// StreamMessages continuously reads messages from Redis stream and sends them to the channel
func (s *MessageService) StreamMessages(ctx context.Context) (<-chan *model.Message, <-chan error) {
	msgChan := make(chan *model.Message)
	errChan := make(chan error, 1)

	go func() {
		defer close(msgChan)
		defer close(errChan)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				streams, err := s.redis.XRead(ctx, &redis.XReadArgs{
					Streams: []string{constants.RedisStreamRoom, "$"},
					Count:   constants.RedisStreamCount,
					Block:   0,
				}).Result()

				if err != nil {
					if errors.Is(err, context.Canceled) {
						return
					}
					errChan <- fmt.Errorf("failed to stream messages: %w", err)
					return
				}

				if len(streams) > 0 {
					stream := streams[0]
					if len(stream.Messages) > 0 {
						msgValue, ok := stream.Messages[0].Values[constants.RedisMessageField].(string)
						if !ok {
							errChan <- fmt.Errorf("invalid message format in stream")
							return
						}

						msg := &model.Message{
							ID:      stream.Messages[0].ID,
							Message: msgValue,
						}

						select {
						case msgChan <- msg:
						case <-ctx.Done():
							return
						}
					}
				}
			}
		}
	}()

	return msgChan, errChan
}
