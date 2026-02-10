package main

import (
	"context"
	"errors"
	"log"

	"github.com/labstack/echo/v5"

	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/graph"
	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/internal/datastore"
	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/internal/graphql"
	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/internal/router"
)

// Version is a constant variable containing the version
const Version = "v0.0.1"

func main() {
	ctx := context.Background()

	client, err := datastore.NewRedisClient(ctx, "localhost:6379")
	if !errors.Is(err, nil) {
		log.Fatalln(err)
	}
	defer client.Close()

	r := graph.NewResolver(client)
	r.SubscribeRedis(ctx)
	srv := graphql.NewGraphQLServer(r)

	e := router.NewRouter(echo.New(), srv)
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
