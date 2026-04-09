package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/v5"

	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/graph"
	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/internal/datastore"
	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/internal/graphql"
	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/internal/router"
)

// Version is a constant variable containing the version
const Version = "v0.0.1"

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func run() error {
	ctx := context.Background()

	redisURL := getEnv("REDIS_URL", "localhost:6379")
	serverPort := getEnv("PORT", ":8080")

	client, err := datastore.NewRedisClient(ctx, redisURL)
	if err != nil {
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

	// Echo v5 Start() handles graceful shutdown internally via signal.NotifyContext
	log.Printf("Starting server on %s", serverPort)
	if err := e.Start(serverPort); err != nil {
		return fmt.Errorf("server stopped: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
