//go:generate go run github.com/99designs/gqlgen

package graph

import (
	"context"
	"log"
	"sync"
	"time"

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
		const maxRetries = 5
		retries := 0

		for {
			if retries > 0 {
				backoff := time.Duration(retries) * time.Second
				if backoff > 30*time.Second {
					backoff = 30 * time.Second
				}
				log.Printf("Retrying Redis stream in %v (attempt %d/%d)", backoff, retries, maxRetries)
				select {
				case <-ctx.Done():
					log.Println("Redis stream context cancelled during retry backoff")
					return
				case <-time.After(backoff):
				}
			}

			msgChan, errChan := r.messageService.StreamMessages(ctx)

			for {
				select {
				case <-ctx.Done():
					log.Println("Redis stream context cancelled")
					return
				case err, ok := <-errChan:
					if !ok {
						log.Println("Error channel closed")
						goto retry
					}
					if err != nil {
						log.Printf("Error streaming messages: %v", err)
						goto retry
					}
				case msg, ok := <-msgChan:
					if !ok {
						log.Println("Message channel closed")
						goto retry
					}
					retries = 0 // reset on successful message
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

		retry:
			retries++
			if retries > maxRetries {
				log.Printf("Redis stream failed after %d retries, giving up", maxRetries)
				return
			}
		}
	}()
}
