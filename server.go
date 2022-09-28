package main

import (
	"errors"
	"log"

	"github.com/labstack/echo/v4"

	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/graph"
	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/infrastructure/datastore"
	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/infrastructure/graphql"
	"github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/infrastructure/router"
)

// Version is a constant variable containing the version
const Version = "v0.0.1"

func main() {
	client, err := datastore.NewRedisClient("localhost:6379")
	if !errors.Is(err, nil) {
		log.Fatalln(err)
	}
	defer client.Close()

	r := graph.NewResolver(client)
	r.SubscribeRedis()
	srv := graphql.NewGraphQLServer(r)

	e := router.NewRouter(echo.New(), srv)
	e.Logger.Fatal(e.Start(":8080"))
}
