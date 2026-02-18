//go:generate go run github.com/99designs/gqlgen

package graph

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/graph/model"
	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/internal/datastore"
	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/internal/service"
)

type Resolver struct {
	RedisClient     datastore.RedisClient
	messageService  *service.MessageService
	messageChannels map[string]chan *model.Message
	mutex           sync.Mutex
}

func NewResolver(client datastore.RedisClient) *Resolver {
	return &Resolver{
		RedisClient:     client,
		messageService:  service.NewMessageService(client),
		messageChannels: map[string]chan *model.Message{},
		mutex:           sync.Mutex{},
	}
}

func (r *Resolver) SubscribeRedis(ctx context.Context) {
	log.Println("Start Redis Stream...")

	go func() {
		msgChan, errChan := r.messageService.StreamMessages(ctx)

		for {
			select {
			case <-ctx.Done():
				log.Println("Redis stream context cancelled")
				return
			case err, ok := <-errChan:
				if !ok {
					log.Println("Error channel closed")
					return
				}
				if !errors.Is(err, nil) {
					log.Printf("Error streaming messages: %v", err)
					return
				}
			case msg, ok := <-msgChan:
				if !ok {
					log.Println("Message channel closed")
					return
				}
				log.Printf("Received message: %s", msg.Message)

				r.mutex.Lock()
				for _, ch := range r.messageChannels {
					select {
					case ch <- msg:
					default:
						log.Println("Channel full, skipping message")
					}
				}
				r.mutex.Unlock()
			}
		}
	}()
}
