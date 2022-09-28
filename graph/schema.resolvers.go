package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"log"

	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/graph/generated"
	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/graph/model"
	redis "github.com/go-redis/redis/v9"
	"github.com/thanhpk/randstr"
)

// CreateMessage is the resolver for the createMessage field.
func (r *mutationResolver) CreateMessage(ctx context.Context, message string) (*model.Message, error) {
	m := model.Message{
		Message: message,
	}

	r.RedisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: "room",
		ID:     "*",
		MaxLen: 1, // max elements to store
		Values: map[string]interface{}{
			"message": m.Message,
		},
	})

	return &m, nil
}

// Messages is the resolver for the messages field.
func (r *queryResolver) Messages(ctx context.Context) ([]*model.Message, error) {
	streams, err := r.RedisClient.XRead(ctx, &redis.XReadArgs{
		Streams: []string{"room", "$"}, // "0" and remove Block to read all
		Block:   0,
	}).Result()
	if !errors.Is(err, nil) {
		return nil, err
	}

	stream := streams[0]

	ms := make([]*model.Message, len(stream.Messages))
	for i, v := range stream.Messages {
		ms[i] = &model.Message{
			ID:      v.ID,
			Message: v.Values["message"].(string),
		}
	}

	return ms, nil
}

// MessageCreated is the resolver for the messageCreated field.
func (r *subscriptionResolver) MessageCreated(ctx context.Context) (<-chan *model.Message, error) {
	token := randstr.Hex(16)
	mc := make(chan *model.Message, 1)
	r.mutex.Lock()
	r.messageChannels[token] = mc
	r.mutex.Unlock()

	go func() {
		<-ctx.Done()
		r.mutex.Lock()
		delete(r.messageChannels, token)
		r.mutex.Unlock()
		log.Println("Deleted")
	}()

	log.Println("Subscription: message created")

	return mc, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
