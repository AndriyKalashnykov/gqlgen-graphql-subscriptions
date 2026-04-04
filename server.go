package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/v5"

	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/graph"
	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/internal/constants"
	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/internal/datastore"
	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/internal/graphql"
	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/internal/router"
)

// Version is a constant variable containing the version
const Version = "v0.0.1"

const redisURL = "localhost:6379"

func run() error {
	ctx := context.Background()

	client, err := datastore.NewRedisClient(ctx, redisURL)
	if !errors.Is(err, nil) {
		return fmt.Errorf("failed to connect to Redis at %s: %w", redisURL, err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Error closing Redis client: %v", err)
		}
	}()

	r := graph.NewResolver(client)
	r.SubscribeRedis(ctx)
	srv := graphql.NewGraphQLServer(r)

	e := router.NewRouter(echo.New(), srv)

	log.Printf("Starting server on %s", constants.ServerPort)
	if err := e.Start(constants.ServerPort); err != nil {
		return fmt.Errorf("server failed to start: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
