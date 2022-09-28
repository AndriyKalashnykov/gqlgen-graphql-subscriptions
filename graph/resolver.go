//go:generate go run github.com/99designs/gqlgen

package graph

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/go-redis/redis/v9"

	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/graph/model"
)

type Resolver struct {
	RedisClient     *redis.Client
	messageChannels map[string]chan *model.Message
	mutex           sync.Mutex
}

func NewResolver(client *redis.Client) *Resolver {
	return &Resolver{
		RedisClient:     client,
		messageChannels: map[string]chan *model.Message{},
		mutex:           sync.Mutex{},
	}
}

func (r *Resolver) SubscribeRedis(ctx context.Context) {
	log.Println("Start Redis Stream...")

	go func() {
		for {
			log.Println("Stream starting...")
			streams, err := r.RedisClient.XRead(ctx, &redis.XReadArgs{
				Streams: []string{"room", "$"},
				Count:   1,
				Block:   0,
			}).Result()
			if !errors.Is(err, nil) {
				panic(err)
			}

			stream := streams[0]
			m := &model.Message{
				ID:      stream.Messages[0].ID,
				Message: stream.Messages[0].Values["message"].(string),
			}
			r.mutex.Lock()
			for _, ch := range r.messageChannels {
				ch <- m
			}
			r.mutex.Unlock()

			log.Println("Stream finished...")
		}
	}()
}
